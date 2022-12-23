# Auth script for OpenVPN for support AWS Cognito

That script make auth VPN user with user credentional like AWS Cognito and etc

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

* Go to [Cloud Build](https://console.cloud.google.com/cloud-build/triggers?project=asset-management-ci-cd&pageState=(%22triggers%22:(%22f%22:%22%255B%257B_22k_22_3A_22_22_2C_22t_22_3A10_2C_22v_22_3A_22_5C_22OpenVPN-build_5C_22_22_2C_22s_22_3Atrue%257D%255D%22)))

* clone project into your machine
* install dependencies: `go mod download`
* copy `.env.example` into `.env`, tune values if necessary
* run project: `go run .`


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

## Improvements

### TODO

### BUGS

### CHANGELOG

- 2022-08-09 -- development start

- 2021-08-10 -- run in production 

## Authors

- [@Rolands Kalpins](https://hodlersassetm-iro9206.slack.com/team/U027H05SHS4)
- Maksim Dogonov