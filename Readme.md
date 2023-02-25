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
make compose-up
```
and run
```
POST /bfadd-test HTTP/1.1
Host: localhost:8080
```
on cURL or Postman to test the test route `/bfadd-test` and it should return a `it works!!`  response in stdout.


