--
-- SERVER
--

-- name: CreateServer :exec
INSERT INTO server
    (id, type, name, state, labels, ipv4, port, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: ServerByID :one
SELECT * FROM server WHERE id = $1;

-- name: UpdateServer :exec
UPDATE server SET
    name = $1,
    labels = $2,
    ipv4 = $3,
    port = $4,
    updated_at = $5
WHERE id = $6;

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