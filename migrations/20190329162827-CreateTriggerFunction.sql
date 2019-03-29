
-- +migrate Up

-- +migrate StatementBegin

CREATE OR REPLACE FUNCTION check_if_balanced()
RETURNS TRIGGER
AS $$
DECLARE
  total integer;
BEGIN
  total := (SELECT SUM(CASE WHEN direction = 'outgoing' THEN amount * -1 ELSE amount END) FROM payments WHERE transaction_id = NEW.id);
  IF (total != 0) THEN
    RAISE EXCEPTION 'Balance of payments does not match for given transaction';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +migrate StatementEnd

-- +migrate Down

DROP FUNCTION IF EXISTS check_if_balanced();
