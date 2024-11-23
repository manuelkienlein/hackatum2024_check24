package internal

// here are the models for insertion in the database
type Offer struct {
	RegionID              int    `json:"regionID"`
	TimeRangeStart        int    `json:"timeRangeStart"`
	TimeRangeEnd          int    `json:"timeRangeEnd"`
	NumberDays            int    `json:"numberDays"`
	SortOrder             string `json:"sortOrder"`
	Page                  int    `json:"page"`
	PageSize              int    `json:"pageSize"`
	PriceRangeWidth       int    `json:"priceRangeWidth"`
	MinFreeKilometerWidth int    `json:"minFreeKilometerWidth"`
	MinNumberSeats        int    `json:"minNumberSeats"`
	MinPrice              int    `json:"minPrice"`
	MaxPrice              int    `json:"maxPrice"`
	CarType               string `json:"carType"`
	OnlyVollkasko         bool   `json:"onlyVollkasko"`
	MinFreeKilometer      int    `json:"minFreeKilometer"`
}

// here are the models for the response
type ResponseOffer struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type PriceRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Count int `json:"count"`
}

type CarTypeCounts struct {
	Small  int `json:"small"`
	Sports int `json:"sports"`
	Luxury int `json:"luxury"`
	Family int `json:"family"`
}

type SeatsCount struct {
	NumberSeats int `json:"numberSeats"`
	Count       int `json:"count"`
}

type FreeKilometerRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Count int `json:"count"`
}

type VollkaskoCount struct {
	TrueCount  int `json:"trueCount"`
	FalseCount int `json:"falseCount"`
}

type OfferFilterParams struct {
	RegionID              int
	TimeRangeStart        int
	TimeRangeEnd          int
	NumberDays            int
	SortOrder             string
	Page                  int
	PageSize              int
	PriceRangeWidth       int
	MinFreeKilometerWidth int
	MinNumberSeats        int
	MinPrice              int
	MaxPrice              int
	CarType               string
	OnlyVollkasko         bool
	MinFreeKilometer      int
}
