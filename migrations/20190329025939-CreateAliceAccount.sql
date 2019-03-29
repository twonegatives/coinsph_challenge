
-- +migrate Up
INSERT INTO accounts(name, balance, currency)
VALUES ('alice456', 0, 'usd');

-- +migrate Down

DELETE FROM accounts WHERE name = 'alice456';
