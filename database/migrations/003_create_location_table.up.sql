CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE location (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       city_id UUID NOT NULL,
                       FOREIGN KEY (city_id) REFERENCES cities(id)
);
