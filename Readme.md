# Local dev
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

### dockerized local dev with redis in a separate container

Run
```
make compose-up
```
and logs should show up in stdout
