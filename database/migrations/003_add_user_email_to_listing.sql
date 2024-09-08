-- Write your migrate up statements here

ALTER TABLE listings
ADD COLUMN seller_email varchar(255) NOT NULL;

---- create above / drop below ----

ALTER TABLE listings
DROP COLUMN seller_email;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
