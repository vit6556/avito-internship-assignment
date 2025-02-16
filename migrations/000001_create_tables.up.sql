CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    balance INTEGER NOT NULL DEFAULT 1000 CHECK (balance >= 0)
);
CREATE INDEX IF NOT EXISTS idx_employees_username ON employees(username);

CREATE TABLE IF NOT EXISTS merch_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0)
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    receiver_id INTEGER NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_transactions_sender ON transactions(sender_id);
CREATE INDEX IF NOT EXISTS idx_transactions_receiver ON transactions(receiver_id);

CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES merch_items(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_purchases_employee ON purchases(employee_id);
CREATE INDEX IF NOT EXISTS idx_purchases_item ON purchases(item_id);