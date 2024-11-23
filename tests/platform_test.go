package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/internal/models"
	"testing"
)

func TestPostOffersPlatform(t *testing.T) {

	// Setup der App
	app := setupApp()

	// JSON-Daten als String
	offerJSON := `{"offers":[{"ID":"87b57605-1ed2-43be-9613-e279d446466c","carType":"family","data":"LeMxLnrv9bMYI0iSDjUn3DCHo1y/SDeAC4ZFHUDO41nQHNxUR5nOgQx3db9Pt8TESK50BEPZzzu4tcVg2qujF3aN0oMr5Bhmc19vgnu533HzIYlmE7454fBL2ercKABOrO3B1ntpSCpAa4zl2H3QEwWvRePE85hzof9HMnngpv1/9abUOzWutvZNZtUae9XFHEPc1sf6+GAESw6HxKbHB2LlG40bM9+jZujlfB535q1UgIVG8S25zG49k6+IB+Lc3enyXuL6F+acT+przcsvcMzgefPXujERGprnqHfCfdnKWg3mRe9bDtgrqT/4Oaw+Cev0+yMgY58WB5yCPCbP6Q==","endDate":1673568000000,"freeKilometers":707,"hasVollkasko":true,"mostSpecificRegionID":118,"numberSeats":2,"price":3796,"startDate":1673395200000},{"ID":"03cfd034-ab0f-43bf-ba74-403c675a2752","carType":"family","data":"7GfdC6mskC0+sTOvK3KbNXRywM8DPu+nLSn8pvGQexaDH/bPpJXloGHT7gGpsnHXK/3QPXim4RzY+DvRJKgxABpQq9SVD/pLhF4WmZC0Z5mY9l7abN8c28lNKV9IGn6Ngx+T0qG+PY2/ug98EtcAr6B0jQ11+mXobbk7sW9lYxkEqUjAdNDoANJ1MpOjW+Wm+tYWs/+qwHiUKhdUivgeYWfHBxs8gNDurHhzGKquyHUs+c1WZd3LLWqk7zXbwIhqokAVaFJ9WC5wRJRvSOdzldLYmGxNUmX1cGntGNytpvtRnbfP4MoAlLiLhzDrIqm7E5hVqCcU5KsKcZlf2iymyw==","endDate":1673568000000,"freeKilometers":465,"hasVollkasko":false,"mostSpecificRegionID":110,"numberSeats":3,"price":5636,"startDate":1673308800000},{"ID":"56d019d0-0ec8-4cc0-bbb4-1930ab03b05e","carType":"small","data":"JsPIwtHF+AD00ryVxDfXQRDq36j+KXpaWRxYEFSelbyWDTwa3QOtz7jYN1j80IZa1QePIhySW50zUO8E6hLWT9/k/C0ePALfuiq0L9f1j1v+VYQVL2H1fDMpW3dqzkz4tEnaLYHjtAP1rmlKdzMACILdqdlGygAcw2s8VSQucvvENyoTEFFb3On3oNq3K31yE2Xbb1Jc0OaG7rBeau7Nc7+wDf8ZhackEEg4171LZ/dW9+0ncnE3rtZ+RQK530NgOTTCCEWCQZzvvl7lDqeWdX8s/NqnAwkEcGISMR809h9hK96aXpLfdcPkyTLi5LCEZNyjNOaqm3BPbk7Rp9u9NQ==","endDate":1673568000000,"freeKilometers":840,"hasVollkasko":true,"mostSpecificRegionID":122,"numberSeats":6,"price":4192,"startDate":1673222400000},{"ID":"a6b115b8-8115-4ca0-b340-108f242cb2d0","carType":"family","data":"/RI57PSZIlftAWS98kz86tD4uZkRT89t24cDFSj7Gy8kymkGpjczlqv36D4C1UfcW32PoO+D49sbEC9daWGt5TvmRLcUO2YqzTi+S2gBxkjtL3ytRLlkPIm65T6pu3A+eA8UqxSl++TRvCJ9AP4XXnkPFIVdSzoK9aoOPjOIDbPh9+1MMEo7EEOTUKGzgt3hP9k49KDe2itAPE9e0lpb8XlaLwJdcBhdo7fWep15EufWWPUI2UefQA9QOVXMScymdjd5JoGcJ+RcLkJ2r8y4tcHY7xKfeb1nZImGOcpSY7NKtc9TOsEowZZptU5F8a6kBZucF5mlTzrSbV6bd+hetw==","endDate":1673481600000,"freeKilometers":987,"hasVollkasko":false,"mostSpecificRegionID":101,"numberSeats":2,"price":7829,"startDate":1673308800000},{"ID":"8a020d24-9508-4192-907e-3ec8f1579777","carType":"sports","data":"OcsZ4lCCLyVzE3oMgHBsgo0osk/dDNCnydCXUrlkEManJvwV2h7cGvnNn0DVYWUsWlLeNwPSe4XC06UVeJgmONdvbfqhDzMKspoKYNCqptLDUQm45PQ88a9sZ7CjX6goiRgA7qElIu6NCgoHiPgCjOV5fk+hZLxmCOF2ldM8xvlcY+Zt1ae8YOJpPmY+cpYppedc+uYCqKY+8EYkK82tAfa8AjvNj23hQEPkcWTzK7vQU1HZXmOP7BGBblkBlZWL3tuZrdFrpgoGYo5XGi9Q1DeZnVCd6jKV4lO8KvWmGOIkHnMKYeoq8Hzr8toGFksCEZAe2P5I2yklMqICR4M77Q==","endDate":1673568000000,"freeKilometers":552,"hasVollkasko":false,"mostSpecificRegionID":88,"numberSeats":3,"price":13676,"startDate":1673222400000},{"ID":"2f2942e4-7973-4b41-90ec-bd961ff721dd","carType":"small","data":"Al4gVqRkaInqG2Uiq5OaD2ESPNfr/AobonuBzWnw+s5X+5LDMby48+6Qs7Fm9zAA4qsORJWhz3TweNenzlvyozbZUGHepMM3SLRxUYTT1t4gDNn+SAkQ1BBpEsJC/5t3S7I/5syS3lMKsUh/Nn01rU0y/VU7zN1WSxvz0sz4yudiW43V+ok7VlXRywGeJUtDJeK4MdnqJtM/ITPi4hR1JyIudswNYEvo/SBE2Ey/6X3kVMGKxwQkI0P1WmcReDLvmTey3eHTtUiGmyKk9Mi+U/20acbPg0eXwK03gGabrs0Sh9I8GKEatUEL107CLwsBkLOfb/RaVACnf0cko/WRmQ==","endDate":1673568000000,"freeKilometers":889,"hasVollkasko":true,"mostSpecificRegionID":90,"numberSeats":6,"price":3949,"startDate":1673308800000},{"ID":"7dd014c0-f473-4210-ad1c-6aa2aca1e6ca","carType":"small","data":"ycmO4hCt6DId5TQoahhsXU220wwGUxzy5qMl4WziRLa1rFgSd2dS5a0zWEgdmYGX6gEAH48AlhCzUS1609wFoPT4sb3pJCEXJm1lQ7k44iI1XUDTKYuLfShxbwyumcfjYgThp3/yQEMdis+PtOr5RYjdPI258cH+VgEuM7+TLaf2mwVROs+c9A/1pFRd8PI2lq8JE7n8Gwr6w6luLseWw/cdYjhO0W/xR0R+KJJwl+n7UN2Y7rwEYNHvHFE1Bj9uPTYr71pH05e4VUmBJamRh3mzNmN+h+wgp9ndG1NdX0o4VvssWekwrurthW6aZi5zIQjlREzYrHX8s9eGJ+gH2w==","endDate":1673481600000,"freeKilometers":895,"hasVollkasko":false,"mostSpecificRegionID":120,"numberSeats":3,"price":2077,"startDate":1673222400000},{"ID":"296757ff-e327-4f77-9703-932b2386d0bc","carType":"luxury","data":"5x9HqvGQIytOxK08ZdzsOmzhmwnVY2BEdF8E1XApBv21YNx6MCLRNHSggd/P9B/0t2EakUGhchsSXlplGSkitleKLFK6rW/onnNGJw9URdTCrOB77PdX92pMg7PAGmJsr01jfU+EqbBOitzbRyaoWybTrcYt1k8uqQ7mn0YQLPEaPw6WRmLkMOiZv2RkkA4Ufw5jazKVK/dbGpggrWK8l/60YKUAq/V6VUpbGa2JB0JMDsGBQG68DrhGizA49+o5ngB0ikkxbq79i7WamR7DuGVsaKJEGl4jvoQIabIHgyE6UzuWL0tMLjMnpbOJrvqUp7FA5ugqVCXuIpLQEN2PNw==","endDate":1673568000000,"freeKilometers":586,"hasVollkasko":true,"mostSpecificRegionID":73,"numberSeats":2,"price":5733,"startDate":1673308800000},{"ID":"1f350c5b-6fcd-4f78-99b8-2b51ee632abe","carType":"sports","data":"vPwqjsM2rE8JfcEjreTio3QsutLp6cJRb3bbeaACvLhtMdP53GtvunrtrI2JTVLXk1LZvUn6pcVZ+/XDLKnCxD3e3vt5peUoU6U5Dxk8+9vZ6icW4d2dcXnbmpjAp7TXdHHqTdQ0mxU7kOlx7DfTUphQm5cs1R3v+9zsPvRTHAzfJ3J7E2Uv5dmJvYYKR5Voq6CIq49PaP8Yktpwh0OfQNAetQxPDJmFa4WOKOnNobYBRctfmnw8hV0A1jH05wDdsoNxYE8TDmRXEe4yVjCK9HSM7euZEUdgzTDH6t7G0e9k/YipjxTaH8JC/tbgz7Dd2pRVU9AMBZgKum7I3qDUjw==","endDate":1673568000000,"freeKilometers":544,"hasVollkasko":false,"mostSpecificRegionID":109,"numberSeats":5,"price":1727,"startDate":1673395200000},{"ID":"b412549a-69d5-4774-ae19-5819366cacac","carType":"family","data":"6Be14/fRDD55mVIHxzEDIEVpUndFp2VFa8CTn6tleexSZsFMu5adIEP0dmugMASXW7vzU2Q7AuTtD0GDM/oruXxmzdmVh649iif7mAER6ZDJGT/YJLmq0dGvVTFpKI/dUPj51uwQ/YAgQpSWNfBzIDxXOtEE18aAKeTCcr7tdWOIm/f/WthufXs5w6/mr/Cge4tTI0dshSBfsTCVfvd2fCSRX0IQV0facPLCVMGgb26KNY6YH0Oj/fXcUUXElGohSbCwmGRro85jqWPO5aGbIvVkEJGXNiB1posH4D+5hELtt/nGArA2P+ljw/DFC65dMOJisxNc0evde5u2c3Ispw==","endDate":1673481600000,"freeKilometers":255,"hasVollkasko":true,"mostSpecificRegionID":109,"numberSeats":6,"price":9346,"startDate":1673308800000}]}`

	// Führe eine POST-Anfrage an /api/offers durch
	req := httptest.NewRequest("POST", "/api/offers", bytes.NewReader([]byte(offerJSON)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Überprüfen, ob keine Fehler beim Testen aufgetreten sind
	assert.NoError(t, err)

	// Überprüfen, ob der Statuscode = 200 ist
	assert.Equal(t, resp.StatusCode, 200)
}

// This test should return all ten offers
func TestGetOffersPlattformBase(t *testing.T) {
	app := setupApp()
	TestPostOffersPlatform(t)

	req := httptest.NewRequest(http.MethodGet, "/api/offers?regionID=0&timeRangeEnd=1673568000000&timeRangeStart=0&numberDays=3&sortOrder=price-asc&page=0&pageSize=100&priceRangeWidth=20000&minFreeKilometerWidth=1000", nil)
	resp, _ := app.Test(req)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error during GET request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be 200")

	var actualResponse models.OfferQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	expectedResponse := models.OfferQueryResponse{
		Offers: []models.ResponseOffer{
			{
				ID:   "1f350c5b-6fcd-4f78-99b8-2b51ee632abe",
				Data: "vPwqjsM2rE8JfcEjreTio3QsutLp6cJRb3bbeaACvLhtMdP53GtvunrtrI2JTVLXk1LZvUn6pcVZ+/XDLKnCxD3e3vt5peUoU6U5Dxk8+9vZ6icW4d2dcXnbmpjAp7TXdHHqTdQ0mxU7kOlx7DfTUphQm5cs1R3v+9zsPvRTHAzfJ3J7E2Uv5dmJvYYKR5Voq6CIq49PaP8Yktpwh0OfQNAetQxPDJmFa4WOKOnNobYBRctfmnw8hV0A1jH05wDdsoNxYE8TDmRXEe4yVjCK9HSM7euZEUdgzTDH6t7G0e9k/YipjxTaH8JC/tbgz7Dd2pRVU9AMBZgKum7I3qDUjw==",
			}, {
				ID:   "7dd014c0-f473-4210-ad1c-6aa2aca1e6ca",
				Data: "ycmO4hCt6DId5TQoahhsXU220wwGUxzy5qMl4WziRLa1rFgSd2dS5a0zWEgdmYGX6gEAH48AlhCzUS1609wFoPT4sb3pJCEXJm1lQ7k44iI1XUDTKYuLfShxbwyumcfjYgThp3/yQEMdis+PtOr5RYjdPI258cH+VgEuM7+TLaf2mwVROs+c9A/1pFRd8PI2lq8JE7n8Gwr6w6luLseWw/cdYjhO0W/xR0R+KJJwl+n7UN2Y7rwEYNHvHFE1Bj9uPTYr71pH05e4VUmBJamRh3mzNmN+h+wgp9ndG1NdX0o4VvssWekwrurthW6aZi5zIQjlREzYrHX8s9eGJ+gH2w==",
			}, {
				ID:   "87b57605-1ed2-43be-9613-e279d446466c",
				Data: "LeMxLnrv9bMYI0iSDjUn3DCHo1y/SDeAC4ZFHUDO41nQHNxUR5nOgQx3db9Pt8TESK50BEPZzzu4tcVg2qujF3aN0oMr5Bhmc19vgnu533HzIYlmE7454fBL2ercKABOrO3B1ntpSCpAa4zl2H3QEwWvRePE85hzof9HMnngpv1/9abUOzWutvZNZtUae9XFHEPc1sf6+GAESw6HxKbHB2LlG40bM9+jZujlfB535q1UgIVG8S25zG49k6+IB+Lc3enyXuL6F+acT+przcsvcMzgefPXujERGprnqHfCfdnKWg3mRe9bDtgrqT/4Oaw+Cev0+yMgY58WB5yCPCbP6Q==",
			}, {
				ID:   "2f2942e4-7973-4b41-90ec-bd961ff721dd",
				Data: "Al4gVqRkaInqG2Uiq5OaD2ESPNfr/AobonuBzWnw+s5X+5LDMby48+6Qs7Fm9zAA4qsORJWhz3TweNenzlvyozbZUGHepMM3SLRxUYTT1t4gDNn+SAkQ1BBpEsJC/5t3S7I/5syS3lMKsUh/Nn01rU0y/VU7zN1WSxvz0sz4yudiW43V+ok7VlXRywGeJUtDJeK4MdnqJtM/ITPi4hR1JyIudswNYEvo/SBE2Ey/6X3kVMGKxwQkI0P1WmcReDLvmTey3eHTtUiGmyKk9Mi+U/20acbPg0eXwK03gGabrs0Sh9I8GKEatUEL107CLwsBkLOfb/RaVACnf0cko/WRmQ==",
			}, {
				ID:   "56d019d0-0ec8-4cc0-bbb4-1930ab03b05e",
				Data: "JsPIwtHF+AD00ryVxDfXQRDq36j+KXpaWRxYEFSelbyWDTwa3QOtz7jYN1j80IZa1QePIhySW50zUO8E6hLWT9/k/C0ePALfuiq0L9f1j1v+VYQVL2H1fDMpW3dqzkz4tEnaLYHjtAP1rmlKdzMACILdqdlGygAcw2s8VSQucvvENyoTEFFb3On3oNq3K31yE2Xbb1Jc0OaG7rBeau7Nc7+wDf8ZhackEEg4171LZ/dW9+0ncnE3rtZ+RQK530NgOTTCCEWCQZzvvl7lDqeWdX8s/NqnAwkEcGISMR809h9hK96aXpLfdcPkyTLi5LCEZNyjNOaqm3BPbk7Rp9u9NQ==",
			}, {
				ID:   "03cfd034-ab0f-43bf-ba74-403c675a2752",
				Data: "7GfdC6mskC0+sTOvK3KbNXRywM8DPu+nLSn8pvGQexaDH/bPpJXloGHT7gGpsnHXK/3QPXim4RzY+DvRJKgxABpQq9SVD/pLhF4WmZC0Z5mY9l7abN8c28lNKV9IGn6Ngx+T0qG+PY2/ug98EtcAr6B0jQ11+mXobbk7sW9lYxkEqUjAdNDoANJ1MpOjW+Wm+tYWs/+qwHiUKhdUivgeYWfHBxs8gNDurHhzGKquyHUs+c1WZd3LLWqk7zXbwIhqokAVaFJ9WC5wRJRvSOdzldLYmGxNUmX1cGntGNytpvtRnbfP4MoAlLiLhzDrIqm7E5hVqCcU5KsKcZlf2iymyw==",
			}, {
				ID:   "296757ff-e327-4f77-9703-932b2386d0bc",
				Data: "5x9HqvGQIytOxK08ZdzsOmzhmwnVY2BEdF8E1XApBv21YNx6MCLRNHSggd/P9B/0t2EakUGhchsSXlplGSkitleKLFK6rW/onnNGJw9URdTCrOB77PdX92pMg7PAGmJsr01jfU+EqbBOitzbRyaoWybTrcYt1k8uqQ7mn0YQLPEaPw6WRmLkMOiZv2RkkA4Ufw5jazKVK/dbGpggrWK8l/60YKUAq/V6VUpbGa2JB0JMDsGBQG68DrhGizA49+o5ngB0ikkxbq79i7WamR7DuGVsaKJEGl4jvoQIabIHgyE6UzuWL0tMLjMnpbOJrvqUp7FA5ugqVCXuIpLQEN2PNw==",
			}, {
				ID:   "a6b115b8-8115-4ca0-b340-108f242cb2d0",
				Data: "/RI57PSZIlftAWS98kz86tD4uZkRT89t24cDFSj7Gy8kymkGpjczlqv36D4C1UfcW32PoO+D49sbEC9daWGt5TvmRLcUO2YqzTi+S2gBxkjtL3ytRLlkPIm65T6pu3A+eA8UqxSl++TRvCJ9AP4XXnkPFIVdSzoK9aoOPjOIDbPh9+1MMEo7EEOTUKGzgt3hP9k49KDe2itAPE9e0lpb8XlaLwJdcBhdo7fWep15EufWWPUI2UefQA9QOVXMScymdjd5JoGcJ+RcLkJ2r8y4tcHY7xKfeb1nZImGOcpSY7NKtc9TOsEowZZptU5F8a6kBZucF5mlTzrSbV6bd+hetw==",
			}, {
				ID:   "b412549a-69d5-4774-ae19-5819366cacac",
				Data: "6Be14/fRDD55mVIHxzEDIEVpUndFp2VFa8CTn6tleexSZsFMu5adIEP0dmugMASXW7vzU2Q7AuTtD0GDM/oruXxmzdmVh649iif7mAER6ZDJGT/YJLmq0dGvVTFpKI/dUPj51uwQ/YAgQpSWNfBzIDxXOtEE18aAKeTCcr7tdWOIm/f/WthufXs5w6/mr/Cge4tTI0dshSBfsTCVfvd2fCSRX0IQV0facPLCVMGgb26KNY6YH0Oj/fXcUUXElGohSbCwmGRro85jqWPO5aGbIvVkEJGXNiB1posH4D+5hELtt/nGArA2P+ljw/DFC65dMOJisxNc0evde5u2c3Ispw==",
			}, {
				ID:   "8a020d24-9508-4192-907e-3ec8f1579777",
				Data: "OcsZ4lCCLyVzE3oMgHBsgo0osk/dDNCnydCXUrlkEManJvwV2h7cGvnNn0DVYWUsWlLeNwPSe4XC06UVeJgmONdvbfqhDzMKspoKYNCqptLDUQm45PQ88a9sZ7CjX6goiRgA7qElIu6NCgoHiPgCjOV5fk+hZLxmCOF2ldM8xvlcY+Zt1ae8YOJpPmY+cpYppedc+uYCqKY+8EYkK82tAfa8AjvNj23hQEPkcWTzK7vQU1HZXmOP7BGBblkBlZWL3tuZrdFrpgoGYo5XGi9Q1DeZnVCd6jKV4lO8KvWmGOIkHnMKYeoq8Hzr8toGFksCEZAe2P5I2yklMqICR4M77Q==",
			},
		},
		PriceRanges: []models.PriceRange{
			{Start: 0, End: 20000, Count: 10},
		},
		CarTypeCounts: models.CarTypeCounts{
			Small:  3,
			Sports: 2,
			Luxury: 1,
			Family: 4,
		},
		SeatsCount: []models.SeatsCount{
			{NumberSeats: 5, Count: 1},
			{NumberSeats: 3, Count: 3},
			{NumberSeats: 2, Count: 3},
			{NumberSeats: 6, Count: 3},
		},
		FreeKilometerRange: []models.FreeKilometerRange{
			{Start: 0, End: 1000, Count: 10},
		},
		VollkaskoCount: models.VollkaskoCount{
			TrueCount:  5,
			FalseCount: 5,
		},
	}

	// Assertions for the response structure
	assert.Equal(t, expectedResponse.Offers, actualResponse.Offers, "Offers do not match")
	assert.Equal(t, expectedResponse.PriceRanges, actualResponse.PriceRanges, "PriceRanges do not match")
	assert.Equal(t, expectedResponse.CarTypeCounts, actualResponse.CarTypeCounts, "CarTypeCounts do not match")
	assert.Equal(t, expectedResponse.SeatsCount, actualResponse.SeatsCount, "SeatsCount do not match")
	assert.Equal(t, expectedResponse.FreeKilometerRange, actualResponse.FreeKilometerRange, "FreeKilometerRange do not match")
	assert.Equal(t, expectedResponse.VollkaskoCount, actualResponse.VollkaskoCount, "VollkaskoCount does not match")
}

// All 2 cars in Rotterdam
func TestGetOffersPlattformBaseRegion1(t *testing.T) {
	app := setupApp()
	TestPostOffersPlatform(t)

	req := httptest.NewRequest(http.MethodGet, "/api/offers?regionID=18&timeRangeEnd=1673568000000&timeRangeStart=0&numberDays=3&sortOrder=price-asc&page=0&pageSize=100&priceRangeWidth=20000&minFreeKilometerWidth=1000", nil)
	resp, _ := app.Test(req)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error during GET request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be 200")

	var actualResponse models.OfferQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	expectedResponse := models.OfferQueryResponse{
		Offers: []models.ResponseOffer{
			{
				ID:   "7dd014c0-f473-4210-ad1c-6aa2aca1e6ca",
				Data: "ycmO4hCt6DId5TQoahhsXU220wwGUxzy5qMl4WziRLa1rFgSd2dS5a0zWEgdmYGX6gEAH48AlhCzUS1609wFoPT4sb3pJCEXJm1lQ7k44iI1XUDTKYuLfShxbwyumcfjYgThp3/yQEMdis+PtOr5RYjdPI258cH+VgEuM7+TLaf2mwVROs+c9A/1pFRd8PI2lq8JE7n8Gwr6w6luLseWw/cdYjhO0W/xR0R+KJJwl+n7UN2Y7rwEYNHvHFE1Bj9uPTYr71pH05e4VUmBJamRh3mzNmN+h+wgp9ndG1NdX0o4VvssWekwrurthW6aZi5zIQjlREzYrHX8s9eGJ+gH2w==",
			}, {
				ID:   "87b57605-1ed2-43be-9613-e279d446466c",
				Data: "LeMxLnrv9bMYI0iSDjUn3DCHo1y/SDeAC4ZFHUDO41nQHNxUR5nOgQx3db9Pt8TESK50BEPZzzu4tcVg2qujF3aN0oMr5Bhmc19vgnu533HzIYlmE7454fBL2ercKABOrO3B1ntpSCpAa4zl2H3QEwWvRePE85hzof9HMnngpv1/9abUOzWutvZNZtUae9XFHEPc1sf6+GAESw6HxKbHB2LlG40bM9+jZujlfB535q1UgIVG8S25zG49k6+IB+Lc3enyXuL6F+acT+przcsvcMzgefPXujERGprnqHfCfdnKWg3mRe9bDtgrqT/4Oaw+Cev0+yMgY58WB5yCPCbP6Q==",
			},
		},
		PriceRanges: []models.PriceRange{
			{Start: 0, End: 20000, Count: 2},
		},
		CarTypeCounts: models.CarTypeCounts{
			Small:  1,
			Sports: 0,
			Luxury: 0,
			Family: 1,
		},
		SeatsCount: []models.SeatsCount{
			{NumberSeats: 3, Count: 1},
			{NumberSeats: 2, Count: 1},
		},
		FreeKilometerRange: []models.FreeKilometerRange{
			{Start: 0, End: 1000, Count: 2},
		},
		VollkaskoCount: models.VollkaskoCount{
			TrueCount:  1,
			FalseCount: 1,
		},
	}

	// Assertions for the response structure
	assert.Equal(t, expectedResponse.Offers, actualResponse.Offers, "Offers do not match")
	assert.Equal(t, expectedResponse.PriceRanges, actualResponse.PriceRanges, "PriceRanges do not match")
	assert.Equal(t, expectedResponse.CarTypeCounts, actualResponse.CarTypeCounts, "CarTypeCounts do not match")
	assert.Equal(t, expectedResponse.SeatsCount, actualResponse.SeatsCount, "SeatsCount do not match")
	assert.Equal(t, expectedResponse.FreeKilometerRange, actualResponse.FreeKilometerRange, "FreeKilometerRange do not match")
	assert.Equal(t, expectedResponse.VollkaskoCount, actualResponse.VollkaskoCount, "VollkaskoCount does not match")
}

func TestGetOffersPlattform1(t *testing.T) {
	app := setupApp()
	TestPostOffersPlatform(t)

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=2&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673395200000", nil)
	resp, _ := app.Test(req)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error during GET request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be 200")

	var actualResponse models.OfferQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	expectedResponse := models.OfferQueryResponse{
		Offers: []models.ResponseOffer{
			{ID: "01934a57-7988-7879-bb9b-e03bd4e77b9d", Data: "randomString"},
		},
		PriceRanges: []models.PriceRange{
			{Start: 10000, End: 15000, Count: 4},
		},
		CarTypeCounts: models.CarTypeCounts{
			Small:  1,
			Sports: 2,
			Luxury: 1,
			Family: 0,
		},
		SeatsCount: []models.SeatsCount{
			{NumberSeats: 5, Count: 4},
		},
		FreeKilometerRange: []models.FreeKilometerRange{
			{Start: 100, End: 150, Count: 4},
		},
		VollkaskoCount: models.VollkaskoCount{
			TrueCount:  3,
			FalseCount: 1,
		},
	}

	// Assertions for the response structure
	assert.Equal(t, expectedResponse.Offers, actualResponse.Offers, "Offers do not match")
	assert.Equal(t, expectedResponse.PriceRanges, actualResponse.PriceRanges, "PriceRanges do not match")
	assert.Equal(t, expectedResponse.CarTypeCounts, actualResponse.CarTypeCounts, "CarTypeCounts do not match")
	assert.Equal(t, expectedResponse.SeatsCount, actualResponse.SeatsCount, "SeatsCount do not match")
	assert.Equal(t, expectedResponse.FreeKilometerRange, actualResponse.FreeKilometerRange, "FreeKilometerRange do not match")
	assert.Equal(t, expectedResponse.VollkaskoCount, actualResponse.VollkaskoCount, "VollkaskoCount does not match")
}

/*
func TestGetOffersPlattform2(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=4&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673222400000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOffersPlattform3(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=4&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673222400000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOffersPlattform4(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=2&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673395200000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOffersPlattform5(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=3&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673308800000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOffersPlattform6(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=4&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673222400000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOffersPlattform7(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=2&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673395200000&timeRangeStart=1673222400000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOffersPlattform8(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=4&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673222400000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOffersPlattform9(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=4&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673222400000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOffersPlattform10(t *testing.T) {
	app := setupApp()

	//TODO
	req := httptest.NewRequest(http.MethodGet, "/api/offers?minFreeKilometerWidth=50&numberDays=3&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673308800000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
*/
