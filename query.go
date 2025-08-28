package bnm

import (
	"fmt"
	"time"
)

// Supported API languages.
const (
	LANG_EN = "en"
	LANG_RO = "ro"
	LANG_RU = "ru"
)

// Query represents a request for exchange rates on a specific date and in a specific language.
type Query struct {
	Date time.Time
	Lang string
}

// NewQuery creates a new Query for the given date and language.
func NewQuery(date time.Time, lang string) Query {
	return Query{
		Date: date,
		Lang: lang,
	}
}

// RequestURL returns the URL used to request exchange rates from the BNM API.
func (q Query) RequestURL() string {
	return fmt.Sprintf("http://www.bnm.md/%s/official_exchange_rates?get_xml=1&date=%s", q.Lang, q.dateToStr())
}

// ID returns a unique identifier for the query.
func (q Query) ID() string {
	return q.Lang + "_" + q.dateToStr()
}

func (q Query) dateToStr() string {
	const dateFormat = "02.01.2006"
	return q.Date.Format(dateFormat)
}
