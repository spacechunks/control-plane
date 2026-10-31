-- migrate:up
CREATE TYPE server_state AS ENUM ('PENDING', 'CREATING', 'RUNNING', 'UPDATING', 'DELETING', 'ERROR');
CREATE TYPE server_type AS ENUM ('EPHEMERAL', 'PERSISTENT');

CREATE TABLE server (
    id UUID PRIMARY KEY,
    type server_type NOT NULL DEFAULT 'EPHEMERAL',
    name VARCHAR(256) NOT NULL DEFAULT '',
    state server_state NOT NULL DEFAULT 'PENDING',
    labels JSONB NOT NULL DEFAULT '{}',
    ipv4 INET,
    port INTEGER
);

-- migrate:down

