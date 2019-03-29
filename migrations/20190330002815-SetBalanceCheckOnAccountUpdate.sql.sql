
-- +migrate Up

CREATE CONSTRAINT TRIGGER check_account_update
AFTER UPDATE
ON accounts
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE PROCEDURE check_if_account_balanced();

-- +migrate Down

DROP TRIGGER check_account_update ON accounts;
