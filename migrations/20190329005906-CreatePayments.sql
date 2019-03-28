
-- +migrate Up
CREATE TYPE direction AS ENUM('incoming','outgoing');

CREATE TABLE payments (
  id serial,
  transfer_id     integer  REFERENCES transactions(id)  ON DELETE RESTRICT NOT NULL,
  account_id      integer  REFERENCES accounts(id)      ON DELETE RESTRICT NOT NULL,
  participant_id  integer  REFERENCES accounts(id)      ON DELETE RESTRICT NOT NULL,
  direction       direction NOT NULL,
  PRIMARY KEY(id)
);

-- +migrate Down

DROP TABLE IF EXISTS payments;
DROP TYPE IF EXISTS direction;
