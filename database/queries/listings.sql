-- name: ListingByID :one
SELECT l.*
FROM listing_with_image_urls l
WHERE l.id = @listing_id::text;

-- name: ListingsByViews :many
SELECT l.*
FROM listing_with_image_urls l
JOIN listing_views lv ON lv.listing_id = l.id
ORDER BY lv.views DESC
LIMIT $1
OFFSET $2;

-- name: ListingsByLikeName :many
SELECT l.*
FROM listing_with_image_urls l
WHERE UPPER(l.name) LIKE UPPER('%' || @listing_name::text || '%');

-- name: ListingsBySellerEmail :many
SELECT l.*
FROM listing_with_image_urls l
WHERE UPPER(l.seller_email) = UPPER(@seller_email::text);

-- name: RecordListing :one
INSERT INTO listings(id, seller_email, name, description, price) VALUES(
    @id::text,
    @seller_email::text,
    @listing_name::text,
    @description::text,
    @price::int
)
RETURNING *;

-- name: RecordListingImages :many
INSERT INTO listing_images(listing_id, image_url) VALUES(
    @listing_id::text,
    unnest(@image_url_array::text[])
)
RETURNING *;

-- name: DeleteListing :one
DELETE FROM listings l
WHERE l.id = @listing_id::text
RETURNING *;
