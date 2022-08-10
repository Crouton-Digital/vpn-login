FROM golang:1.18.5 as builder

COPY . /opt/
RUN cd /opt && rm -f vpn-login && go mod vendor && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o vpn-login .


# Smallest base image
FROM alpine:3.15.5

# Install needed packages
RUN echo "http://dl-4.alpinelinux.org/alpine/edge/community/" >> /etc/apk/repositories && \
    echo "http://dl-4.alpinelinux.org/alpine/edge/testing/" >> /etc/apk/repositories && \
    apk update && apk add openssl easy-rsa openvpn openvpn-auth-ldap openvpn-auth-pam openvpn-auth-script iptables bash  && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

RUN mkdir /opt/openvpn
# ADD server.conf for OpenVPN
#COPY openvpn/ /opt/openvpn
COPY --from=builder /opt/vpn-login /opt/openvpn/
RUN chmod a+x /opt/openvpn/vpn-login
#ADD settings.yml /

ADD entrypoint.sh /
RUN chmod u+x entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
