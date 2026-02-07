-- migrate:up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL
);

CREATE TABLE financial_periods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255),
    start_dt TIMESTAMP WITH TIME ZONE NOT NULL,
    end_dt TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE envelopes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    financial_period_id UUID NOT NULL REFERENCES financial_periods(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    envelope_id UUID NOT NULL REFERENCES envelopes(id) ON DELETE CASCADE,
    category VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL, -- stored in cents
    description TEXT,
    date TIMESTAMP WITH TIME ZONE NOT NULL
);

-- migrate:down

DROP TABLE transactions;
DROP TABLE envelopes;
DROP TABLE financial_periods;
DROP TABLE users;
