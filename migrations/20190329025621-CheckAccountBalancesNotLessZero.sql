
-- +migrate Up
ALTER TABLE accounts ADD CONSTRAINT valid_balance CHECK (balance >= 0 OR name = 'SYSTEM');

-- +migrate Down

ALTER TABLE accounts DROP CONSTRAINT valid_balance;
