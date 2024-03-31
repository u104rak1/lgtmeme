# lgtmeme

## Redis
```
docker exec -it lgtmeme_redis redis-cli
keys *
get [key]
```

## Grant execution permission to the script
```
chmod +x ./script/*.sh
```

## Generate RSA key pair
```
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
cat private_key.pem | base64 | tr -d '\n'
cat public_key.pem | base64 | tr -d '\n'
```