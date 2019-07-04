package testtool

import (
	"fmt"
	"runtime/debug"
	"testing"
)

// FailOnErr used in testing assert
func FailOnErr(t *testing.T, e error, msg string) {
	if e != nil {
		fmt.Printf("[Fail] on error, %s, %v\n", msg, e)
		debug.PrintStack()
		t.FailNow()
	}
}

// FailOnFlag false时t.Fatal
func FailOnFlag(t *testing.T, flag bool, params ...interface{}) {
	if flag {
		fmt.Printf("[Fail] on flag, %v\n", params)
		debug.PrintStack()
		t.FailNow()
	}
}
