-- name: GetCampaign :one
SELECT id, public_id, name, status, rules, created_at, updated_at
FROM campaigns
WHERE public_id = ?
AND org_id = ?
AND project_id = ?
AND deleted_at IS NULL;

-- name: CreateCampaign :one
INSERT INTO campaigns (name, public_id, status, rules, project_id, org_id)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, public_id, name, status, rules, project_id, org_id, created_at, updated_at;

-- name: ListActiveCampaignsByProject :many
SELECT id, public_id, name, status, rules, created_at, updated_at
FROM campaigns
WHERE org_id = ?
AND project_id = ?
AND status = 'active'
AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateCampaign :one
UPDATE campaigns
SET
    public_id = ?,
    name = ?,
    status = ?,
    rules = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
AND deleted_at IS NULL
RETURNING id, public_id, name, status, rules, created_at, updated_at;

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

-- name: CreateProjectAPIKey :one
INSERT INTO api_keys (name, public_id, org_id, project_id, key_hash, prefix, status)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id, public_id, org_id,  project_id, name, key_hash, prefix, status, created_at, updated_at;

-- name: GetActiveAPIKeyByPublicID :one
SELECT
    ak.id,
    ak.public_id,
    ak.key_hash,
    ak.project_id,
    p.org_id
FROM api_keys ak
JOIN projects p ON p.id = ak.project_id
WHERE ak.public_id = ?
  AND ak.revoked_at IS NULL
  AND ak.deleted_at IS NULL
  AND p.deleted_at IS NULL;
  