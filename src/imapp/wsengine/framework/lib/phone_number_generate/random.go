package phoneNumberGenerate

import (
	"math"
)

func (g *Generator) randRune(runeGroup []rune) rune {
	r := g.randInt(len(runeGroup) / 2)
	return rune(
		g.randInRange(
			int(runeGroup[2*r]),
			int(runeGroup[2*r+1])))
}

func (g *Generator) randAscii() rune {
	return rune(' ' + g.rand.Intn(0x5f))
}

func (g *Generator) randInRange(min, max int) int {
	return min + g.rand.Intn(max-min+1)
}

const defaultMax = math.MaxInt8 >> 1

func (g *Generator) randRepeat(min, max int) int {
	if max < 0 {
		max = min + defaultMax
	}
	return g.randInRange(min, max)
}

const halfMaxInt16 = math.MaxInt16 >> 1

func (g *Generator) randBool() bool {
	return g.rand.Intn(math.MaxInt16) > halfMaxInt16
}

func (g *Generator) randInt(max int) int {
	if max <= 0 {
		return 0
	}
	return g.rand.Intn(max)
}
