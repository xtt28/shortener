package handlers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xtt28/shortener/api/handlers"
	"github.com/xtt28/shortener/database"
	"github.com/xtt28/shortener/database/models"
)

func TestRedirect(t *testing.T) {
	db, err := database.ConnectToTestDatabase()
	assert.NoError(t, err)

	database.MigrateAllModels(db)
	database.DB = db

	router := gin.New()
	router.GET("/:id", handlers.Redirect)

	t.Run("ValidRecordID", func(t *testing.T) {
		record := models.ShortLink{
			Destination: "https://www.google.com",
		}
		dbRes := database.DB.Create(&record)
		assert.NoError(t, dbRes.Error)
		assert.NotEmpty(t, record.ID)
		assert.False(t, record.Alias.Valid)

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/"+record.ID, nil)
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusFound, rec.Code)
			assert.Equal(t, "https://www.google.com", rec.Header().Get("Location"))
		}
	})

	t.Run("ValidRecordAlias", func(t *testing.T) {
		record := models.ShortLink{
			Destination: "https://www.bing.com",
			Alias:       sql.NullString{String: "alias_test", Valid: true},
		}
		dbRes := database.DB.Create(&record)
		assert.NoError(t, dbRes.Error)
		assert.NotEmpty(t, record.ID)
		assert.NotEmpty(t, record.Alias.String)

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/alias_test", nil)
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusFound, rec.Code)
			assert.Equal(t, "https://www.bing.com", rec.Header().Get("Location"))
		}
	})

	t.Run("InvalidRecordID", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/invalid_id", nil)
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Empty(t, rec.Header().Get("Location"))
		}
	})
}
