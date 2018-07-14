package bnm

import (
	"fmt"
	"time"
)

// Query struct used to describe request query
type Query struct {
	Locale string
	Date   time.Time
}

// NewQuery is used to create new Query instance
func NewQuery(locale string, date time.Time) Query {
	return Query{Locale: locale, Date: date}
}

// GenerateURI is used to generate URI used for GET request
func (q Query) GenerateURI() string {
	return fmt.Sprintf("http://www.bnm.md/%s/official_exchange_rates?get_xml=1&date=%s", q.Locale, q.DateToString())
}

// DateToString converts date to string
func (q Query) DateToString() string {
	const dateFormat = "02.01.2006"
	return q.Date.Format(dateFormat)
}

// GetID returns query id
func (q Query) GetID() string {
	return q.Locale + "_" + q.DateToString()
}
