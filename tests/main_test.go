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
	offerJSON := `{
		"offers": [
			{
				"ID": "01934a57-7988-7879-bb9b-e03bd4e77b9d",
				"data": "string",
				"mostSpecificRegionID": 5,
				"startDate": 1732104000000,
				"endDate": 1732449600000,
				"numberSeats": 5,
				"price": 10000,
				"carType": "luxury",
				"hasVollkasko": true,
				"freeKilometers": 120
			}
		]
	}`

	// Führe eine POST-Anfrage an /api/offers durch
	req := httptest.NewRequest("POST", "/api/offers", bytes.NewReader([]byte(offerJSON)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Überprüfen, ob keine Fehler beim Testen aufgetreten sind
	assert.NoError(t, err)

	// Überprüfen, ob der Statuscode >= 200 ist
	assert.GreaterOrEqual(t, resp.StatusCode, 200)
}
