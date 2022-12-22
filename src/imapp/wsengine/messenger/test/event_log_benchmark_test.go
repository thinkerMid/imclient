package test

import (
	"testing"
	"ws/framework/application/constant/types"
	appContainer "ws/framework/application/container"
	"ws/framework/application/core/wam"
	eventSerialize "ws/framework/plugin/event_serialize"
)

// ----------------------------------------------------------------------------

/*
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkEventNewFromMap
BenchmarkEventNewFromMap-12          538           2182259 ns/op        12911659 B/op       2258 allocs/op
BenchmarkEventNewFromMap-12          544           2178372 ns/op        12911663 B/op       2258 allocs/op
BenchmarkEventNewFromMap-12          560           2213397 ns/op        12911686 B/op       2259 allocs/op
*/
func BenchmarkEventNewFromMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for min := wam.ET_MIN; min <= wam.ET_MAX; min++ {
			wam.NewWAMEvent(min, nil)
		}
	}
}

func BenchmarkEventSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		types.SerializeBuf(types.OPTION_NUMBER, 1336, 3.1, "", 1)
	}
}

func BenchmarkEventBufferSerialize(b *testing.B) {
	eventBuffer := eventSerialize.AcquireEventBuffer()

	for i := 0; i < b.N; i++ {
		eventBuffer.Reset()
		eventBuffer.Header().SerializeNumber(1336, 3.1)
	}
}

func BenchmarkEventBufferGenerate(b *testing.B) {
	ioc := appContainer.NewAppIocContainer()
	ioc.Inject(appContainer.Channel0EventCache, &Cache{Buffer: eventSerialize.AcquireEventBuffer()})
	ioc.Inject(appContainer.Channel2EventCache, &Cache{Buffer: eventSerialize.AcquireEventBuffer()})

	for i := 0; i < b.N; i++ {
		wam.LogManager().LogContactAdd(ioc, true, true)
		wam.LogManager().LogNotifyContactAvatar(ioc, 1024)
		wam.LogManager().LogSessionNew(ioc)
		wam.LogManager().LogSendText(ioc, 0, true)
	}
}

/**
pkg: ws/messenger/test
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkEventSerialize
BenchmarkEventSerialize-12               4344254               272.3 ns/op                128 B/op               10 allocs/op
BenchmarkEventBufferSerialize
BenchmarkEventBufferSerialize-12        38711799               28.80 ns/op                  0 B/op                0 allocs/op
BenchmarkEventGenerate
BenchmarkEventGenerate-12                  87267               13511 ns/op               5192 B/op              370 allocs/op
BenchmarkEventBufferGenerate
BenchmarkEventBufferGenerate-12           444580                2536 ns/op               1770 B/op               16 allocs/op
*/
