-- name: GetCampaign :one
SELECT id, public_id, name, status, rule, created_at, updated_at
FROM campaigns
WHERE public_id = ?
AND deleted_at IS NULL;

-- name: CreateCampaign :one
INSERT INTO campaigns (name, public_id, status, rule)
VALUES (?, ?, ?, ?)
RETURNING id, public_id, name, status, rule, created_at, updated_at;

-- name: ListCampaigns :many
SELECT id, public_id, name, status, rule, created_at, updated_at
FROM campaigns
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListActiveCampaigns :many
SELECT id, public_id, name, status, rule, created_at, updated_at
FROM campaigns
WHERE status = 'active'
AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateCampaign :one
UPDATE campaigns
SET
    public_id = ?,
    name = ?,
    status = ?,
    rule = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
AND deleted_at IS NULL
RETURNING id, public_id, name, status, rule, created_at, updated_at;

-- name: DeleteCampaign :exec
UPDATE campaigns
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ?;
