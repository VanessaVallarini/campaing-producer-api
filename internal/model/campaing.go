package model

import (
	"time"

	"github.com/google/uuid"
)

type Campaing struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	SlugId      uuid.UUID
	MerchantId  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Active      bool
	Lat         float64
	Long        float64
	Clicks      int
	Impressions int
}

type CampaingListingFilters struct {
	Ids    []uuid.UUID `query:"ids"`
	Page   int         `query:"page"`
	Size   int         `query:"size"`
	Active *bool       `query:"active"`
}

type CampaingCreatingRequest struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
	Ip   string  `json:"ip"`
}
