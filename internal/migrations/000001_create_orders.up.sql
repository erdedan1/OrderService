CREATE TABLE IF NOT EXISTS orders
(
    id UUID NOT NULL PRIMARY KEY
    user_id UUID NOT NULL
    market_id UUID NOT NULL
    quantity    int NOT NULL
    order_type    VARCHAR(16) NOT NULL
    order_status  VARCHAR(32) NOT NULL
    price   NUMERIC NOT NULL
    created_at  timestamptz      DEFAULT NOW(),
    updated_at   timestamptz      DEFAULT NULL
    deleted_at  timestamptz      DEFAULT NULL
);