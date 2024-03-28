# lgtmeme

## Redis
```
docker exec -it lgtmeme_redis redis-cli
keys *
get [key]
```

## /api/connect/authorize
[login](http://localhost:8080/api/connect/authorize?response_type=code&client_id=a74983c2-c578-41fd-993b-9e4716d244ac&redirect_uri=http://localhost:3000/api/auth/callback&scope=images.read%20images.create%20images.update%20images.delete&state=xyz&nonce=abc123)

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