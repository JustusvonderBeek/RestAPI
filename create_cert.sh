#!/bin/bash

echo "Creating a new certificate and key"

# Create a new RSA key and certificate for 3 years (3*365 days)
openssl req -x509 -newkey rsa:4096 -keyout vocabulary_local.key -out vocabulary_local.cer -sha256 -days 1095 -nodes -subj "/C=DE/ST=Bavaria/L=Munich/O=Cloudsheeptech/OU=Vocabulary/CN=10.0.2.2"

echo "Created new certificate under 'vocabulary.cer' and 'vocabulary.key'"