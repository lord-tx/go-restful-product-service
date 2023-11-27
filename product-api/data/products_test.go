package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := Product{
		Name:  "Beans",
		Price: 42,
		SKU:   "abs-abc-def",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
