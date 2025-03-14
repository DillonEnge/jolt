// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: negotiations.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const negotiationByListingIDAndBuyerEmail = `-- name: NegotiationByListingIDAndBuyerEmail :one
SELECT n.id, n.listing_id, n.buyer_email, n.bid, n.ask
FROM negotiations n
WHERE n.listing_id = $1::text
AND n.buyer_email = $2::text
`

type NegotiationByListingIDAndBuyerEmailParams struct {
	ListingID  string `json:"listing_id"`
	BuyerEmail string `json:"buyer_email"`
}

func (q *Queries) NegotiationByListingIDAndBuyerEmail(ctx context.Context, arg NegotiationByListingIDAndBuyerEmailParams) (Negotiation, error) {
	row := q.db.QueryRow(ctx, negotiationByListingIDAndBuyerEmail, arg.ListingID, arg.BuyerEmail)
	var i Negotiation
	err := row.Scan(
		&i.ID,
		&i.ListingID,
		&i.BuyerEmail,
		&i.Bid,
		&i.Ask,
	)
	return i, err
}

const negotiationsByEmail = `-- name: NegotiationsByEmail :many
SELECT n.id, n.listing_id, n.buyer_email, n.bid, n.ask, l.name, l.seller_email
FROM negotiations n
LEFT JOIN listings l ON l.id = n.listing_id
WHERE l.seller_email = $1::text
OR n.buyer_email = $1::text
`

type NegotiationsByEmailRow struct {
	ID          string      `json:"id"`
	ListingID   string      `json:"listing_id"`
	BuyerEmail  string      `json:"buyer_email"`
	Bid         pgtype.Int4 `json:"bid"`
	Ask         pgtype.Int4 `json:"ask"`
	Name        pgtype.Text `json:"name"`
	SellerEmail pgtype.Text `json:"seller_email"`
}

func (q *Queries) NegotiationsByEmail(ctx context.Context, email string) ([]NegotiationsByEmailRow, error) {
	rows, err := q.db.Query(ctx, negotiationsByEmail, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []NegotiationsByEmailRow
	for rows.Next() {
		var i NegotiationsByEmailRow
		if err := rows.Scan(
			&i.ID,
			&i.ListingID,
			&i.BuyerEmail,
			&i.Bid,
			&i.Ask,
			&i.Name,
			&i.SellerEmail,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const recordNegotiation = `-- name: RecordNegotiation :one
INSERT INTO negotiations(id, listing_id, buyer_email)
VALUES(
    uuid_generate_v4(),
    $1::text,
    $2::text
)
ON CONFLICT(listing_id, buyer_email)
DO NOTHING
RETURNING id, listing_id, buyer_email, bid, ask
`

type RecordNegotiationParams struct {
	ListingID  string `json:"listing_id"`
	BuyerEmail string `json:"buyer_email"`
}

func (q *Queries) RecordNegotiation(ctx context.Context, arg RecordNegotiationParams) (Negotiation, error) {
	row := q.db.QueryRow(ctx, recordNegotiation, arg.ListingID, arg.BuyerEmail)
	var i Negotiation
	err := row.Scan(
		&i.ID,
		&i.ListingID,
		&i.BuyerEmail,
		&i.Bid,
		&i.Ask,
	)
	return i, err
}
