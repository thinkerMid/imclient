package lightning

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const (
	RunTimes           = 1000000
	BenchParam         = 10
	BenchAntsSize      = 500
	DefaultExpiredTime = 10 * time.Second
)

var atotal int64 = 0
var total = 0

var lock = sync.RWMutex{}

func antsDemoFunc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

func demoFunc(_ interface{}) {
	//lock.Lock()
	//defer lock.Unlock()
	//total += i.(int)
	//atomic.AddInt64(&atotal, 1)
	//fmt.Println(i)
	//if total == RunTimes {
	//}
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

func BenchmarkSemaphoreWorkerPoolThroughput(b *testing.B) {
	p := New(BenchAntsSize, 1024, demoFunc)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			p.Invoke(j, false)
		}
	}
	b.StopTimer()
}

func BenchmarkSemaphoreWorkerPool(b *testing.B) {
	p := New(BenchAntsSize, 1024, demoFunc)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			p.Invoke(j, false)
		}
	}
	p.Await()
	b.StopTimer()
}

func TestNew(t *testing.T) {
	p := New(BenchAntsSize, 1024, demoFunc)
	//for i := 0; i < b.N; i++ {
	for j := 0; j < RunTimes; j++ {
		//fmt.Println(j)
		p.Invoke(j, false)
	}
	//}
	p.Await()
	fmt.Println(total)
}
