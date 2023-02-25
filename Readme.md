# Local dev setup
just run:

```
make run
```
and run
```
POST /bfadd-test HTTP/1.1
Host: localhost:8080
```
on cURL or Postman to test the test route `/bfadd-test` and it should return a `it works!!` response

# Dockerized local dev with redis as a separate service

Run
```
make up
```
and run
```
POST /bfadd-setup HTTP/1.1
Host: localhost:8080
```
on cURL or Postman to test the test route `/bfadd-setup`.
The first time this runs, you should create a key of name `testBF` and insert key named `works` with status 201. The second time you run this, it should return 409 status



