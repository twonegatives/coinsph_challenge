
-- +migrate Up
CREATE TYPE currency AS ENUM('usd');

CREATE TABLE accounts (
  id serial,
  name varchar UNIQUE NOT NULL,
  balance decimal NOT NULL,
  currency currency NOT NULL,
  PRIMARY KEY(id)
);

-- +migrate Down

DROP TABLE IF EXISTS accounts;
DROP TYPE IF EXISTS currency;
