package bytePool

import "sync"

// Alloc .
func Alloc(size int) []byte {
	if size <= MaxSize {
		sz := uint32(size)
		for i, v := range sizes {
			if sz <= v {
				bs := pools[i].Get().([]byte)
				return bs[0:size]
			}
		}
	}
	return nil
}

// Free .
func Free(bs []byte) {
	bs = bs[0:cap(bs)]
	l := uint32(len(bs))
	for i, v := range sizes {
		if l == v {
			pools[i].Put(bs)
			return
		}
	}
}

const stage = 18

// MaxSize .
const MaxSize = 2 << 23

var pools = [stage]sync.Pool{}
var sizes = [stage]uint32{}

func init() {
	var byteSize uint32

	for i := range pools {
		byteSize = 2 << (uint32(i) + 5)

		pools[i].New = create(byteSize).create
		sizes[i] = byteSize
	}
}

type create int32

func (c create) create() interface{} {
	return make([]byte, int(c))
}
