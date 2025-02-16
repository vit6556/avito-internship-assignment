DROP INDEX IF EXISTS idx_employees_username;
DROP INDEX IF EXISTS idx_transactions_sender;
DROP INDEX IF EXISTS idx_transactions_receiver;
DROP INDEX IF EXISTS idx_purchases_employee;
DROP INDEX IF EXISTS idx_purchases_item;

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS merch_items;
DROP TABLE IF EXISTS employees;