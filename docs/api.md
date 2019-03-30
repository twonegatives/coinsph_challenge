# Wallet REST API

## Accounts

There are no pre-generated accounts (except `SYSTEM`) in the application database.
In order to obtain access to the whole application functionality you're recommended to create a couple of new accounts first.

### Create account

- _Method_: `POST`
- _URL_: `/api/v1/accounts`
- _Payload_: Nested JSON object containing account name
- _Response_: JSON struct of created account
- _Exception_: `400` on request with blank account name
- _Exception_: `500` on request with duplicated account name

_Examples_:
```bash
> curl -v -X POST localhost:8090/api/v1/accounts -d '{"account": {"name": "john_doe"}}'
< HTTP/1.1 200 OK
< {"account":{"name":"john_doe","balance":"0","currency":"usd"}}
```

```bash
> curl -v -X POST localhost:8090/api/v1/accounts -d '{"account": {"name": "john_doe"}}'
< HTTP/1.1 500 Internal Server Error
{"error":"failed to create new account in database: can't create new account: pq: duplicate key value violates unique constraint \"accounts_name_key\""}
```

```bash
> curl -v -X POST localhost:8090/api/v1/accounts -d '{"account": {"name": ""}}'
< HTTP/1.1 400 Bad Request
< {"error":"bad request"}
```

_Note_: account name is the only attribute consumed by Account creation API. Balance and currency are automatically set up to `0` and `usd` respectively.

### Get accounts list

- _Method_: `GET`
- _URL_: `/api/v1/accounts`
- _Response_: JSON array of existing accounts

_Examples_:
```bash
> curl -v localhost:8090/api/v1/accounts
< HTTP/1.1 200 OK
< {"accounts":[{"name":"SYSTEM","balance":"-190","currency":"usd"},{"name":"john_doe","balance":"190","currency":"usd"}]}
```

## Payments

### Create payment

- _Method_: `GET`
- _URL_: `/api/v1/payments`
- _Payload_: Nested JSON object containing sender/receiver names and amount
- _Response_: Blank JSON
- _Exception_: `400` on request with blank sender/receiver names
- _Exception_: `400` on payment amount less or equal to zero
- _Exception_: `500` on payment which sets user balance below zero
- _Exception_: `500` when sender and receiver is the same person

_Examples_:
```bash
> curl -v -X POST localhost:8090/api/v1/payments -d '{"payment" : {"from": "SYSTEM", "to": "john_doe", "amount": 10.12}}'
< HTTP/1.1 200 OK
< {}
```

### Get payments list

- _Method_: `GET`
- _URL_: `/api/v1/payments`
- _Response_: JSON array of existing payments

_Examples_:
```bash
> curl -v localhost:8090/api/v1/payments
< HTTP/1.1 200 OK
< {"payments":[{"account":"SYSTEM","amount":"180","currency":"usd","direction":"outgoing","to_account":"john_doe"},{"account":"john_doe","amount":"180","currency":"usd","direction":"incoming","from_account":"SYSTEM"}]}
```
