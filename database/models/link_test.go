package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xtt28/shortener/database/models"
)

func TestGenerateShortLinkID(t *testing.T) {
	id, err := models.GenerateShortLinkID()
	if assert.NoError(t, err) {
		assert.Len(t, id, models.ShortLinkIDLength)
	}
}
