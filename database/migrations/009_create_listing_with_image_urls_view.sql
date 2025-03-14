CREATE VIEW listing_with_image_urls AS
SELECT l.*, COALESCE(array_agg(li.image_url) FILTER (WHERE li.image_url IS NOT NULL), ARRAY[]::text[])::text[] AS image_urls
FROM listings l
LEFT JOIN listing_images li ON li.listing_id = l.id
GROUP BY l.id;
---- create above / drop below ----
DROP VIEW listing_with_image_urls;
