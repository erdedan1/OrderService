CREATE TABLE IF NOT EXISTS orders
(
    id UUID NOT NULL PRIMARY KEY
    user_id UUID NOT NULL
    market_id UUID NOT NULL
    quantity    int NOT NULL
    type    TEXT NOT NULL
    status  TEXT NOT NULL
    price   TEXT NOT NULL
    created_at  timestamptz      DEFAULT NOW(),
    updated_at   timestamptz      DEFAULT NULL
    deleted_at  timestamptz      DEFAULT NULL
)