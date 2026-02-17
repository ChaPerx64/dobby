-- migrate:up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE financial_periods (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    start_dt TIMESTAMPTZ NOT NULL,
    end_dt TIMESTAMPTZ NOT NULL
);

CREATE TABLE envelopes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    CONSTRAINT fk_envelopes_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    financial_period_id UUID NOT NULL,
    user_id UUID NOT NULL,
    envelope_id UUID NOT NULL,
    category VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL,
    description TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_transactions_period FOREIGN KEY (financial_period_id) REFERENCES financial_periods(id),
    CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_transactions_envelope FOREIGN KEY (envelope_id) REFERENCES envelopes(id)
);

-- Indices for performance
CREATE INDEX idx_envelopes_user_id ON envelopes(user_id);
CREATE INDEX idx_transactions_period_id ON transactions(financial_period_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_envelope_id ON transactions(envelope_id);

-- migrate:down

DROP TABLE transactions;
DROP TABLE envelopes;
DROP TABLE financial_periods;
DROP TABLE users;
