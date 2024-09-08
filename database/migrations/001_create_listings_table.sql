-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE listings(
    id varchar(255),
    name varchar(255) NOT NULL,
    description varchar(255),
    price int NOT NULL,
    PRIMARY KEY(id)
);

---- create above / drop below ----

DROP TABLE listings;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
