package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xtt28/shortener/database"
	"github.com/xtt28/shortener/database/models"
	"gorm.io/gorm"
)

// Redirect will redirect the user to the URL that has the given ID or alias. If
// no such URL has been identified in the database, the handler will return a
// 404 Not Found error.
func Redirect(c *gin.Context) {
	id := c.Param("id")

	target := models.ShortLink{}
	res := database.DB.
		Where(&models.ShortLink{ID: id}).
		Or(&models.ShortLink{Alias: sql.NullString{String: id, Valid: true}}).
		First(&target)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, "could not find link with given ID")
		} else {
			c.String(http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.Header("Location", target.Destination)
	c.String(http.StatusFound, "redirecting")
}
