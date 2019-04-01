[![Build Status](https://travis-ci.com/twonegatives/coinsph_challenge.svg?branch=master)](https://travis-ci.com/twonegatives/coinsph_challenge)
[![codecov](https://codecov.io/gh/twonegatives/coinsph_challenge/branch/master/graph/badge.svg)](https://codecov.io/gh/twonegatives/coinsph_challenge)
[![Maintainability](https://api.codeclimate.com/v1/badges/531e7ad5551a8ed8bd49/maintainability)](https://codeclimate.com/github/twonegatives/coinsph_challenge/maintainability)

# Coins.ph challenge

Golang code challange was to implement a generic Wallet service with a RESTful API.
Here is a list of requested user stories for the MVP phase:
1. See all payments;
2. See all accounts;
3. Send a payment from one account to another (same currency).

## Design decisions
### Bookkeeping system type
This wallet implementation uses a so-called [double-entry bookkeeeping system](https://en.wikipedia.org/wiki/Double-entry_bookkeeping_system).
Each transaction from one account to another gets stored as two equal and oppositee payment entries (the `outgoing` for sender and `incoming` for receiver).
With every transaction recorded as both a debit and a credit, the totals of each should always balance.
In other words, a value difference between `outgoing` and `incoming` payments of the same transaction indicate an error.

### Money flows over the system borders
It's a pretty rare case when the system exists in a vacuum and does not deal with the outer world.
Users would typically like to deposit their funds on their account (or withdraw it as a cash).
In order for money not to appear from nowhere, there is a special Account named `SYSTEM`.
Any money transfer from/to the outer world is done with the participation of this Account.
Please note that this account has a difference to all other (user) Accounts: `SYSTEM` may have its balance go below zero.
In fact, the less `SYSTEM` balance is, the more money users deposited into the wallet, so it's rather a happy scenario, yay!

### Data integrity checks
Any bookkeeeping system has a number of possible data integrity issues, to name a few:
- Difference betweeen `outgoing` and `incoming` payments of the same transaction;
- User account's balance goes below zero (for debit users);
- User account's balance does not match with all of his `incoming` and `outgoing` payments etc.

In order to address the mentioned issues wallet makes use of the following techniques:
1. Database transactions and row locks to overcome concurrent use cases;
2. Database constraints to guarantee the correctness of unique and non-zero fields;
3. Database `check` triggers (
[#1](https://github.com/twonegatives/coinsph_challenge/blob/master/migrations/20190329162827-CreateTxTriggerFunction.sql),
[#2](https://github.com/twonegatives/coinsph_challenge/blob/master/migrations/20190330002201-CreateAccountsTriggerFunctionn.sql)
) on insert/update to guarantee business rules fulfillment.

These arrangements provide a solid confidence in data integrity, though do not cover some nasty cases which may arise if someone makes changes using the db client directly on production servers.
A complete bulletproof solution would require more restrictive trigger policies which was intentionally left out of the scope for this phase.

## Web API

Check [api.md](https://github.com/twonegatives/coinsph_challenge/blob/master/docs/api.md) file for API documentation.

## Installation
* [Golang v. 1.11.1+](https://golang.org/dl)
* [SqlMigrate](https://github.com/rubenv/sql-migrate)
* [PostgreSQL 9.6+](https://www.postgresql.org/download/)

Setup database:

```bash
createdb coinsph
sql-migrate up
```

[Go modules](https://github.com/golang/go/wiki/Modules) and its builtin dependency management system are used to fetch
all third-party libraries. Just `build`:

```bash
go build ./cmd/service
```

or `run`:

```bash
LISTEN=:8090 go run cmd/service/*.go
```

to get it started.

## Testing

Wallet uses:

* [testing](https://golang.org/pkg/testing/) package from stdlib
* [testify](https://github.com/stretchr/testify) for `assert`/`require` syntatic sugar
* [mockgen](https://github.com/golang/mock) for interfaces auto-creation

Wallet does not require you to create a test database manually.
Instead on tests run it does the following automatically:
- connects to the development database;
- creates a test db;
- migrates a test db.
Please note that this requires you to have a dev db (default is `coinsph`) to be set up prior to firing up tests.

Run tests:

```bash
go test -v ./...
```

## Configuration
A set of environment variables might be provided to alter the Waller behaviour:

- `LISTEN` - `host:port` for server. Default: `:80`
- `APP_ENV` - application environment. Used by `sql-migrate` to pick according db configuration from `dbconf.yml` on migrations run. Default: `dev`
- `DB` - database connection string. Application server connects to this database on startup. Default: `postgres://localhost/coinsph?sslmode=disable`
- `SHUTDOWN_TIMEOUT` - timeout for gracefull server stop on exceptional cases (e.g. interruption). Default: `2s`

## Deployment
There is a [Dockerfile](https://github.com/twonegatives/coinsph_challenge/blob/master/Dockerfile) to help you get up and running:

```bash
docker build -t coinsph-challenge -f Dockerfile .
docker run -it --rm --name coinsph-running-app -p 4000:80 --env DB="postgres://DB_USER:DB_PASSWORD@DB_HOST:DB_PORT/DB_NAME?sslmode=disable" coinsph-challenge
```

Please note that database creation and migration are left out of `Dockerfile` scope and should be implemented separately.
Either `docker-compose`, `kubernetes` or a manual `docker run` should do the trick.
Database migration script is accessible at [bin/migrate](https://github.com/twonegatives/coinsph_challenge/blob/master/bin/migrate).
Please pay attention to `ENV` variables listed at [dbconfig.yml](https://github.com/twonegatives/coinsph_challenge/blob/master/dbconfig.yml) and `APP_ENV`.
```bash
docker run --env APP_ENV=stage --env STAGE_DB_USER_NAME=twonegatives --env STAGE_DB_HOST=docker.for.mac.host.internal:5432 --env STAGE_DB_NAME=coinsph -it coinsph-challenge /bin/migrate
```

## Suggestions? Bugs? Contributions?
If you've got a question, feature suggestion or found a bug please add an [issue](https://github.com/twonegatives/coinsph_challenge/issues) on GitHub or fork the project and send a pull request.
