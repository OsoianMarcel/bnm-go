package bnm

import (
	"encoding/xml"
	"fmt"
)

// Currency represents a currency as returned by the official API.
type Currency struct {
	ID      string  `xml:"ID,attr" json:"id"`
	Code    string  `xml:"CharCode" json:"code"`
	NumCode int     `json:"num_code"`
	Nominal int     `json:"nominal"`
	Name    string  `json:"name"`
	Value   float32 `json:"value"`
}

// Response represents the API response containing exchange rates for multiple currencies.
type Response struct {
	Date       string     `xml:"Date,attr" json:"date"`
	Name       string     `xml:"name,attr" json:"name"`
	Currencies []Currency `xml:"Valute" json:"currencies"`
}

// FindByCode searches for a currency by its three-letter code.
// It returns the currency and true if found, or an empty Currency and false otherwise.
func (r Response) FindByCode(code string) (Currency, bool) {
	for _, val := range r.Currencies {
		if val.Code == code {
			return val, true
		}
	}

	return Currency{}, false
}

// unmarshalResponse parses XML data into a Response struct.
// Returns an error if the XML cannot be decoded.
func unmarshalResponse(data []byte) (Response, error) {
	var res Response
	err := xml.Unmarshal(data, &res)
	if err != nil {
		return Response{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return res, nil
}
