-- name: ListingViewsByID :one
SELECT l.views
FROM listing_views l
WHERE l.listing_id = @listing_id::text;

-- name: UpsertListingViews :one
INSERT INTO listing_views(listing_id) VALUES(
    @listing_id::text
)
ON CONFLICT(listing_id)
DO UPDATE SET
views = listing_views.views+1
RETURNING *;
