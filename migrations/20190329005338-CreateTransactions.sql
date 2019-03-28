
-- +migrate Up
CREATE TABLE transactions (
  id serial,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY(id)
);

-- +migrate Down

DROP TABLE IF EXISTS transactions;
