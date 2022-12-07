package msisdn

import (
	"fmt"
	"testing"
)

func Test_validateNpaCodeFromPhoneNumber(t *testing.T) {
	code := []string{
		"3436771175",
		"6394691699",
		"6394694309",
		"7825099861",
		"5812729417",
		"7827151353",
		"3434491721",
		"3434482434",
		"3434482078",
		"3434482296",
		"2367619340",
		"2268591411",
		"3434482746",
		"7826241546",
		"2268614692",
		"5812820500",
		"2268634519",
		"7826251991",
		"5812709302",
		"3434472966",
		"7826231471",
		"3434492754",
		"3434492426",
		"7825089797",
		"3434472784",
		"3434472662",
		"7826241623",
		"2268594010",
		"6394693623",
		"7826231186",
		"6394694480",
		"2367615225",
		"3434471845",
		"7826231579",
		"7827151620",
		"7826231901",
		"3434472169",
		"3434482623",
		"7825099604",
		"7826241304",
		"6394693042",
		"5812830699",
		"5812850123",
		"2268624666",
		"3434472665",
		"7825521186",
		"2367618969",
		"2367619508",
		"7823659514",
		"7825521359",
		"7826241016",
		"2367617894",
		"2268621960",
		"2268611022",
		"5812850360",
		"3434462098",
		"2367615124",
		"5484779956",
		"5812840804",
		"3436771079",
		"2268591448",
		"7825089796",
		"5797074289",
		"7826231425",
		"5485307307",
		"5812689510",
		"7825099583",
		"3434481918",
		"2268631346",
		"3434492544",
		"7826251265",
		"3434472812",
		"2268614382",
		"2367618147",
		"5797444559",
		"7827151299",
		"2367618252",
		"7823669119",
		"2268634217",
		"7826241776",
		"2367615465",
		"7827151452",
		"5812920621",
		"5797444878",
		"7826241265",
		"7827151285",
		"2367616801",
		"5485809197",
		"7826231059",
		"5812840878",
		"5797074234",
		"2268611997",
		"2268614767",
		"3434492917",
		"2268611722",
		"5812659843",
		"7826251600",
		"2268614716",
		"3434492464",
		"2367618675",
	}

	t.Run("tt.name", func(t *testing.T) {
		for _, v := range code {
			got, err := validateNdcCode(v)
			fmt.Println(got, err)
		}
	})
}