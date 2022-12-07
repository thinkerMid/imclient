package firmware

import (
	"testing"
)

func Test_genData(t *testing.T) {
	t.Run("tt.name", func(t *testing.T) {
		//genData()
		//fetch()
	})
}

func TestNewAppleFirmware(t *testing.T) {
	t.Run("tt.name", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			got := NewAppleFirmware()
			if len(got.production) == 0 {
				t.Fatal(got)
			}
			if len(got.osVersion) == 0 {
				t.Fatal(got)
			}
			if len(got.buildNumber) == 0 {
				t.Fatal(got)
			}
		}
	})
}
