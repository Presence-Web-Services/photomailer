# photomailer

Include `.env` in local install to setup email config.

How to run:
```
go mod download
go build -o ./photomailer-server
./photomailer-server
```

Running in Docker container:
```
docker build -t photomailer .
docker run -p 80:80 -d photomailer
```
