package phoneNumberGenerate

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"regexp/syntax"
	"time"
	byteBufferPool "ws/framework/plugin/byte_buffer_pool"
)

// AcquireGenerator .
func AcquireGenerator(iso2 string) *Generator {
	src := rand.NewSource(time.Now().UnixNano())

	g := Generator{
		iso2:   iso2,
		buffer: byteBufferPool.AcquireBuffer(),
		rand:   rand.New(src),
	}

	return &g
}

// ReleaseGenerator .
func ReleaseGenerator(g *Generator) { byteBufferPool.ReleaseBuffer(g.buffer) }

// ----------------------------------------------------------------------------

// Generator 根据正则生成符合规则的号码
type Generator struct {
	iso2   string
	regexp *regexpPattern
	buffer *bytes.Buffer
	rand   *rand.Rand
}

// GenerateNumber .
func (g *Generator) GenerateNumber() (string, error) {
	regexp, err := g.parse()
	if err != nil {
		return "", err
	}

	// 预先写入国家号码
	g.buffer.WriteString(g.regexp.CC)

	// 生成号码
	err = g.generate(regexp)
	if err != nil {
		return "", err
	}

	return g.buffer.String(), err
}

// GenerateMultipleNumber .
func (g *Generator) GenerateMultipleNumber(amount int64) ([]string, error) {
	if amount < 1 {
		return nil, errors.New("amount must be greater than or equal to 1")
	}

	regexp, err := g.parse()
	if err != nil {
		return nil, err
	}

	result := make([]string, amount)
	var i int64 = 0

	for ; i < amount; i++ {
		g.buffer.Reset()

		// 预先写入国家号码
		g.buffer.WriteString(g.regexp.CC)

		// 生成号码
		if err = g.generate(regexp); err != nil {
			return nil, err
		}

		result[i] = g.buffer.String()
	}

	return result, nil
}

func (g *Generator) parse() (*syntax.Regexp, error) {
	var ok bool

	g.regexp, ok = searchRegexp(g.iso2)
	if !ok {
		return nil, fmt.Errorf("not found iso2=%s phone number pattern", g.iso2)
	}

	reg, err := syntax.Parse(g.regexp.Pattern, syntax.Perl)
	if err != nil {
		return nil, err
	}

	return reg.Simplify(), nil
}

func (g *Generator) generate(regexp *syntax.Regexp) error {
	switch regexp.Op {
	case syntax.OpNoMatch,
		syntax.OpEmptyMatch,
		syntax.OpNoWordBoundary,
		syntax.OpBeginLine,
		syntax.OpBeginText,
		syntax.OpEndText:

	case syntax.OpLiteral:
		g.buffer.WriteString(string(regexp.Rune))

	case syntax.OpCharClass:
		if len(regexp.Rune)%2 != 0 || len(regexp.Rune) == 0 {
			return nil
		}
		g.buffer.WriteRune(g.randRune(regexp.Rune))

	case syntax.OpAnyCharNotNL:
		g.buffer.WriteRune(g.randAscii())

	case syntax.OpAnyChar:
		r := g.randInRange(32, 128)
		if r == 127 {
			g.buffer.WriteRune('\n')
		} else {
			g.buffer.WriteRune(rune(r))
		}
	case syntax.OpEndLine:
		g.buffer.WriteRune('\n')

	case syntax.OpWordBoundary:
		g.buffer.WriteRune(0x20)

	case syntax.OpStar:
		return g.repeat(regexp.Sub, 0, -1)

	case syntax.OpPlus:
		return g.repeat(regexp.Sub, 1, -1)

	case syntax.OpQuest:
		if g.randBool() {
			if err := g.traversal(regexp.Sub); err != nil {
				return err
			}
		}
	case syntax.OpRepeat:
		return g.repeat(regexp.Sub, regexp.Min, regexp.Max)

	case syntax.OpConcat, syntax.OpCapture:
		return g.traversal(regexp.Sub)

	case syntax.OpAlternate:
		return g.generate(regexp.Sub[g.randInt(len(regexp.Sub))])
	}

	return nil
}

func (g *Generator) traversal(children []*syntax.Regexp) error {
	for _, child := range children {
		if err := g.generate(child); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) repeat(list []*syntax.Regexp, min, max int) error {
	for count := g.randRepeat(min, max); count > 0; count-- {
		if err := g.traversal(list); err != nil {
			return err
		}
	}
	return nil
}
