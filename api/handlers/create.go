package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xtt28/shortener/database"
	"github.com/xtt28/shortener/database/models"
	"gorm.io/gorm"
)

// validAliasPattern is a pattern representing valid short URL aliases; only
// alphanumeric characters and underscores are permitted. The database has a
// limit of 20 characters per alias, but this is not handled by this pattern
// and must be checked in combination with this.
var validAliasPattern = regexp.MustCompile(`\w+$`)

// createErrorTemplate is the HTML template to use when sending a response to
// indicate an error.
const createErrorTemplate = "create_error.go.html"

// IsURLValid returns whether the given string is a valid HTTP URL that will be
// accepted by the link shortener.
func IsURLValid(urlText string) bool {
	uri, err := url.ParseRequestURI(urlText)
	if err != nil {
		return false
	}

	return (uri.Scheme == "http" || uri.Scheme == "https") && uri.Host != ""
}

// CreateView serves the "create" template, which provides a form for creating
// short links.
func CreateView(c *gin.Context) {
	c.HTML(http.StatusOK, "create.go.html", gin.H{})
}

// Create is a POST API endpoint that creates a link with the given form data
// and returns an HTML value with the result or an error for use with htmx.
func Create(c *gin.Context) {
	dest := c.PostForm("destination")
	alias := strings.ToLower(c.PostForm("alias"))

	if dest == "" {
		c.HTML(http.StatusBadRequest, createErrorTemplate, gin.H{
			"Description": "Please provide a destination.",
		})
		return
	}

	if !IsURLValid(dest) {
		c.HTML(http.StatusBadRequest, createErrorTemplate, gin.H{
			"Description": "Please provide a valid URL.",
		})
		return
	}

	if len(alias) > 20 {
		c.HTML(http.StatusBadRequest, createErrorTemplate, gin.H{
			"Description": "The alias must have a length of 20 characters or less.",
		})
		return
	}

	if alias != "" && !validAliasPattern.MatchString(alias) {
		c.HTML(http.StatusBadRequest, createErrorTemplate, gin.H{
			"Description": "Only letters, numbers and underscores are permitted in the alias.",
		})
		return
	}

	aliasNullable := sql.NullString{
		String: alias,
		Valid:  alias != "",
	}

	target := &models.ShortLink{
		Destination: dest,
		Alias:       aliasNullable,
	}

	res := database.DB.Create(target)
	err := res.Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.HTML(http.StatusConflict, createErrorTemplate, gin.H{
				"Description": "A link with this alias already exists.",
			})
			return
		} else {
			log.Print(err)
			c.HTML(http.StatusInternalServerError, createErrorTemplate, gin.H{
				"Description": "A system error occurred. Please try again later.",
			})
			return
		}
	}

	id := target.ID
	root := os.Getenv("ROOT") + "/v/"
	shortLinkText := root + id
	var aliasText string
	if alias != "" {
		aliasText = root + alias
	}
	c.HTML(http.StatusOK, "create_success.go.html", gin.H{
		"Link":          shortLinkText,
		"LinkWithAlias": aliasText,
	})
}
