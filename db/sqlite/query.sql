-- name: GetCampaign :one
SELECT id, public_id, name, status, rule, created_at, updated_at
FROM campaigns
WHERE public_id = ?
AND deleted_at IS NULL;

-- name: CreateCampaign :one
INSERT INTO campaigns (name, public_id, status, rule, project_id, org_id)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, public_id, name, status, rule, project_id, org_id, created_at, updated_at;

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

-- name: GetOrgByPublicID :one
SELECT id, public_id, name
FROM organizations
WHERE public_id = ?
AND deleted_at IS NULL;

-- name: GetUserByEmailAndOrg :one
SELECT id, public_id, email, org_id, password_hash, role
FROM users
WHERE email = ?
AND org_id = ?
AND deleted_at IS NULL;

-- name: GetProjectByPublicID :one
SELECT id, public_id, name
FROM projects
WHERE public_id = ?
AND org_id = ?
AND deleted_at IS NULL;
