-- name: RecordNegotiation :one
INSERT INTO negotiations(id, listing_id, buyer_email)
VALUES(
    uuid_generate_v4(),
    @listing_id::text,
    @buyer_email::text
)
ON CONFLICT(listing_id, buyer_email)
DO NOTHING
RETURNING *;

-- name: NegotiationsByEmail :many
SELECT n.*, l.name, l.seller_email
FROM negotiations n
LEFT JOIN listings l ON l.id = n.listing_id
WHERE l.seller_email = @email::text
OR n.buyer_email = @email::text;

-- name: NegotiationByListingIDAndBuyerEmail :one
SELECT n.*
FROM negotiations n
WHERE n.listing_id = @listing_id::text
AND n.buyer_email = @buyer_email::text;
