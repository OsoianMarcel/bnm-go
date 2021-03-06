package bnm

// Rate model
type Rate struct {
	Code  string  `xml:"CharCode" json:"code"`
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}

// Result used to parse XML
type Result struct {
	Rates []Rate `xml:"Valute" json:"rates"`
}

// FindByCode finds one rate by code
func (r Result) FindByCode(code string) (Rate, bool) {
	for _, val := range r.Rates {
		if val.Code == code {
			return val, true
		}
	}

	return Rate{}, false
}
