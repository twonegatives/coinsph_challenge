
-- +migrate Up
CREATE TYPE direction AS ENUM('incoming','outgoing');

CREATE TABLE payments (
  id serial,
  transaction_id  integer   REFERENCES transactions(id)  ON DELETE RESTRICT NOT NULL,
  account_id      integer   REFERENCES accounts(id)      ON DELETE RESTRICT NOT NULL,
  counterparty_id integer   REFERENCES accounts(id)      ON DELETE RESTRICT NOT NULL,
  direction       direction NOT NULL,
  amount          decimal   NOT NULL,
  PRIMARY KEY(id)
);

ALTER TABLE payments ADD CONSTRAINT valid_amount CHECK (amount > 0);

-- +migrate Down

DROP TABLE IF EXISTS payments;
DROP TYPE IF EXISTS direction;
