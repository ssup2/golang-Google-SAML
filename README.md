# golang-Google-SAML

* For studying SAML
* Reference Code : https://github.com/crewjam/saml/blob/main/example/trivial/trivial.go

* Generate Certificate
```
# openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=myservice.example.com"
```

* Generate SAML Metadata
```
# go run main.go
# curl localhost:8000/saml/metadata > metadata
```