# Schema

```sql
CREATE TABLE chains (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    name TEXT NOT NULL,
    chain_id TEXT NOT NULL UNIQUE
);

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    chain_id INT NOT NULL,
    name TEXT NOT NULL,
    denom TEXT NOT NULL,
    decimals INT NOT NULL,

    CONSTRAINT fk_tokens_chain
        FOREIGN KEY (chain_id)
            REFERENCES chains(id)
            ON DELETE CASCADE
);

CREATE UNIQUE INDEX ux_tokens_chain_id_denom ON tokens(chain_id, denom);

CREATE TABLE prices (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    token_id INT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    price TEXT NOT NULL
);

CREATE UNIQUE INDEX ux_prices_token_id_timestamp ON prices(token_id, timestamp);
```

# Queries

## Chains

### Create/Update Chains

```sql
INSERT INTO chains (created_at, updated_at, name, chain_id)
VALUES
    (NOW(), NOW(), 'Chain A', 'chain-a'),
    (NOW(), NOW(), 'Chain B', 'chain-b'),
    (NOW(), NOW(), 'Chain C', 'chain-c')
ON CONFLICT (chain_id)
DO
    UPDATE SET updated_at=NOW(), name=EXCLUDED.name
RETURNING id, created_at, updated_at, name, chain_id;
```

### Get Chains

```sql
SELECT
    id, created_at, updated_at, name, chain_id
FROM chains
WHERE chain_id IN ('chain-a', 'chain-b')
ORDER BY chain_id;
```

## Tokens

### Create/Update Tokens

```sql
INSERT INTO tokens (created_at, updated_at, chain_id, name, denom, decimals)
VALUES
    (NOW(), NOW(), (SELECT id FROM chains WHERE chain_id='chain-a'), 'Token A', 'utokena', 6),
    (NOW(), NOW(), (SELECT id FROM chains WHERE chain_id='chain-b'), 'Token B1', 'utokenb1', 6),
    (NOW(), NOW(), (SELECT id FROM chains WHERE chain_id='chain-b'), 'Token B2', 'utokenb2', 6)
ON CONFLICT (chain_id, denom)
DO UPDATE SET
        updated_at=NOW(),
        chain_id=EXCLUDED.chain_id,
        name=EXCLUDED.name,
        denom=EXCLUDED.denom,
        decimals=EXCLUDED.decimals
RETURNING id, created_at, updated_at, chain_id, name, denom, decimals;
```

### Get Tokens

```sql
SELECT
    tokens.id,
    tokens.created_at,
    tokens.updated_at,
    tokens.chain_id,
    tokens.name,
    tokens.denom,
    tokens.decimals
FROM tokens
JOIN chains ON tokens.chain_id = chains.id
WHERE (chains.chain_id, denom) IN (
    ('chain-a', 'utokena'),
    ('chain-b', 'utokenb2')
ORDER BY (chains.chain_id, denom);
);
```

## Price

### Create/Update Prices

```sql
INSERT INTO prices (created_at, updated_at, token_id, price, timestamp)
VALUES
    (NOW(), NOW(), (
        SELECT tokens.id
            FROM tokens
            JOIN chains ON tokens.chain_id = chains.id
            WHERE chains.chain_id='chain-a' AND tokens.denom='utokena'
    ), '4.20', '2023-06-01 12:00:00 US/Eastern'::TIMESTAMPTZ),
    (NOW(), NOW(), (
        SELECT tokens.id
            FROM tokens
            JOIN chains ON tokens.chain_id = chains.id
            WHERE chains.chain_id='chain-a' AND tokens.denom='utokena'
    ), '4.30', '2023-06-02 12:00:00 US/Eastern'::TIMESTAMPTZ),
    (NOW(), NOW(), (
        SELECT tokens.id
            FROM tokens
            JOIN chains ON tokens.chain_id = chains.id
            WHERE chains.chain_id='chain-a' AND tokens.denom='utokena'
    ), '4.40', '2023-06-03 12:00:00 US/Eastern'::TIMESTAMPTZ)
ON CONFLICT (token_id, timestamp)
DO UPDATE SET
    updated_at=NOW(),
    price=EXCLUDED.price
RETURNING id, created_at, updated_at, token_id, price, timestamp;
```

### Get Prices

```sql
SELECT
    prices.id,
    prices.created_at,
    prices.updated_at,
    prices.token_id,
    prices.price,
    prices.timestamp
FROM prices
JOIN tokens ON prices.token_id=tokens.id
JOIN chains ON tokens.chain_id=chains.id
WHERE (chains.chain_id, tokens.denom) IN (
    ('chain-a', 'utokena')
)
ORDER BY prices.timestamp;
```
