package database

import "context"

type NegotiationQuerier interface {
	RecordNegotiation(ctx context.Context, arg RecordNegotiationParams) (Negotiation, error)
	NegotiationsByListingIDAndBuyerEmail(ctx context.Context, arg NegotiationByListingIDAndBuyerEmailParams) (Negotiation, error)
}

func GetOrCreateNegotiation(ctx context.Context, db NegotiationQuerier, params RecordNegotiationParams) {
	// negotiation, err := db.NegotiationsByListingIDAndBuyerEmail()
}
