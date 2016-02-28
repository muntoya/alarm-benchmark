package main

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lhiredis
#include <hiredis/hiredis.h>
*/
import (
	"C"
)

func main() {
	var _ *C.redisContext = C.redisConnect(C.CString("127.0.0.1"), 6379)
}