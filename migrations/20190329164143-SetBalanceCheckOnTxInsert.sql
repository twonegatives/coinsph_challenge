
-- +migrate Up

CREATE CONSTRAINT TRIGGER check_tx_insert
AFTER INSERT
ON transactions
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE PROCEDURE check_if_tx_balanced();

-- +migrate Down

DROP TRIGGER check_tx_insert ON transactions;
