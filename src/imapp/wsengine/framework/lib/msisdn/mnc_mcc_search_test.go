package msisdn

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("tt.name", func(t *testing.T) {
		got, err := Parse("12292708736")
		if err == nil {
			fmt.Println(got.CC)
			fmt.Println(got.Source)
			fmt.Println(got.GetLanguage())
			fmt.Println(got.GetISO())
		}
		//got, err = Parse("1252535097")
		//if err == nil {
		//	fmt.Println(got.GetLanguage())
		//	fmt.Println(got.GetISO())
		//}
	})
}

func Test_genData(t *testing.T) {
	t.Run("tt.name", func(t *testing.T) {
		//genData()
	})
}
