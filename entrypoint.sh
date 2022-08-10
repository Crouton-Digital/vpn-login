#!/bin/bash

set -e

mkdir -p /dev/net
mknod /dev/net/tun c 10 200
chmod 600 /dev/net/tun

iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE

/bin/sleep 2
openvpn --config /opt/openvpn/server.conf

