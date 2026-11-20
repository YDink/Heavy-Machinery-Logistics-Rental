package domain

import "github.com/google/uuid"

// MatchedCarrier represents a dimensionally compatible lowboy transport fleet candidate.
type MatchedCarrier struct {
	FleetID                 uuid.UUID
	OwnerID                 uuid.UUID
	Make                    string
	Model                   string
	OwnerTradingName        string
	OwnerPhoneNumber        string
	PayloadCapacityKg       float64
	DeckWellLengthMm        int
	DeckWidthMm             int
	OutriggersExtensionMm   int
	DeckHeightGroundMm      int
	TareWeightKg            float64
	AxleCount               int
	DailyRateMinor          int64
	Currency                string
	DistanceFromDepotMeters float64
	AverageRating           float64
}
