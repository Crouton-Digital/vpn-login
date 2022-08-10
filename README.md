# VPN Login script

# Generate certificate for OpenVPN Server
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
# Generate certificate for client

```
mkdir ../client
./easyrsa build-client-full client nopass  
cp pki/private/client.key ../client/ 
cp pki/issued/client.crt ../client/
```