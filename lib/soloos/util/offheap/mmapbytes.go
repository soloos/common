package offheap

import (
	"syscall"
	"unsafe"
)

type mmapbytes struct {
	readOff      uintptr
	readBoundary uintptr
	data         []byte
	UData        uintptr
}

func AllocMmapBytes(size int) (mmapbytes, error) {
	var (
		ret mmapbytes
		err error
	)
	prot := syscall.PROT_READ | syscall.PROT_WRITE
	flags := syscall.MAP_ANON | syscall.MAP_PRIVATE

	ret.data, err = syscall.Mmap(-1, 0, size, prot, flags)
	if err != nil {
		panic(err)
	}
	ret.UData = *((*uintptr)((unsafe.Pointer)(&ret.data)))
	ret.readOff = ret.UData
	ret.readBoundary = ret.UData + uintptr(size)
	return ret, err
}
