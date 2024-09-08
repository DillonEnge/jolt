-- name: ListingByID :many
SELECT l.*
FROM listings l
WHERE l.id = @listing_id::text;

-- name: ListingsByLikeName :many
SELECT l.*
FROM listings l
WHERE UPPER(l.name) LIKE UPPER('%' || @listing_name::text || '%');

-- name: ListingsBySellerEmail :many
SELECT l.*
FROM listings l
WHERE UPPER(l.seller_email) = UPPER(@seller_email::text);

-- name: RecordListing :one
INSERT INTO listings(id, seller_email, name, description, price) VALUES(
    uuid_generate_v4(),
    @seller_email::text,
    @listing_name::text,
    @description::text,
    @price::int
)
RETURNING *;

-- name: DeleteListing :one
DELETE FROM listings l
WHERE l.id = @listing_id::text
RETURNING *;
