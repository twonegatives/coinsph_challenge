
-- +migrate Up
CREATE TYPE currency AS ENUM('usd');

CREATE TABLE accounts (
  id serial,
  name varchar NOT NULL,
  balance decimal NOT NULL,
  currency currency NOT NULL,
  PRIMARY KEY(id)
);

INSERT INTO accounts(name, balance, currency)
VALUES ('SYSTEM', 0, 'usd');

ALTER TABLE accounts ADD CONSTRAINT valid_balance CHECK (balance >= 0 OR id = 1);

-- +migrate Down

DROP TABLE IF EXISTS accounts;
DROP TYPE IF EXISTS currency;
