package bnm

import (
	"testing"
)

func TestResponse_unmarshalResponse_Err(t *testing.T) {
	_, err := unmarshalResponse([]byte("invalid json"))

	if err == nil {
		t.Error("unmarshaling of an invalid error should return an error")
	}
}

func TestResponse_unmarshalResponse_Ok(t *testing.T) {
	xmlData := `
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

	res, err := unmarshalResponse([]byte(xmlData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Currencies) != 3 {
		t.Fatalf("expected 3 currencies, got %d", len(res.Currencies))
	}

	tests := []struct {
		code    string
		name    string
		value   float32
		nominal int
		numCode int
	}{
		{"EUR", "Euro", 21.2997, 1, 978},
		{"USD", "Dolar S.U.A.", 17.9948, 1, 840},
		{"RON", "Leu romanesc", 4.6680, 1, 946},
	}

	for _, tt := range tests {
		curr, ok := res.FindByCode(tt.code)
		if !ok {
			t.Errorf("currency %s not found", tt.code)
			continue
		}
		if curr.Name != tt.name || curr.Nominal != tt.nominal || curr.NumCode != tt.numCode || curr.Value != tt.value {
			t.Errorf("currency %s: expected %+v, got %+v", tt.code, tt, curr)
		}
	}

	if res.Date != "05.08.2017" {
		t.Errorf("expected date '05.08.2017', got '%s'", res.Date)
	}

	if res.Name != "Cursul oficial de schimb" {
		t.Errorf("expected name 'Cursul oficial de schimb', got '%s'", res.Name)
	}
}
