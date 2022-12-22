package phoneNumberGenerate

import (
	"testing"
)

func TestAcquireGenerator(t *testing.T) {
	for ios2 := range regexpData() {
		g := AcquireGenerator(ios2)

		t.Run(ios2, func(t *testing.T) {
			_, err := g.GenerateNumber()
			if err != nil {
				t.Fatal(err)
			}

			_, err = g.GenerateMultipleNumber(3)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
