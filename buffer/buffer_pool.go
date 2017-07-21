package buffer

import (
	"bytes"
	"sync"
)

const (
	DefaultBufferCap               = 40
	DefaultBufferPoolInitCap       = 50
	DefaultBufferPoolMaxCap        = 5000
	DefaultBufferPoolEnlargeFactor = 50
)

type Buffer struct {
	*bytes.Buffer
	index int
}

type BufferPool struct {
	mutex             sync.Mutex
	busyFlag          []bool
	buffer            []*Buffer
	bufferCap         int
	poolLen           int
	poolCap           int
	poolMaxCap        int
	poolEnlargeFactor int
}

func (bp *BufferPool) newBuffer(index int) *Buffer {
	b := new(Buffer)
	b.Buffer = bytes.NewBuffer(make([]byte, 0, bp.bufferCap*1024))
	b.index = index
	return b
}

func NewBufferPool(bufferCap, poolInitCap, poolMaxCap, poolEnlargeFactor int) *BufferPool {
	if bufferCap == 0 {
		bufferCap = DefaultBufferCap
	}
	if poolInitCap == 0 {
		poolInitCap = DefaultBufferPoolInitCap
	}
	if poolMaxCap == 0 {
		poolMaxCap = DefaultBufferPoolMaxCap
	}
	if poolEnlargeFactor == 0 {
		poolEnlargeFactor = DefaultBufferPoolEnlargeFactor
	}
	bp := new(BufferPool)
	// Loading parameters
	bp.bufferCap = bufferCap
	bp.poolCap = poolInitCap
	bp.poolMaxCap = poolMaxCap
	bp.poolEnlargeFactor = poolEnlargeFactor
	// Buffer pool initialization
	bp.busyFlag = make([]bool, bp.poolCap)
	bp.buffer = make([]*Buffer, 0, bp.poolCap)
	for i := 0; i < bp.poolCap; i++ {
		bp.buffer = append(bp.buffer, bp.newBuffer(i))
	}
	return bp
}

func (bp *BufferPool) BufferCap() int     { return bp.bufferCap }
func (bp *BufferPool) Len() int           { return bp.poolLen }
func (bp *BufferPool) Cap() int           { return bp.poolCap }
func (bp *BufferPool) MaxCap() int        { return bp.poolMaxCap }
func (bp *BufferPool) EnlargeFactor() int { return bp.poolEnlargeFactor }

// Enlarge the buffer pool by poolEnlargeFactor
func (bp *BufferPool) enlarge() {
	// Have to make sure that the capacity of the BufferPool WOULD NOT LARGER than the max capacity.
	enlargeFactor := bp.poolEnlargeFactor
	if bp.poolCap+bp.poolEnlargeFactor > bp.poolMaxCap {
		enlargeFactor = bp.poolMaxCap - bp.poolCap
	}
	tempFlag := make([]bool, 0, bp.poolCap+enlargeFactor)
	tempBuffer := make([]*Buffer, 0, bp.poolCap+enlargeFactor)
	tempFlag = append(tempFlag, bp.busyFlag...)
	tempBuffer = append(tempBuffer, bp.buffer...)
	for i := 0; i < enlargeFactor; i++ {
		tempFlag = append(tempFlag, false)
		tempBuffer = append(tempBuffer, bp.newBuffer(bp.poolCap+i))
	}
	bp.busyFlag = tempFlag
	bp.buffer = tempBuffer
	// update capacity
	bp.poolCap += enlargeFactor
}

/*
Get fetch a single buffer from the buffer pool and returns:
	1. The pointer of a bytes.Buffer;
	2. The buffer index of the buffer pool.
Notice:	Caller has to keep the second return value for purpose of the Release(i)
*/
func (bp *BufferPool) Get() *Buffer {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	bp.poolLen++
	for i := 0; i < bp.poolCap; i++ {
		if !bp.busyFlag[i] {
			bp.busyFlag[i] = true
			return bp.buffer[i]
		}
	}
	if bp.poolCap < bp.poolMaxCap {
		tempIndex := bp.poolCap
		bp.enlarge()
		bp.busyFlag[tempIndex] = true
		return bp.buffer[tempIndex]
	}
	bp.poolLen--
	return bp.newBuffer(bp.poolMaxCap * 2)
}

/*
Release put a single buffer back into the buffer pool.
Notice:	Entrance parameter of Release() is the SECOND return value of Get()	not the FIRST one.
*/
func (bp *BufferPool) Release(b *Buffer) {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	// Boundary conditions
	if b.index < 0 ||
		(b.index >= bp.poolCap && b.index < bp.poolMaxCap) {
		return
	} else if b.index >= bp.poolMaxCap {
		return
	}
	// Request to release an empty buffer
	if !bp.busyFlag[b.index] {
		return
	}
	// Release
	bp.busyFlag[b.index] = false
	bp.buffer[b.index].Reset()
	bp.poolLen--
}
