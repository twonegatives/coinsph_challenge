
-- +migrate Up

CREATE CONSTRAINT TRIGGER check_insert
AFTER INSERT
ON transactions
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE PROCEDURE check_if_balanced();

-- +migrate Down

DROP TRIGGER check_insert ON transactions;
