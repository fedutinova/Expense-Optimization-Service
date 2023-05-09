CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE wallets
(
    id SERIAL PRIMARY KEY,
    isFamily BOOLEAN DEFAULT FALSE,
    balance NUMERIC(10, 2),
    description TEXT
);

CREATE TABLE users_wallets
(
    is_holder BOOLEAN DEFAULT FALSE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    wallet_id INTEGER REFERENCES wallets(id) ON DELETE CASCADE NOT NULL
 );

CREATE TABLE transactions
(
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER,
    date           DATE NOT NULL DEFAULT CURRENT_DATE,
    category       VARCHAR(50) NOT NULL DEFAULT 'other',
    amount         NUMERIC(10, 2),
    description    VARCHAR(255),
	FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE wallets_transactions
(
    wallet_id INTEGER REFERENCES wallets(id) ON DELETE CASCADE NOT NULL,
    transaction_id INTEGER REFERENCES transactions(id) ON DELETE CASCADE NOT NULL
);