package models

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
