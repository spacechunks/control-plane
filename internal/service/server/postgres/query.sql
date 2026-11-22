--
-- SERVER
--

-- name: CreateServer :exec
INSERT INTO server
    (id, type, name, state, labels, ipv4, port, created_at, updated_at, cluster_id)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: ServerByID :one
SELECT
    sqlc.embed(server),
    sqlc.embed(cluster)
FROM server
    JOIN cluster ON cluster.id = server.cluster_id
WHERE server.id = $1;

-- name: UpdateServer :exec
UPDATE server SET
    name = $1,
    labels = $2,
    ipv4 = $3,
    port = $4,
    updated_at = $5,
    state = $6
WHERE id = $7;

-- name: DeleteServer :exec
DELETE FROM server WHERE id = $1;

--
-- CLUSTER
--

-- name: BestCluster :one
SELECT
    cluster.*,
    COUNT(server.id) AS server_count
FROM cluster
    LEFT JOIN server ON server.cluster_id = cluster.id
WHERE cluster.slots > (SELECT COUNT(*) FROM server WHERE cluster_id = cluster.id)
GROUP BY cluster.id
ORDER BY server_count ASC, random() ASC
LIMIT 1;