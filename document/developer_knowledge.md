# Developer knowledge

## List of services in use

- [Github](https://github.com/ucho456job/lgtmeme)
- [Render](https://dashboard.render.com)
- [supabase](https://supabase.com/dashboard/projects)
- [Upstash](https://console.upstash.com/redis)
- [UptimeRobot](https://dashboard.uptimerobot.com/monitors)

## Frequently used commands

### Redis
```
docker exec -it lgtmeme_redis redis-cli   // Start redis-cli
keys *                                    // list keys
get [key]                                 // show value
set [key] [value]                         // save key value
FLUSHALL                                  // reset keys
```

### Grant execution permission to the script
```
chmod +x ./script/*.sh
```

### Generate new RSA key pair
```
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
cat private_key.pem | base64 | tr -d '\n'
cat public_key.pem | base64 | tr -d '\n'
```