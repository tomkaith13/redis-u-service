# Redis BloomFilter + CuckooFilter Service

basics on BloomFilters https://redis.io/docs/stack/bloom/
This is minimalistic u-service that spins up a redis DB and uses it for membership operations in  the bloom filter.
One additional op that cuckoo filter has that bloom filter does not is deleting an item from a key.
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

## Testing setup

Run
```
make up
```
and run
```
POST /bf-test-setup HTTP/1.1
Host: localhost:8080
```
on cURL or Postman to test the test route `/bf-test-setup`.
The first time this runs, you should create a key of name `testBF` and insert key named `works` with status 201. The second time you run this, it should return 409 status

## Reserve a new BloomFilter

See https://redis.io/commands/bf.reserve/ for details on the params.

```
POST /bf-reserve HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 65

{
    "name": "BF", 
    "errorRate": 0.01,
    "capacity": 1000,
    "ttl_in_secs": 0
}
```
By default, the ttl can be `0` if not provided in body and the filter is not ephemeral.

### Status codes
#### 201
If a bloomfilter with name is created
#### 409
If one with same name already exists

#### 500
If other errors


## Insert a item to an existing BloomFilter

See this for details https://redis.io/commands/bf.insert/

```
POST /bf-insert HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 40

{
    "keyName": "BF",
    "item": "b"
}
```
### Status Codes
#### 201
Added a new item

#### 409 
Item may already exist in BloomFilter

#### 404
BloomFilter with keyName does not exist.
User needs to use POST /bf-reserve to create a new one

## Check membership
See https://redis.io/commands/bf.exists/ for details
The HTTP request looks like :
```
GET /bf-exists HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 40

{
    "keyName": "BF",
    "item": "a"
}
```

### Status Codes
#### 200
Returns `definitely does not exist` or `maybe exists` based on the entries

#### 500
Returned if there is an error.

## Delete a BloomFilter key
See https://redis.io/commands/del/

This route can be used to delete a bloomfilter from the redis server.
You can use this:
```
DELETE /bf HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 23

{
    "keyName": "BF"
}
```

### Status Codes
#### 200
Deletion worked

#### 404
Key not found

## Cuckoo-filter delete an item

We can use `DEL /cf` endpoint to delete an item from a key
```
DELETE /cf HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 23

{
    "keyName": "CF"
}
```

### Status Codes
#### 200
If the delete was successful

#### 404
 if the item does not exist


## Cleanup
Run `make clean` to tear down the services and remove all containers.
