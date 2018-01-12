package bnm

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"bytes"
)

// Test parseResponse() when error occured
func Test_parseResponse_Error(t *testing.T) {
	_, err := parseResponse([]byte{})
	if err == nil {
		t.Error("the function shoud return error")
	}
}

// Test parseResponse() by parsing XML
func Test_parseResponse_Result(t *testing.T) {
	xml := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ValCurs Date="05.08.2017" name="Cursul oficial de schimb">
		  <Valute ID="47">
			<NumCode>978</NumCode>
			<CharCode>EUR</CharCode>
			<Nominal>1</Nominal>
			<Name>Euro</Name>
			<Value>21.2997</Value>
		  </Valute>
		  <Valute ID="44">
			<NumCode>840</NumCode>
			<CharCode>USD</CharCode>
			<Nominal>1</Nominal>
			<Name>Dolar S.U.A.</Name>
			<Value>17.9948</Value>
		  </Valute>
		  <Valute ID="35">
			<NumCode>946</NumCode>
			<CharCode>RON</CharCode>
			<Nominal>1</Nominal>
			<Name>Leu romanesc</Name>
			<Value>4.6680</Value>
		  </Valute>
		</ValCurs>
	`

	res, err := parseResponse([]byte(xml))
	if err != nil {
		t.Error(err)
	}

	expectedLen := 3
	resultLen := len(res.Rates)
	if resultLen != expectedLen {
		t.Errorf("expected length: %d, result: %d", expectedLen, resultLen)
	}
}

// Test getRequest() method when error
func Test_getRequest_Error(t *testing.T) {
	inexistentSite := "http://localhost:-41"
	_, err := getRequest(inexistentSite)
	if err == nil {
		t.Error("the function shoud return an error")
	}
}

// Test getRequest() method with test server
func Test_getRequest_Success(t *testing.T) {
	expectedResponse := []byte("Ok")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(expectedResponse)
	}))
	defer ts.Close()

	res, err := getRequest(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(res, expectedResponse) {
		t.Errorf("expected response: %s, result: %s", expectedResponse, res)
	}
}