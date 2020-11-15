package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "joojoojuice",
		Price: 1.00,
		SKU:"abs-abc-abd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
