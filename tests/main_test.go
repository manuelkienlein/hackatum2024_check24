package tests

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
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
		//log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

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
