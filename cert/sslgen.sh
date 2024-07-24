#!/bin/bash
COUNTRY="RU" #Country code
STATE="Primorsky Region" #State
LOCATION="Vladivostok" #City for ex.
ORGANIZATION_UNIT="IT-Express Ltd" #Your organization
SERVER_CN="localhost" #your domain. most important field. for dev localhost is ok

#Creating directory structure

mkdir -p ./{ca,server,client}/{private,public}

#CA
openssl genrsa -passout pass:asdrfv26 -des3 -out ca/private/ca.key 2048
openssl req -passin pass:asdrfv26 -new -x509 -days 365000 -key ca/private/ca.key -out ca/public/ca.crt -subj "/C=${COUNTRY}/ST=${STATE}/L=${LOCATION}/O=${ORGANIZATION_UNIT}/CN=${SERVER_CN}"
#SERVER
openssl genrsa -passout pass:asdrfv26 -des3 -out server/private/server.key 2048
openssl req -passin pass:asdrfv26 -new -key server/private/server.key -out server/private/server.csr -subj "/C=${COUNTRY}/ST=${STATE}/L=${LOCATION}/O=${ORGANIZATION_UNIT}/CN=${SERVER_CN}"
openssl x509 -req -passin pass:asdrfv26 -extfile <(printf "subjectAltName=DNS:${SERVER_CN}") -days 365000 -in server/private/server.csr -CA ca/public/ca.crt -CAkey ca/private/ca.key -set_serial 01 -out server/public/server.crt
openssl pkcs8 -topk8 -nocrypt -passin pass:asdrfv26 -in server/private/server.key -out server/private/server.pem
#CLIENT
openssl genrsa -passout pass:asdrfv26 -des3 -out client/private/client.key 2048
openssl req -passin pass:asdrfv26 -new -key client/private/client.key -out client/private/client.csr --subj "/C=${COUNTRY}/ST=${STATE}/L=${LOCATION}/O=${ORGANIZATION_UNIT}/CN=${SERVER_CN}"
openssl x509 -req -passin pass:asdrfv26 -extfile <(printf "subjectAltName=DNS:${SERVER_CN}") -days 365000 -in client/private/client.csr -CA ca/public/ca.crt -CAkey ca/private/ca.key -set_serial 02 -out client/public/client.crt
openssl pkcs8 -topk8 -nocrypt -passin pass:asdrfv26 -in client/private/client.key -out client/private/client.pem

