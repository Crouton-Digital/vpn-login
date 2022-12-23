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
```
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

```
mkdir ../client
./easyrsa build-client-full client nopass  
cp pki/private/client.key ../client/ 
cp pki/issued/client.crt ../client/
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