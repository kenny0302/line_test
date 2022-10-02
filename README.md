# Start up

## Create MongoDB Container
```
docker run -d -name testdb -p 27017:27017 mongo:4.4
```

## Set Test DB and Collection

```
docker exec -it testdb bash
mongo
use line
db.createCollection("message")
db.createCollection("user")
```

## Fix Config.yaml

Create a test line dev official account & Message API

Fix line.secret & line.token

## Run Service
```
go run main.go
```

## Test Video

### Spark

```
curl --location --request POST '127.0.0.1:8080/push' \
--form 'UserId="$UserId"'
```

https://drive.google.com/file/d/1b4kTe0RY5kwcwQEO-N_2UDMvanV7KaC_/view?usp=sharing

### Save User info and Message from line bot

https://drive.google.com/file/d/1P7bOGYC2t3TtN2MJGBcH39WGxuPWbkux/view?usp=sharing

### List User collection

```
curl --location --request GET '127.0.0.1:8080/list'
```

https://drive.google.com/file/d/1U_iD4vrJFMmqYWKFkOVFM7wx6bYh3uIt/view?usp=sharing
