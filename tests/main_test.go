package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"server/internal/controller"
	"server/internal/database"
	"server/internal/framework"
	"server/internal/repository"
	"server/internal/service"
	"sync"
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

	// Init components
	offerRepo := repository.NewOfferRepository(dbPool)
	offerService := service.NewOfferService(offerRepo)
	offerController := controller.NewOfferController(offerService)

	framework.RegisterRoutes(app, offerController)

	return app
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
	q.Add("regionID", "0")
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

func TestGetOffers2(t *testing.T) {
	// Setup der App
	app := setupApp()

	TestPostOffers(t)

	// Führe eine GET-Anfrage an /api/offers durch
	// /api/offers?minFreeKilometerWidth=50&numberDays=4&page=0&pageSize=100&priceRangeWidth=10&regionID=0&sortOrder=price-asc&timeRangeEnd=1673568000000&timeRangeStart=1673222400000
	req := httptest.NewRequest("GET", "/api/offers?minFreeKilometerWidth=50&numberDays=1&page=0&pageSize=100&priceRangeWidth=10&regionID=37&sortOrder=price-asc&timeRangeEnd=1692179200000&timeRangeStart=1491920000000&onlyVollkasko=true", nil)
	resp, err := app.Test(req)

	// Überprüfen, ob keine Fehler beim Testen aufgetreten sind
	assert.NoError(t, err)

	// Überprüfen, ob der Statuscode = 200 ist
	assert.Equal(t, 200, resp.StatusCode)
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
	offerJSON := `{"offers":[{"ID":"68ed0a29-a7ae-42b4-bdfd-2f35462828ca","carType":"sports","data":"9Wm0fKJQxnqV+1R8Eie8g2D2hXjsX3H1Yri3S46Ev+Sil/L1sbn6UeaKQLxif+IlR1Srgx12Xpyaqx+gS+T4UnnSxyNKnpK49Jdag6i2BqmdIt7c3LZz891HVCMtrpTqevN5HaiRXy5yLGdix8PMAdxfj7qY5wgR1EXcn5oCypNZoMrbmV378Hd5vPWwsHR1p9ZfbMRAFZMbcfZxomM4kwacugBUy0rsw2CzhXj7B4z88jwZJVSWTzKiBbVdisdQpCRgEgM+I6nvDk8WEms32VIWUGaE8apqmMMH+DxpJVVnudGlb2nZvKzTPoUmRHv58O1l2BhwhA393t7WstCMGA==","endDate":1592179200000,"freeKilometers":508,"hasVollkasko":true,"mostSpecificRegionID":81,"numberSeats":2,"price":4481,"startDate":1591920000000},{"ID":"1b993f42-6450-4ed1-a834-7c4a024f10f6","carType":"sports","data":"/QNVpWu5uBIREBD1noR89RSYgnSQX5af7fyrp4u1GCwdKPd+5SMZk3VMg+WQ5rF4scRKfHwGTWYm05ayO+iNHwL5iQD6rIMg5R6q4JsGQpE88ppPnAYHJZl9TvUsZt2gChMxgGbnn+3UlYg75I3GlqdjCODs4sXuGgngV59M2TLMRQF6fKcwLaHaYm27vi4dsvaw1Jvh4T8za9zSHbSq/A4Q77qjZGR+26FtG8sDrkKTNS1y+Z8k20qlH+8/Ud39k3++Q+K1rD30zD1sLCaGWtIF8E+QxniBQynIjKaccVXsvmJUUH9Vh6+YJMqmCmxsHI8WsRKyD4LCZfgW9UcYAQ==","endDate":1592092800000,"freeKilometers":695,"hasVollkasko":true,"mostSpecificRegionID":120,"numberSeats":2,"price":2362,"startDate":1591920000000},{"ID":"f72473c8-b00e-46f2-bd23-6b2111b15a2d","carType":"sports","data":"ue+0jMLw6kVmt4wsyJDVZgWPyu9j+RxNmVC1V7ogi1gaxEas7cLpzmzbT+9dZOmp95+fhVKKkkGNkKwX1nSTJzs6myg5in1FyhnK0WSfZX+Z0hpb3y0UDQEdDmKMOln109nbtb9y+u8Jb2SBi4UVPJ/4I7PwssS+edsvfkdD6aeLVYEKZiglHeFDLaM/bbWmV8Z5JiLHUzvPsj1Q2alKBmmcS7uaD9XozSsY4Y61Hs+OryrM06Apt6eevof4+06sWKwfRMGGt3wylzXJoTANbIESC7WQwi4TLQMgiyCuplxgdQPWohgqswo+YvX+1V0QJoZG7VLERwfi+Rf67d/7rQ==","endDate":1592006400000,"freeKilometers":878,"hasVollkasko":false,"mostSpecificRegionID":60,"numberSeats":5,"price":5154,"startDate":1591833600000},{"ID":"607031c8-7507-40c7-9a24-582d7519e37c","carType":"family","data":"WRM24YCYo+bo/WFsYkog120uRXx61KqvdhxbvggdhTkVMd6awSO6W8lyS3SYh1K1hOyZ532w8aQOpC3lQ62QnJE07AuAovRghmSu7z2cqC9tr6VlCbCIIwQurSn5rfdIYWt/77A+nBhKoN0nOXboZplTK4D44TuDo9XSgZnYzlI3ZPbzWwZakHfWLaqNUYIVvkJMNoWIoj7b11lfY76VMdK78JuXYUci79BrVL/tQV8a4m/g7NGZzhJot4GA1K7TAtk/kynvrpfQ/B+FbBg4AN17fIF0F8mZdSS9JKEW3e+/xec8Cm8widDITRTB50/SJzsgBEQwCowmXS6Oi5YV2A==","endDate":1592092800000,"freeKilometers":236,"hasVollkasko":false,"mostSpecificRegionID":105,"numberSeats":4,"price":4438,"startDate":1591920000000},{"ID":"d77bd5e1-1617-495a-b2c9-d325f7d8e1de","carType":"family","data":"M1lLibnJiAhLpnbD5gMcJGl2ygpPzw6VcnImPxcrWl1rSDwwmImPWeC/ZeSVBuSMHdOzSq4BzWPQdvyUQNjNcIpZMx9IUj2Z8damdU2GA4K5OBqmu8sk3W5BMy2cKf8VNOqXdcZ7k/zP3hgqKX6uXNG/qquIuKkb1MEbcXP1C6Kks2+dnlQjFJbJVL9OBkk/ihMwYYszD+TdNurf4mWIGKM12lOZk4uktzVR2QqmTTkdyBQas7NVMxvgjaEcaq0Ox3+KJkIlhg01tYTehbSOrNXLBHsHv6SeM34mbQgWp5F+XmWQgbLVaDedDvlkOCOevFh3nAe19K+CmslY/4E52g==","endDate":1592179200000,"freeKilometers":102,"hasVollkasko":false,"mostSpecificRegionID":88,"numberSeats":5,"price":2189,"startDate":1591833600000},{"ID":"6562eef4-cd57-468b-affe-07b3a834af62","carType":"sports","data":"3c742E0pUpd6DVoDaEQSGBf0MvIfKPcbBkhi02Oj0hajU1E3ISOKDi8yRs0zoYpTiVFghVIKqRZBstGXD6KFEsedY7jLre28MlEY/RT7GU8vO2GUY03xM9maHeW9xhB3E8dh/n9NxksyHmq/pMzlG9XmHk8pfB5j+1c5Qx69hK8tu9M9Zz9eUTkHgEe4+XL0tlvl1qdZHQAgAntxjTVbzBPDbXhExwzv3xZqPOZumjZSNpY/mm8iS95gxqavEbP67bBLWk/SlAG4QEiQztIkj+m1mTgQ5ez5ifpFnzQAM3JuAinvOyvmdOixaBjxvpbyvWTltwD0gZGVZv8Yr7ipWA==","endDate":1592179200000,"freeKilometers":314,"hasVollkasko":true,"mostSpecificRegionID":108,"numberSeats":4,"price":5408,"startDate":1592006400000},{"ID":"8d499624-78ca-4a4f-aad4-4cee50fb732a","carType":"family","data":"14D1rIwmBYViwyWDdbHCPY62/WatAXdJ4ILExRQD7kTczaHOK5y5s7SgLPCQPYkbhJLUcqpUPYFO5Q+vMMHISBMejNMIwajyySC/wZwYKbk2pvWrE8GHwvzrSGr7ZGtA5l9Za1KaLQvaokjk3kDZRcO/r8Z+75Uqbp6iHOSPO86RozwiAIkcagOIIqm9a0B2y2Kak+zq4tob4bj71yeGaXSLY0B2Q3ro9qPQNpUX9doV4eo5ns7qQWVm3tyZr+u9n5IGF/fNwajTys04fMQKnSufFKCrxEL44irjW8Ra9ImD5zVAB9GKX8fRapY+0lVsTU+dPC8gOyQ4tInv8K04zQ==","endDate":1592092800000,"freeKilometers":530,"hasVollkasko":false,"mostSpecificRegionID":104,"numberSeats":2,"price":2019,"startDate":1591833600000},{"ID":"b7d5e28e-aa11-44d1-a13e-2ea82ad8b144","carType":"small","data":"mbo55qiV8s3wh8Kjkj6rOtG1P2mzT3J3VNqXqBYUb3kMWhf38PcDKzGt389aHmhKzTfGFPX/3hhHKcE6VOeONoEn/plECqNvxZc+/+F4mp3HNqcEZw4CfdSBkrfIVxP9dxiqQrskZ8YDR7VFvwU9EAGR+pE3a5Pb3BZVbmrti4u/uocAznHoDbW1niiEeecDnOGDwg0FhwaAOW2dVUPrj5N/VwD5fEOghigXNcMRDPofWgey0arMl7v6jYk7+OiTvvw+TWtgLIMDuIl+dfuZuGlNRpdtbngUqtkZVSZIl2tIUczRNMtaxhgf903Yzr8DLsHnBaZa/ft7fYUt0uMTXg==","endDate":1592092800000,"freeKilometers":290,"hasVollkasko":false,"mostSpecificRegionID":100,"numberSeats":4,"price":10007,"startDate":1591920000000},{"ID":"8db1cafc-d328-4f0b-8d8f-edf340ed2043","carType":"small","data":"FNhoPrOZ85qvmmOfy8gVU5mOBY4xRtZruEJUoYS30PiDvyQCmzrbqBOVDVi0DZGtaxx638k9Bh/ZdrVwIg2mzoa/u1GRYD5a8WOpO5ImrURw0XEqzpcd0iI3axmcEr24AiDtlutFWnkb4D3rOwJIKg+8ec4pXDstiAVmYVK6sWa6jtjUYDaSlhKMvu163aNS9M50PBbsngHzRdq7rske6XMh7zxFcmsPngOMxmyErCd2jxGU0Wxr+RBSTACsjrsjqbuc02NphuhMzNj+cgsIWHInE9NZlX0t/11OQiETaVgjzZCLETi6G6YibaQOYmCkGdWqVjLTb5yEFAi3WlE2Nw==","endDate":1592179200000,"freeKilometers":745,"hasVollkasko":false,"mostSpecificRegionID":93,"numberSeats":3,"price":1811,"startDate":1591833600000},{"ID":"f63c607a-9140-42c8-bbe7-6e078cc54e75","carType":"family","data":"QdZoDKwHUOM7sIOLTrwCtwkDvVZroHxT35aGFgs8hsMfv4tcJZUxK4Pqs7+YYFzNWMVISTzIEhxAG4Vra6caGfQK620rDLp6hidPmpBjmRidBLbEbkUe9iQahw5zBGXCh3okVyA7sfe47HJXMMMV5FFrB8gZLXc3TT3hXtzBP+qUt8xdOV65oBpU/Yl8jVHM2/dmxUev5udzfJ4iqgmNHSKkTXgDuQxkrWmgepMhM5ENuce9uLjpzALJSVeJkPmZMdyKo353Moiafbea03+bt7Holgnh94q5A2O+/7ftlpqgLSMExYX5qwWn90R2f9SywBAi9K+yDRxm9+8m+01BPg==","endDate":1592179200000,"freeKilometers":133,"hasVollkasko":false,"mostSpecificRegionID":120,"numberSeats":5,"price":7335,"startDate":1591833600000}]}`
	// Führe eine POST-Anfrage an /api/offers durch
	req := httptest.NewRequest("POST", "/api/offers", bytes.NewReader([]byte(offerJSON)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Überprüfen, ob keine Fehler beim Testen aufgetreten sind
	assert.NoError(t, err)

	// Überprüfen, ob der Statuscode = 200 ist
	assert.Equal(t, resp.StatusCode, 200)
}

func TestPostPerf(t *testing.T) {
	// Setup der App
	app := setupApp()

	numOffers := 1000

	// Erzeuge die JSON-Daten für die POST-Anfrage
	offers := []byte(`{"offers": [`)

	for i := 0; i < numOffers; i++ {
		// Erzeuge eine zufällige UUID für jedes Angebot
		id := uuid.New().String()

		// Baue das Angebot mit der zufälligen UUID
		offer := fmt.Sprintf(`
		{
			"ID": "%s",
			"data": "string",
			"mostSpecificRegionID": 5,
			"startDate": 1732104000000,
			"endDate": 1732449600000,
			"numberSeats": 5,
			"price": 10000,
			"carType": "luxury",
			"hasVollkasko": true,
			"freeKilometers": 120
		}`, id)

		// Füge das Angebot zur JSON-Datenstruktur hinzu
		if i > 0 {
			offers = append(offers, []byte(",")...)
		}
		offers = append(offers, []byte(offer)...)
	}

	// Beende das Array und das JSON-Objekt
	offers = append(offers, []byte(`]}`)...)

	// Führe eine POST-Anfrage an /api/offers durch
	req := httptest.NewRequest("POST", "/api/offers", bytes.NewReader(offers))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Überprüfen, ob keine Fehler beim Testen aufgetreten sind
	assert.NoError(t, err)

	// Überprüfen, ob der Statuscode = 200 ist
	assert.Equal(t, resp.StatusCode, 200)
}

func TestPostPerfConcurrency(t *testing.T) {
	// Setup der App
	app := setupApp()

	// Anzahl der zu generierenden Angebote
	numOffersPerBatch := 1000

	// Anzahl der parallelen Anfragen
	concurrency := 32

	// Anzahl der Batches
	numBatches := 10 //100

	// Warteschleife für parallele Tests
	var wg sync.WaitGroup

	// Führe parallele Tests aus
	for n := 0; n < numBatches; n++ {

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				// Erzeuge die JSON-Daten für die POST-Anfrage
				offers := []byte(`{"offers": [`)

				for j := 0; j < numOffersPerBatch/concurrency; j++ {
					// Erzeuge eine zufällige UUID für jedes Angebot
					id := uuid.New().String()

					// Baue das Angebot mit der zufälligen UUID
					offer := fmt.Sprintf(`
				{
					"ID": "%s",
					"data": "string",
					"mostSpecificRegionID": 5,
					"startDate": 1732104000000,
					"endDate": 1732449600000,
					"numberSeats": 5,
					"price": 10000,
					"carType": "luxury",
					"hasVollkasko": true,
					"freeKilometers": 120
				}`, id)

					// Füge das Angebot zur JSON-Datenstruktur hinzu
					if j > 0 {
						offers = append(offers, []byte(",")...)
					}
					offers = append(offers, []byte(offer)...)
				}

				// Beende das Array und das JSON-Objekt
				offers = append(offers, []byte(`]}`)...)

				// Führe eine POST-Anfrage an /api/offers durch
				req := httptest.NewRequest("POST", "/api/offers", bytes.NewReader(offers))
				req.Header.Set("Content-Type", "application/json")
				resp, err := app.Test(req)

				// Überprüfen, ob keine Fehler beim Testen aufgetreten sind
				assert.NoError(t, err)

				// Überprüfen, ob der Statuscode = 200 ist
				assert.Equal(t, resp.StatusCode, 200)
			}(i)
		}

		// Warten, bis alle Goroutinen abgeschlossen sind
		wg.Wait()
	}
}
