
-- +migrate Up

-- +migrate StatementBegin

CREATE OR REPLACE FUNCTION check_if_account_balanced()
RETURNS TRIGGER
AS $$
DECLARE
  total integer;
BEGIN
  total := (SELECT COALESCE(SUM(CASE WHEN direction = 'outgoing' THEN amount * -1 ELSE amount END), 0) FROM payments WHERE account_id = NEW.id);
  IF (total != NEW.balance) THEN
    RAISE EXCEPTION 'Account balance (%) does not correspond to its payments (%)', NEW.balance, total;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +migrate StatementEnd

-- +migrate Down

DROP FUNCTION IF EXISTS check_if_account_balanced();
