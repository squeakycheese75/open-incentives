CREATE TABLE campaigns (   
    id       BIGSERIAL PRIMARY KEY,
    slug     VARCHAR(16),
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'inactive')),
    rule JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
    );

CREATE INDEX campaigns_status_idx ON campaigns(status);
