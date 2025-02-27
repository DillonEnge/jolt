-- Write your migrate up statements here

CREATE TABLE negotiations(
    id varchar(255),
    listing_id varchar(255) NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    buyer_email varchar(255) NOT NULL,
    bid int,
    ask int,
    PRIMARY KEY(id),
    UNIQUE(listing_id, buyer_email)
);

---- create above / drop below ----

DROP TABLE negotiations;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
