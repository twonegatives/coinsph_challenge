# Wallet REST API

## Accounts

There are no pre-generated accounts (except `SYSTEM`) in the application database.
In order to obtain access to the whole application functionality you're recommended to create a couple of new accounts first.

### Create account

- __Method__: `POST`
- __URL__: `/api/v1/accounts`
- __Payload__: Nested JSON object containing account name
- __Response__: JSON struct of created account
- __Exception__: `400` on request with blank account name
- __Exception__: `500` on request with duplicated account name

__Examples__:
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

__Note__: account name is the only attribute consumed by Account creation API. Balance and currency are automatically set up to `0` and `usd` respectively.

### Get accounts list

- __Method__: `GET`
- __URL__: `/api/v1/accounts`
- __Response__: JSON array of existing accounts

__Examples__:
```bash
> curl -v localhost:8090/api/v1/accounts
< HTTP/1.1 200 OK
< {"accounts":[{"name":"SYSTEM","balance":"-190","currency":"usd"},{"name":"john_doe","balance":"190","currency":"usd"}]}
```

## Payments

### Create payment

- __Method__: `GET`
- __URL__: `/api/v1/payments`
- __Payload__: Nested JSON object containing sender/receiver names and amount
- __Response__: Blank JSON
- __Exception__: `400` on request with blank sender/receiver names
- __Exception__: `400` on payment amount less or equal to zero
- __Exception__: `500` on payment which sets user balance below zero
- __Exception__: `500` when sender and receiver is the same person

__Examples__:
```bash
> curl -v -X POST localhost:8090/api/v1/payments -d '{"payment" : {"from": "SYSTEM", "to": "john_doe", "amount": 10.12}}'
< HTTP/1.1 200 OK
< {}
```

### Get payments list

- __Method__: `GET`
- __URL__: `/api/v1/payments`
- __Response__: JSON array of existing payments

__Examples__:
```bash
> curl -v localhost:8090/api/v1/payments
< HTTP/1.1 200 OK
< {"payments":[{"account":"SYSTEM","amount":"180","currency":"usd","direction":"outgoing","to_account":"john_doe"},{"account":"john_doe","amount":"180","currency":"usd","direction":"incoming","from_account":"SYSTEM"}]}
```
