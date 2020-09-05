package ptool

import (
	"runtime/pprof"
)

func GetGoroutineCount() interface{} {
	return GetPProfLookup("goroutine")
}

func GetThreadCreateCount() interface{} {
	return GetPProfLookup("threadcreate")
}

func GetBlockCount() interface{} {
	return GetPProfLookup("block")
}

func GetMutexCount() interface{} {
	return GetPProfLookup("mutex")
}

func GetHeapCount() interface{} {
	return GetPProfLookup("heap")
}

func GetPProfLookup(name string) int {
	p := pprof.Lookup(name)
	return p.Count()
}