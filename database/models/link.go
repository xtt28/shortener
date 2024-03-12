package models

import (
	"crypto/rand"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

// ShortLinkIDLength is the amount of characters that will be contained in a
// short link's ID.
const ShortLinkIDLength = 6

// ShortLink represents a shortened link that will be saved in persistent
// storage.
type ShortLink struct {
	gorm.Model
	// ID is the unique ID of the record in the format of a series of
	// letters and numbers.
	ID string
	// Alias is the optional, user-defined alias that will redirect visitors
	// to the link in combination with its ID.
	Alias sql.NullString `gorm:"uniqueIndex;size:25"`
	// Destination is the full link that the short link will redirect
	// visitors to.
	Destination string
}

// BeforeCreate runs before the creation of a short link and fills the ID field
// with a series of random letters and numbers.
func (l *ShortLink) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := GenerateShortLinkID()
	if err != nil {
		return
	}

	l.ID = id
	return
}

// GenerateShortLinkID creates a random sequence of letters and numbers for use
// as a unique identifier in the ShortLink model.
func GenerateShortLinkID() (string, error) {
	buf := make([]byte, ShortLinkIDLength/2)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}
