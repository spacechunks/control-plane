-- migrate:up
CREATE TYPE server_state AS ENUM ('PENDING', 'CREATING', 'RUNNING', 'UPDATING', 'DELETING', 'ERROR');
CREATE TYPE server_type AS ENUM ('EPHEMERAL', 'PERSISTENT', 'MINIGAME');

CREATE TABLE cluster (
    id UUID PRIMARY KEY,
    name VARCHAR(256) NOT NULL DEFAULT '',
    slots INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE server (
    id UUID PRIMARY KEY,
    type server_type NOT NULL DEFAULT 'EPHEMERAL',
    name VARCHAR(256) NOT NULL DEFAULT '',
    state server_state NOT NULL DEFAULT 'PENDING',
    labels JSONB NOT NULL DEFAULT '{}',
    ipv4 INET,
    port INTEGER,
    image_url TEXT NOT NULL,
    cluster_id UUID REFERENCES cluster(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- migrate:down

