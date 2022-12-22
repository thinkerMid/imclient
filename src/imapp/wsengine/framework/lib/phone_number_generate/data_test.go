package phoneNumberGenerate

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/andybalholm/brotli"
	"os"
	"strings"
	"testing"
	"ws/framework/plugin/json"
)

func Test_genData(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"generate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := os.Open("./countries_v2.tsv")
			sc := bufio.NewScanner(f)

			m := make(map[string]regexpPattern)

			for sc.Scan() {
				text := sc.Text()

				splitText := strings.Split(text, "\t")

				iso2 := splitText[0]
				cc := splitText[2]
				pattern := splitText[6]

				m[iso2] = regexpPattern{
					CC:      cc,
					Pattern: pattern,
				}
			}

			b, _ := json.Marshal(m)
			buf := bytes.NewBuffer(make([]byte, 0))
			w := brotli.NewWriterLevel(buf, brotli.BestCompression)
			w.Write(b)
			w.Flush()

			fmt.Println(buf.Bytes())
		})
	}
}
