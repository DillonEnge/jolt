CREATE TABLE listing_images(
    listing_id varchar(255) NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    image_url varchar(255) NOT NULL,
    PRIMARY KEY(listing_id, image_url),
    UNIQUE(listing_id, image_url)
);
---- create above / drop below ----
DROP TABLE listing_images;
