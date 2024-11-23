package internal

// Offer here are the models for insertion in the database
type Offer struct {
	ID                   string `json:"ID"`
	Data                 string `json:"data"`
	MostSpecificRegionID int    `json:"mostSpecificRegionID"`
	StartDate            int64  `json:"startDate"`
	EndDate              int64  `json:"endDate"`
	NumberSeats          int    `json:"numberSeats"`
	Price                int    `json:"price"`
	CarType              string `json:"carType"`
	OnlyVollkasko        bool   `json:"hasVollkasko"`
	FreeKilometers       int    `json:"freeKilometers"`
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
