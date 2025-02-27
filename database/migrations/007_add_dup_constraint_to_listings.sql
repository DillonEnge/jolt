ALTER TABLE listings
ADD CONSTRAINT listings_unique UNIQUE(name, description);

---- create above / drop below ----

ALTER TABLE listings
DROP CONSTRAINT listings_unique;
