package utils

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

// 100000 set values, and 100000 read test
var times int = 100000

// BenchmarkTestConcurrentMap
func BenchmarkTestConcurrentMap(b *testing.B) {
	for k := 0; k < b.N; k++ {
		b.StopTimer()
		// generate 100000 keys map
		testKV := map[string]int{}
		for i := 0; i < 10000; i++ {
			testKV[strconv.Itoa(i)] = i
		}

		// new ConcurrentMap
		cMap := NewConcurrentMap()

		// set key and value to cMap
		for k, v := range testKV {
			cMap.Set(k, v)
		}

		// start timer
		b.StartTimer()

		wg := sync.WaitGroup{}
		wg.Add(2)

		// set values
		go func() {
			// rand key set value 100000 counts
			for i := 0; i < times; i++ {
				index := rand.Intn(times)
				cMap.Set(strconv.Itoa(index), index+1)
			}
			wg.Done()
		}()

		// read values
		go func() {
			// rand key, read 100000 counts
			for i := 0; i < times; i++ {
				index := rand.Intn(times)
				cMap.Get(strconv.Itoa(index))
			}
			wg.Done()
		}()

		// wait goroutine done
		wg.Wait()
	}
}
