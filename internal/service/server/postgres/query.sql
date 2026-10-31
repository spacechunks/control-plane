-- name: CreateServer :exec
INSERT INTO server
    (id, type, name, state, labels, ipv4, port)
VALUES
    ($1, $2, $3, $4, $5, $6, $7);
