-- name: GetCampaign :one
SELECT id, slug, name, status, rule, created_at, updated_at
FROM campaigns
WHERE name = $1
AND deleted_at IS NULL;

-- name: CreateCampaign :one
INSERT INTO campaigns (name, slug, status, rule)
VALUES ($1, $2, $3, $4)
RETURNING id, name, slug, status, rule, created_at, updated_at;

-- -- name: ListCampaigns :many
-- SELECT id, name, status, rule, created_at, updated_at
-- FROM campaigns
-- ORDER BY created_at DESC;

-- name: ListActiveByProject :many
SELECT id, name, status, rule, created_at, updated_at
FROM campaigns
WHERE status = 'active'
AND ord_id = ?
AND project_id = ?
ORDER BY created_at DESC;

-- name: UpdateCampaign :one
UPDATE campaigns
SET
    name = $2,
    status = $3,
    rule = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, name, status, rule, created_at, updated_at;

-- name: DeleteCampaign :exec
DELETE FROM campaigns
WHERE id = $1;

