ALTER TABLE transactions
DROP CONSTRAINT IF EXISTS transactions_type_check;

ALTER TABLE transactions
ADD CONSTRAINT transactions_type_check
CHECK (type IN ('income', 'expense', 'transfer', 'informational'));

ALTER TABLE scheduled_movements
DROP CONSTRAINT IF EXISTS scheduled_movements_type_check;

ALTER TABLE scheduled_movements
ADD CONSTRAINT scheduled_movements_type_check
CHECK (type IN ('income', 'expense', 'informational'));

ALTER TABLE scheduled_movements
DROP CONSTRAINT IF EXISTS scheduled_movements_amount_check;

ALTER TABLE scheduled_movements
ADD CONSTRAINT scheduled_movements_amount_check
CHECK (
    (type = 'informational' AND amount = 0)
    OR (type <> 'informational' AND amount > 0)
);
