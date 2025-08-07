-- name: CreateToken :exec
INSERT INTO refresh_tokens(token,created_at,updated_at,user_id,expires_at,revoked_at)
VALUES($1, $2, $3, $4, $5, $6);

-- name: GetTokenByID :one
SELECT * FROM refresh_tokens 
WHERE token = $1 AND revoked_at IS NULL;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = now(), updated_at = now()
WHERE token = $1; 
