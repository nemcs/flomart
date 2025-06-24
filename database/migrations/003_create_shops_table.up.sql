CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE shops (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       name TEXT NOT NULL,
                       description TEXT NOT NULL,
                       location_id UUID NOT NULL,
                       owner_id UUID NOT NULL,
                       FOREIGN KEY (location_id) REFERENCES cities(id),
                       FOREIGN KEY (owner_id) REFERENCES users(id)
);
