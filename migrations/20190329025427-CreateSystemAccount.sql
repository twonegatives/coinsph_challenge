
-- +migrate Up

INSERT INTO accounts(name, balance, currency)
VALUES ('SYSTEM', 0, 'usd');

-- +migrate Down

DELETE FROM accounts WHERE name = 'SYSTEM';
