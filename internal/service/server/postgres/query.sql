-- name: CreateServer :exec
INSERT INTO server
    (id, type, name, state, labels, ipv4, port, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $7, $8);
