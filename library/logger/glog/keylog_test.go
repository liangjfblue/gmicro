package glog

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_KeyLogger(t *testing.T) {

	logger := NewKeyLogger("./test_log", time.Second*3, time.Second*10)

	logger.KPrintf("test_key_a", "test key log a\n")

	logger.Flush()

}

func Benchmark_KeyLog(b *testing.B) {

	logger := NewKeyLogger("./test_log", time.Second*2, time.Second*1)
	defer logger.Flush()

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			logger.KPrintf(fmt.Sprintf("test_key_%d", rand.Intn(10)), "test key log a\n")
		}
	})

}

func Test_KeyLogTimeOutClose(t *testing.T) {

	logger := NewKeyLogger("./test_log", time.Second*1, time.Second*1)
	defer logger.Flush()

	logger.KPrint("timeoutlog", "test\n")
	time.Sleep(time.Second * 2)
	logger.KPrint("timeoutlog", "test reenter\n")

}
