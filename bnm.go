package bnm

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

// Inst it's main package structure
type Inst struct{}

// NewBnm function is used to
func NewBnm() Inst {
	return Inst{}
}

// Request BNM rates by using Query
func (bnm Inst) Request(q Query) (Result, error) {
	res, err := getRequest(q.GenerateURI())

	if err != nil {
		return Result{}, err
	}

	return parseResponse(res)
}

// Do a get request by URI and return response
func getRequest(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return []byte{}, err
	}

	// It is important to defer resp.Body.Close(), else resource leaks will occur
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.New("bad HTTP status code")
	}

	// Will print site contents (HTML) to output
	return io.ReadAll(resp.Body)
}

// Parse XML response
func parseResponse(resp []byte) (Result, error) {
	var c Result
	err := xml.Unmarshal(resp, &c)
	if err != nil {
		return Result{}, err
	}

	return c, nil
}
