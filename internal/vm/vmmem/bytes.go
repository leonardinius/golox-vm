package vmmem

// import (
// 	"bytes"
// 	"sync"
// 	"unsafe"
// )

// // buffer pool to reduce GC.
// var (
// 	gBytes = sync.Pool{
// 		// New is called when a new instance is needed.
// 		New: func() interface{} {
// 			return new(bytes.Buffer)
// 		},
// 	}

// 	gUint64Size = int(unsafe.Sizeof(uint64(0)))
// )

// // mallocBuffer fetches a buffer from the pool.
// func mallocBuffer(size int) *bytes.Buffer {
// 	buffer := gBytes.Get().(*bytes.Buffer)
// 	buffer.Grow(size)
// 	return buffer
// }

// // reclaimBuffer returns a buffer to the pool.
// func freeBuffer(buf *bytes.Buffer) {
// 	buf.Reset()
// 	gBytes.Put(buf)
// }
