package wl

import (
	"encoding/binary"
	"io/ioutil"
	"math"
	"os"
	"syscall"
	"unsafe"
)

var order binary.ByteOrder

func init() {
	var x uint32 = 0x01020304
	if *(*byte)(unsafe.Pointer(&x)) == 0x01 {
		order = binary.BigEndian
	} else {
		order = binary.LittleEndian
	}
}

func CreateAnonymousFile(size int) (*os.File, error) {
	template := "wayland-shared"
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		panic("XDG_RUNTIME_DIR not defined.")
	}
	ret, err := ioutil.TempFile(dir, template)
	if err != nil {
		return nil, err
	}
	err = syscall.Ftruncate(int(ret.Fd()), int64(size))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func fixedToFloat64(fixed int32) float64 {
	dat := ((int64(1023 + 44)) << 52) + (1 << 51) + int64(fixed)
	return math.Float64frombits(uint64(dat)) - float64(3<<43)
}

func float64ToFixed(v float64) int32 {
	dat := v + float64(int64(3)<<(51-8))
	return int32(math.Float64bits(dat))
}
