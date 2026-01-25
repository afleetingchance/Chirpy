-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT *
FROM chirps
WHERE (@user_id::uuid = '00000000-0000-0000-0000-000000000000'::UUID OR user_id = @user_id::uuid)
ORDER BY 
    CASE WHEN @sort::text = 'created_at_asc' THEN created_at END ASC,
    CASE WHEN @sort::text = 'created_at_desc' THEN created_at END DESC;

-- name: GetChirpById :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;