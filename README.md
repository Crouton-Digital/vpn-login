# Auth script for OpenVPN for support AWS Cognito

That script make auth VPN user with user credentional like AWS Cognito and etc

## Table of Contents

[TOC]

## About

This service

- Make auth AWS Cognito users and check security groups

Important notes:

* Only for OpenVPN 2.5 version 
* Build in Docker container

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

To use this project, you need `Docker` installed.

### Building

* Go to Google Cloud Build, search and run `OpenVPN-build` trigger and up version

## Generate certificate for OpenVPN Server
```sh
https://github.com/OpenVPN/easy-rsa/releases/download/v3.1.0/EasyRSA-3.1.0.tgz
tar -xvf EasyRSA-3.1.0.tgz
cp vars.example vars

edit file vars

mv vars pki
. ./pki/vars
./easyrsa init-pki
./easyrsa build-ca nopass
./easyrsa build-server-full server nopass
./easyrsa gen-dh

mkdir ../server-keys
cp pki/dh.pem ../server-keys
cp pki/private/server.key ../server-keys
cp pki/issued/server.crt ../server-keys
cp pki/ca.crt ../server-keys/
```
## Generate certificate for client

```sh
mkdir ../client
./easyrsa build-client-full client nopass  
cp pki/private/client.key ../client/ 
cp pki/issued/client.crt ../client/
```
## Server OpenVPN Configs
Create server.conf file for OpenVPN Server

```sh
port 1194
proto udp
dev tun

persist-key
persist-tun

server 10.8.0.0 255.255.255.0
keepalive 10 120
cipher AES-256-CBC
;compress lz4-v2
;push "compress lz4-v2"
;comp-lzo
# verb 3

status /opt/openvpn/openvpn-status.log
# ifconfig-pool-persist ipp.txt
;server-bridge 10.8.0.4 255.255.255.0 10.8.0.50 10.8.0.100
;server-bridge
; ---- PUSH ROUTES KUBERNETERS CLUSTER, SQL, REDIS ETC SUBNETS 
push "route 10.6.0.0 255.255.0.0"
push "route 10.7.0.0 255.255.0.0"
; ---- PUSH DHCP OPTIONS FOR CLIENTS
push "dhcp-option DOMAIN svc.cluster.local"
push "dhcp-option DOMAIN-SEARCH svc.cluster.local"
push "dhcp-option DOMAIN-SEARCH default.svc.cluster.local"
push "dhcp-option ADAPTER_DOMAIN_SUFFIX svc.cluster.local"
push "dhcp-option ADAPTER_DOMAIN_SUFFIX default.svc.cluster.local"
push "dhcp-option DNS 10.7.0.10"
push "dhcp-option DNS 8.8.8.8"

<key>
-----BEGIN PRIVATE KEY-----
  < server key from file >
-----END PRIVATE KEY-----
</key>
<cert>
-----BEGIN CERTIFICATE-----
  < server certificate from file >
-----END CERTIFICATE-----
</cert>
<ca>
-----BEGIN CERTIFICATE-----
  < CA certificate from file >
-----END CERTIFICATE-----
</ca>
<dh>
-----BEGIN DH PARAMETERS-----
  < DH from file >
-----END DH PARAMETERS-----
</dh>
key-direction 1
# <tls-auth>
# </tls-auth>

#redirect-gateway def1

plugin /usr/lib/openvpn/plugins/openvpn-plugin-auth-script.so /opt/openvpn/vpn-login
username-as-common-name
#verify-client-cert optional
duplicate-cn
```

Create settings.conf for Auth plugin
```sh
LOG_ENABLED: 1
LOG_FILE: access.log
COGNITO_REGION: eu-west-1
COGNITO_USER_POOL_ID: "eu-west-1_*****"
COGNITO_APP_CLIENT_ID: "5guc7*********18eu"
COGNITO_APP_CLIENT_SECRET: "lrjdt*********************************q2l0"
COGNITO_SCOPE: "all"
COGNITO_USER_GROUP: <Cognito VPN Group name>
COGNITO_ISSUER: https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_******
COGNITO_KEYS_URL: https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_*******/.well-known/jwks.json
```

## Client OpenVPN Config

Create file <name connection>.ovpn
```sh
client
persist-tun
persist-key
auth-user-pass
nobind
dev tun
remote < ip adress you can get kubectl get svc -o wide -n vpn > 1194 udp
resolv-retry infinite
 
<key>
-----BEGIN PRIVATE KEY-----
   < client key >
-----END PRIVATE KEY-----
</key>
<cert>
-----BEGIN CERTIFICATE-----
    < client crt >
-----END CERTIFICATE-----
</cert>
<ca>
-----BEGIN CERTIFICATE-----
     < ca crt >
-----END CERTIFICATE-----
</ca>
 
#redirect-gateway def1
```


## Improvements

### TODO

### BUGS

### CHANGELOG

- 2022-08-09 -- development start

- 2021-08-10 -- run in production 

## Authors

- [@Rolands Kalpins](https://hodlersassetm-iro9206.slack.com/team/U027H05SHS4)
- Maksim Dogonov