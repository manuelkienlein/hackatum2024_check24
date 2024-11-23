package internal

// here are the models for insertion in the database
type Offer struct {
	Data                 string `json:"data"`
	MostSpecificRegionID int    `json:"most_specific_region_id"`
	StartDate            int64  `json:"start_date"`
	EndDate              int64  `json:"end_date"`
	NumberSeats          int    `json:"number_seats"`
	Price                int    `json:"price"`
	NumberDays           int    `json:"number_days"`
	CarType              string `json:"car_type"`
	OnlyVollkasko        bool   `json:"only_vollkasko"`
	FreeKilometers       int    `json:"free_kilometers"`
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
