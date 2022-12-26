CREATE DATABASE ranks;
CREATE DATABASE ranks_test;
CREATE USER usr WITH ENCRYPTED PASSWORD 'pwd';

GRANT ALL PRIVILEGES ON DATABASE ranks TO usr;
GRANT ALL PRIVILEGES ON DATABASE ranks_test TO usr;

\connect ranks;

CREATE TABLE stocks (
    code VARCHAR(7),
    sector VARCHAR(30),
    dividend_yield FLOAT,
    daily_liquidity_in_currency FLOAT,
    last_dividend FLOAT,
    current_price FLOAT,
    pb_ratio FLOAT,
    daily_negotiations INTEGER,
    amount_of_properties INTEGER,
    pb_ratio_ranking INTEGER,
    dy_ranking INTEGER,
    s_ranking INTEGER,
    created_at TIMESTAMP
);
CREATE INDEX idx_created_at_s_rank_ranking ON stocks(created_at, s_ranking);
CREATE INDEX idx_daily_liquidity_in_currency ON stocks(daily_liquidity_in_currency);

CREATE TABLE operations (
    id            SERIAL  PRIMARY KEY,
    code          VARCHAR(7) NOT NULL,
    price         FLOAT      NOT NULL,
    amount        INTEGER    NOT NULL,
    purchase_date TIMESTAMP  NOT NULL
);
CREATE INDEX idx_code ON operations(code);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO usr;

\connect ranks_test;

CREATE TABLE stocks (
    code VARCHAR(7),
    sector VARCHAR(30),
    dividend_yield FLOAT,
    daily_liquidity_in_currency FLOAT,
    last_dividend FLOAT,
    current_price FLOAT,
    pb_ratio FLOAT,
    daily_negotiations INTEGER,
    amount_of_properties INTEGER,
    pb_ratio_ranking INTEGER,
    dy_ranking INTEGER,
    s_ranking INTEGER,
    created_at TIMESTAMP
);
CREATE INDEX idx_created_at_s_rank_ranking ON stocks(created_at, s_ranking);
CREATE INDEX idx_daily_liquidity_in_currency ON stocks(daily_liquidity_in_currency);

CREATE TABLE operations (
    id            SERIAL  PRIMARY KEY,
    code          VARCHAR(7) NOT NULL,
    price         FLOAT      NOT NULL,
    amount        INTEGER    NOT NULL,
    purchase_date TIMESTAMP  NOT NULL
);
CREATE INDEX idx_code ON operations(code);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO usr;
