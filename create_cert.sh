#!/bin/bash

echo "Creating a new certificate and key"

keyfile="vocabulary.key"
certfile="vocabulary.cer"

# Create a new RSA key and certificate for 3 years (3*365 days)
openssl req -x509 -newkey rsa:4096 -keyout $keyfile -out $certfile -sha256 -days 365 -nodes -subj "/C=DE/ST=Bavaria/L=Munich/O=Cloudsheeptech/OU=Vocabulary/CN=vocabulary.cloudsheeptech.com"

echo "Created new certificate under '$keyfile' and '$certfile'"
