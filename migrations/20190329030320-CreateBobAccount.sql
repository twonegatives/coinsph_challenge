
-- +migrate Up
INSERT INTO accounts(name, balance, currency)
VALUES ('bob123', 0, 'usd');

-- +migrate Down

DELETE FROM accounts WHERE name = 'bob123';
