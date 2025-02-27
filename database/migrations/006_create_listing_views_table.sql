-- Write your migrate up statements here

CREATE TABLE listing_views(
    listing_id varchar(255) NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    views int NOT NULL DEFAULT 0,
    PRIMARY KEY(listing_id)
);

---- create above / drop below ----

DROP TABLE listing_views;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
