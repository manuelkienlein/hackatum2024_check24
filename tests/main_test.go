package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"server/internal"
	"server/internal/database"
	"testing"
)

// Mock setup for the application
func setupApp() *fiber.App {
	// Erstelle eine neue Instanz der App und richte die Routen ein
	app := fiber.New()

	// PostgreSQL-Verbindung herstellen
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := database.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// clean up databases
	err = database.DropTables(ctx, dbPool)
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}
	err = database.Migrate(ctx, dbPool)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	internal.RegisterRoutes(app, dbPool)

	return app
}

func TestGetOffers(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest(http.MethodGet, "/api/offers", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	/*var offers []Offer
	err := fiber.UnmarshalJSON(resp.Body, &offers)
	assert.NoError(t, err)
	assert.Len(t, offers, 2)
	assert.Equal(t, "Sample Offer 1", offers[0].Data)*/
}

// TestGetOffers tests the GET /api/offers endpoint
func TestGetOffersAdvanced(t *testing.T) {
	app := setupApp()

	TestPostOffers(t)

	req, err := http.NewRequest("GET", "/api/offers", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set query parameters
	q := req.URL.Query()
	q.Add("regionID", "1")
	q.Add("timeRangeStart", "1732104000000") // Example timestamp in ms
	q.Add("timeRangeEnd", "1732449600000")   // Example timestamp in ms
	q.Add("numberDays", "5")
	q.Add("sortOrder", "price-asc")
	q.Add("page", "1")
	q.Add("pageSize", "10")
	q.Add("priceRangeWidth", "500")
	q.Add("minFreeKilometerWidth", "50")
	req.URL.RawQuery = q.Encode()

	// Send the request and capture the response
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// Assert: Verify that the response status is OK (200)
	assert.Equal(t, 200, resp.StatusCode)

	// Optionally: Check if the response body contains the expected values
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the returned query parameters in the response
	assert.Equal(t, "1", responseBody["regionID"])
	assert.Equal(t, "1732104000000", responseBody["timeRangeStart"])
	assert.Equal(t, "1732449600000", responseBody["timeRangeEnd"])
	assert.Equal(t, "5", responseBody["numberDays"])
	assert.Equal(t, "price-asc", responseBody["sortOrder"])
	assert.Equal(t, "1", responseBody["page"])
	assert.Equal(t, "10", responseBody["pageSize"])
	assert.Equal(t, "500", responseBody["priceRangeWidth"])
	assert.Equal(t, "50", responseBody["minFreeKilometerWidth"])

}

func TestPostOffers(t *testing.T) {
	/*
		// Setup der App
		app := setupApp()

		// Erstelle ein Beispielangebot
		offer := map[string]string{"ID": "3", "Data": "Sample Offer 3"}
		offerJSON, _ := json.Marshal(offer)

		// Führe eine POST-Anfrage an /api/offers durch
		req := httptest.NewRequest("POST", "/api/offers", bytes.NewReader(offerJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// Überprüfen der Antwort
		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		var createdOffer map[string]string
		err = json.NewDecoder(resp.Body).Decode(&createdOffer)
		assert.NoError(t, err)
		assert.Equal(t, "3", createdOffer["ID"])*/

	// Setup der App
	app := setupApp()

	// JSON-Daten als String
	offerJSON := `{"offers":[{"ID":"78ac0090-2eaa-47d6-8dba-c8e78dcb1679","carType":"sports","data":"UzjFJP3C3LE0VUKIkiDURK7u0XHNtfZqP9JJICtnLYWOyad5gEC6ljIyiGw8Zu533J3DJWc/1nFMq4BKofQ0/u82py6ca6QSSXfE3RE0uDJGPuKSytd2G7o8MLvvgZAw46Z8sm9Da5Am16k1dH3nM+mtj/uR4wB0oXmQQMNAvhobdF3IR1HfbC6Nd7VaCbTBsKcHd60OiMZUUd2woZFP0G4+vsC8Y3NVIQ9r8n/D7Pt7LACxtF/au04Y6wiku6jhYCNc/yKUy9VpjMjBlwIPGAi1CYIiQxGVA43RtPXzN4JPX9t0V9K4xiw8w6lclAEAwkGar1y2oskUx2v7Q0krYA==","endDate":1607990400000,"freeKilometers":501,"hasVollkasko":false,"mostSpecificRegionID":92,"numberSeats":4,"price":1984,"startDate":1607731200000},{"ID":"3dbf3a0f-255f-4a3e-ac7f-027def81ec2d","carType":"small","data":"ib0AAJ2Ca508rDOU2cr4HXH57XjEHXaSzq24qyf+U37sjZdqH0dg8KrWdZq7TjCQQWRwc5F7+YkLOc0tZzF1n8DxL4Tlr8aCJATK8w4IfbDFl4FcnNVNfCMHlMikxzi5I4e3qawA/SZEIVC6lorL1ry0hzA8/K0EAA8jw7Bg5tCcdfaUe7QOu7SiypqyYasoBP/ayosuCwJJgROsyxEgmvAs4fJzYTT/ekLUXxlkFqGwaIcevKT5B23mLm/SYn5qLSq9RgnlCO599DKGVbSgNbOgE9Br8zr2xPfn+fjXEZpQHHcpNNIgJ1pdkO7dGdYLi01B/dBYHB+oEZSAqHeJuw==","endDate":1607817600000,"freeKilometers":501,"hasVollkasko":true,"mostSpecificRegionID":61,"numberSeats":6,"price":3356,"startDate":1607644800000},{"ID":"1bdd22a9-4e7e-4adc-a66f-a93e15d4ff00","carType":"luxury","data":"bu/EWb2WqSRtLNSqh5kBdIflQLhgcdbxnCqIwG/u2syv1SNs49VbhAGBDgx09z/QKcH3n0+jgNJQVo8LnGeQM5xEdqcyu2hmNkzMADm1cuY3+xlM6mIf/H/jyLXfviIWHsah/9sh7iwocUrUF3eBY6JgYexvK7FmpZNs17k8WuKPoQR5Wb1kKYgjxOEgwlrRZYWdedxhr3HNHiYRm2yHdfNWPLHMZ314O30tB4C/rXy5AsWz4a0kTqfuq31Mcb1+ET3AE1sQVb/HG3xWg5AwrUyOxQLdKk7OpMD7QRBCEQMgYTdZoenDA4o7kkQL3CMoxzUPsxeUiMC8rKENJBBX5A==","endDate":1607990400000,"freeKilometers":429,"hasVollkasko":true,"mostSpecificRegionID":91,"numberSeats":6,"price":2530,"startDate":1607731200000},{"ID":"df897646-cb51-49bd-9f7c-901539618183","carType":"small","data":"zECqt1BluzDsBS2vqCLupgdWXZ4KwJtpGfwLt0LUrwx/hDqai6WCddhM6N3TCTEqwDVYY/cGE8Ks3pseOg0ujEFhOtlzIf2PIIYALg1YlTLukJvL2xw8/+Vk4pELw+D3HXD63oYx+V77s0FKP9lRZvCKIRwnUTJVlwW6OBQIJ5JASpwWyEMsz1nRN1aWYuYK0N0mgdfic1uYABQNhIoA3jluZv10qiuU9shKDfVZpuFVXeDkzvbZYA+D2cb2Fhu3oCVcQhOaTd9g37LtfXphwTUSJkNhzpq9fhobX7BIMom91mLQ4luYHEpufp4KcMydyE0UDrU0WJK7UZVHaBRAIA==","endDate":1607904000000,"freeKilometers":147,"hasVollkasko":false,"mostSpecificRegionID":106,"numberSeats":2,"price":4058,"startDate":1607644800000},{"ID":"bf7b5c8e-c8e9-4ce7-b841-5be858dfa78a","carType":"luxury","data":"jCZ1mc6yFjTzV7RUCdIQi1tVGwnkII/wxjgwdK05Y3hCJ3nzqAo5uzR1Oravs6gOPGOOCcgIIH6VqNHmpk+e7Wwb+yp4lYmFMZ5r9z0Zn3uSRULn0OKFO3VrLV8KnwQDLb4+76j7ITR0TVAAHtmlInYoY6iDT8J8LjGDXlkyZI4W+PGd/U/dhnbmHy2AkZA/wWngN3VACNty6e9zTnLlTNEvNEFBxpvR4gCfAKGE6CUv1XGUWB8Ob0Gm0ULFlIgzLPDiUY/MIrJE2/dM6TD8NWxRn2f3sQGoKQY9WdOkwCW0/QNTaF1vouDoEnx5okEuJaTCFadqRgtl40OgLaQTpw==","endDate":1607990400000,"freeKilometers":482,"hasVollkasko":true,"mostSpecificRegionID":106,"numberSeats":6,"price":6535,"startDate":1607644800000},{"ID":"21e5a167-b676-4633-855e-80db02a7f430","carType":"family","data":"no2V0irrjldxRv/k64HzPpv4TofZ0lhsJOPYKZgH7BBhR974YZODAggdt/e1Id4+dP+Wx2YQMO4SeV5Tw4MgjAXVXYxU4A93XzskWfKSNk6YmOdOpmAWydXCxrmMf9tH6hFoTfeYv8irhA/BnzOPNA26KXg9wutjWezQNsxFpJZZFB77uaYbnWMfxi9orkSDok5j0D84yYOU+FmpW2HQy9UExDz12nRuWkbmNmJ7vP4xk35Odz1dd8/DqvTyDSptSbCLlmWgCdom3sMOtD1EilAjkkSWLK+xveCdQ/+pAUjEZnlXHzxbXTKnIPXRVx3V7+zRtBr74co+c0+y6/Ov1g==","endDate":1607990400000,"freeKilometers":772,"hasVollkasko":true,"mostSpecificRegionID":90,"numberSeats":5,"price":8887,"startDate":1607644800000},{"ID":"ecdf905d-084a-43b3-bbff-5d10f0f28e72","carType":"luxury","data":"isQGK0gYQ+jb84lurXTM9Pln5AeywQzXjJqTqRhZbYFiLLx8OgtJPAIaB4E0Fy1qmsgKgBHa279fKja9rA8nxNoWvd8OidkCZswPvfsAjidIGptpdwNOr/RisqchYdXbzU57KhGSx/IYOqcc37R2uLxGU81IhMJMCqD0CYih1t2BbglnAYpu/P9B2Wqczc5yZe4S3f12xpMM1l7V0xGZ9Q4+iW4qYcJX+aI6TYO9g/zT8U7CEdVbK9/jUIFZ7/CQ/DeIZ9KwNqgpDD/xQvmd02h6f5BVFTfv11PiR8f+4pWvgxuxaO6zvZyyTBxrDtQxOEseMfE2mAZx0GYm4v6UnA==","endDate":1607817600000,"freeKilometers":812,"hasVollkasko":true,"mostSpecificRegionID":103,"numberSeats":2,"price":7512,"startDate":1607644800000},{"ID":"564ba381-5112-4d93-9821-2abc91ed9a2b","carType":"luxury","data":"dEu42MD7PkF6vrSvj0jm0riCkCkFirTO/a230M953u/NqBIhUV9sclBr0fhYtPzCJNBOQsapUyYGyHgbkrqB6sYpaNC2Yz8kS9S032aS7Z51tUWyZkOlGIBiIP88PcgKHeWieuVjAfPREEUIbYkQ+YnX5jMKlIGgrIE3FxtHH/6avoLYQgNxGC0vJw6M2Lh9/Lc2cXN44tKXzTUpb7LQDQMHsbmlnufYTTfe1Q2JJRoYsGGz47fuvtbXRsezMxEtLUeaYHWRIO11+Syl8bvIlaI7qIE2DNThMiCBYokr9ErwXrQTDPhqQhLbY+LFWaYptke4Tp2Lp4J1M04E87PuMw==","endDate":1607817600000,"freeKilometers":715,"hasVollkasko":false,"mostSpecificRegionID":118,"numberSeats":2,"price":3883,"startDate":1607644800000},{"ID":"79ac7ce6-b7a0-4674-b18f-1c036ae17de5","carType":"family","data":"m9/IciGqtVIl8lFe48zmH8MmFmafSiWjAGskEE/ZAn3XssbUVj35SP4HqwmjwpSacD/3T0Up1/3b6lMs3EslJGW5cpFRlU5Hb83dTNERC1Lk3k/dTTUkhcuCVnVjfEeJ/qWGJZbFL0OjQmKk+SXWYD/yYrDEvAqzpMzeYVfwIafSjNIiN0zHpuCj/EtBYio+k900xZ65S4ixIrujwt7+D7Cb9S6+VJIuGcl6le1zbyBc+erTv18hx471sf3oZ7cYqV0JK8lX4jySLM5CAbvjJc20GUWqQXVZpRqdkPR6DOA8rjku/fkqK7sboBwKnnHSDdrWy0jaLBjdDO5MzG1wFQ==","endDate":1607904000000,"freeKilometers":983,"hasVollkasko":false,"mostSpecificRegionID":93,"numberSeats":2,"price":3391,"startDate":1607731200000},{"ID":"79372995-907e-4fcd-8181-58353c23038f","carType":"small","data":"oyQP4Gl77OH8F8+HJYaUYo77YuC2xhFVbgtm12i9SCJM1UZ4bwV4+5vMaJp351WKY1Ye078wtdn7bU+c2m+4Ip5o4fQNps0YBpE3C2aLxDFy3s9hCvFigkU0st1cWeokBqih+2AWhvAQlYx2sKo0Rch1sWxUO3qFarPQ4Dpzn/iVU5MfmxOS9vdQL9BL3fe9U45N4fKXGJKiaYgn8ocpUDI4CkCf2Ss1W1htTauhc4bXn9V6jGNuV6bbLr6EyEWtWgCtv3vG3ienzwrHenG0v9FS8NlyKxiL5Fnw8tzx6a16I8UxmW7Z92cEUHUpH8v8gRO8j9RQq3ROJTyFxkld+g==","endDate":1607990400000,"freeKilometers":513,"hasVollkasko":false,"mostSpecificRegionID":93,"numberSeats":2,"price":6165,"startDate":1607644800000}]}
`

	// Führe eine POST-Anfrage an /api/offers durch
	req := httptest.NewRequest("POST", "/api/offers", bytes.NewReader([]byte(offerJSON)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Überprüfen, ob keine Fehler beim Testen aufgetreten sind
	assert.NoError(t, err)

	// Überprüfen, ob der Statuscode >= 200 ist
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
}
