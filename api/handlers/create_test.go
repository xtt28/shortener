package handlers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xtt28/shortener/api/handlers"
	_ "github.com/xtt28/shortener/api/test/init"
	"github.com/xtt28/shortener/database"
	"github.com/xtt28/shortener/database/models"
)

// initTestDB connects to an in-memory database and sets the database.DB
// variable to the newly created *gorm.DB.
func initTestDB(t *testing.T) {
	db, err := database.ConnectToTestDatabase()
	assert.NoError(t, err)

	database.MigrateAllModels(db)
	database.DB = db
}

func TestIsURLValid(t *testing.T) {
	t.Run("InvalidURLs", func(t *testing.T) {
		assert.False(t, handlers.IsURLValid("google"))
		assert.False(t, handlers.IsURLValid("youtube.com"))
		assert.False(t, handlers.IsURLValid("ftp://"))
		assert.False(t, handlers.IsURLValid("javascript:alert(1)"))
		assert.False(t, handlers.IsURLValid("file://C:/Users/Brad/Pictures/cat.png"))
		assert.False(t, handlers.IsURLValid("jdbc:mysql://127.0.0.1:3306/fake_db"))
		assert.False(t, handlers.IsURLValid("//"))
	})

	t.Run("ValidURLs", func(t *testing.T) {
		assert.True(t, handlers.IsURLValid("https://google.com"))
		assert.True(t, handlers.IsURLValid("https://www.google.com"))
		assert.True(t, handlers.IsURLValid("http://example.com"))
		assert.True(t, handlers.IsURLValid("http://example.net"))
		assert.True(t, handlers.IsURLValid("http://github.com/foo"))
		assert.True(t, handlers.IsURLValid("http://example.net/foo?bar=baz"))
	})
}

func TestCreateView(t *testing.T) {
	router := gin.New()
	router.LoadHTMLGlob("web/templates/*")
	router.GET("/", handlers.CreateView)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)

	if assert.NoError(t, err) {
		assert.NotPanics(t, func() { router.ServeHTTP(rec, req) })
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreate(t *testing.T) {
	router := gin.New()
	router.LoadHTMLGlob("web/templates/*")
	router.POST("/api/create", handlers.Create)

	t.Run("ValidDataNoAlias", func(t *testing.T) {
		initTestDB(t)
		rec := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("destination", "https://github.com")

		req, err := http.NewRequest(http.MethodPost, "/api/create", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			destRecord := &models.ShortLink{
				Destination: "https://github.com",
			}
			res := database.DB.First(&destRecord)
			if assert.NoError(t, res.Error) {
				id := destRecord.ID
				assert.Contains(t, rec.Body.String(), id)
			}
		}
	})

	t.Run("ValidDataWithAlias", func(t *testing.T) {
		initTestDB(t)
		rec := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("destination", "https://go.dev")
		formData.Set("alias", "go")

		req, err := http.NewRequest(http.MethodPost, "/api/create", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "go")

			destRecord := &models.ShortLink{
				Destination: "https://go.dev",
				Alias:       sql.NullString{Valid: true, String: "go"},
			}
			res := database.DB.First(&destRecord)
			if assert.NoError(t, res.Error) {
				id := destRecord.ID
				assert.Contains(t, rec.Body.String(), id)
			}
		}
	})

	t.Run("InvalidNoDestination", func(t *testing.T) {
		initTestDB(t)
		rec := httptest.NewRecorder()

		formData := url.Values{}

		req, err := http.NewRequest(http.MethodPost, "/api/create", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			res := database.DB.First(&models.ShortLink{})
			assert.Error(t, res.Error)
		}
	})

	t.Run("InvalidAliasTooLong", func(t *testing.T) {
		rec := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("destination", "https://example.com")
		formData.Set("alias", "ThisAliasIsVeryVeryLong")

		req, err := http.NewRequest(http.MethodPost, "/api/create", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			res := database.DB.First(&models.ShortLink{})
			assert.Error(t, res.Error)
		}
	})

	t.Run("InvalidAliasForbiddenChars", func(t *testing.T) {
		rec := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("destination", "https://example.com")
		formData.Set("alias", "! @#$%^&*()")

		req, err := http.NewRequest(http.MethodPost, "/api/create", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			res := database.DB.First(&models.ShortLink{})
			assert.Error(t, res.Error)
		}
	})

	t.Run("InvalidURL", func(t *testing.T) {
		rec := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("destination", "invalid_url")

		req, err := http.NewRequest(http.MethodPost, "/api/create", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			res := database.DB.First(&models.ShortLink{})
			assert.Error(t, res.Error)
		}
	})

	t.Run("DuplicateAliases", func(t *testing.T) {
		rec := httptest.NewRecorder()
		rec2 := httptest.NewRecorder()

		formData := url.Values{}
		formData.Set("destination", "https://apple.com")
		formData.Set("alias", "apple")

		req, err := http.NewRequest(http.MethodPost, "/api/create", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if assert.NoError(t, err) {
			router.ServeHTTP(rec, req)
			router.ServeHTTP(rec2, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, http.StatusConflict, rec2.Code)

			res := database.DB.Find(&models.ShortLink{})
			assert.NoError(t, res.Error)
			assert.Equal(t, int64(1), res.RowsAffected)
		}
	})
}
