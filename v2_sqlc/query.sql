-- name: CreateChain :one
INSERT INTO chains (created_at, updated_at, name, chain_id)
VALUES
    (NOW(), NOW(), @name::TEXT, @chain_id::TEXT)
ON CONFLICT (chain_id)
DO
    UPDATE SET updated_at=NOW(), name=EXCLUDED.name
RETURNING *;

-- name: GetChains :one
SELECT *
FROM chains
WHERE chain_id = @chain_id::TEXT
ORDER BY chain_id;

-- name: CreateToken :one
INSERT INTO tokens (created_at, updated_at, chain_id, name, denom, decimals)
VALUES (
    NOW(),
    NOW(),
    (SELECT id FROM chains WHERE chain_id = @chain_id::TEXT),
    @name::TEXT,
    @denom::TEXT,
    @decimals::INT
)
ON CONFLICT (chain_id, denom)
DO UPDATE SET
        updated_at=NOW(),
        chain_id=EXCLUDED.chain_id,
        name=EXCLUDED.name,
        denom=EXCLUDED.denom,
        decimals=EXCLUDED.decimals
RETURNING *;

-- name: GetToken :one
SELECT tokens.*
FROM tokens
JOIN chains ON tokens.chain_id = chains.id
WHERE chains.chain_id = @chain_id::TEXT
  AND tokens.denom = @denom::TEXT
ORDER BY (chains.chain_id, denom);

-- name: CreatePrice :one
INSERT INTO prices (created_at, updated_at, token_id, price, timestamp)
VALUES
    (NOW(), NOW(), (
        SELECT tokens.id
            FROM tokens
            JOIN chains ON tokens.chain_id = chains.id
            WHERE chains.chain_id = @chain_id::TEXT AND tokens.denom = @denom::TEXT
    ), @price::TEXT, @timestamp::TIMESTAMPTZ)
ON CONFLICT (token_id, timestamp)
DO UPDATE SET
    updated_at=NOW(),
    price=EXCLUDED.price
RETURNING *;

-- name: GetPrices :many
SELECT prices.*
FROM prices
JOIN tokens ON prices.token_id=tokens.id
JOIN chains ON tokens.chain_id=chains.id
WHERE chains.chain_id = @chain_id::TEXT
  AND tokens.denom = @denom::TEXT
ORDER BY prices.timestamp;
