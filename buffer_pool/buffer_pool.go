package bufferpool

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	Mutex                   sync.Mutex
	BusyFlag                []bool
	Buffer                  []*bytes.Buffer
	BufferSize              int
	BufferPoolMaxCapacity   int
	BufferPoolEnlargeFactor int
	BufferPoolCurrentSize   int
}

// Enlarge the buffer pool by buffer_pool_enlarge_factor
func (this *BufferPool) enlarge() {
	// Have to make sure that the capacity of the BufferPool WOULD NOT LARGER than the max capacity.
	var enlargeFactor = this.BufferPoolEnlargeFactor
	if this.BufferPoolCurrentSize+this.BufferPoolEnlargeFactor > this.BufferPoolMaxCapacity {
		enlargeFactor = this.BufferPoolMaxCapacity - this.BufferPoolCurrentSize
	}
	for n := 0; n < enlargeFactor; n++ {
		this.BusyFlag = append(this.BusyFlag, false)
		this.Buffer = append(
			this.Buffer,
			bytes.NewBuffer(make([]byte, 0, this.BufferSize*1024)))
	}
	this.BufferPoolCurrentSize = len(this.BusyFlag)
}

func NewBufferPool() *BufferPool {
	var tempPool BufferPool
	// Loading parameters
	tempPool.BufferSize = 40
	tempPool.BufferPoolCurrentSize = 100
	tempPool.BufferPoolMaxCapacity = 2000
	tempPool.BufferPoolEnlargeFactor = 100
	// Buffer pool initialization
	tempPool.BusyFlag = make([]bool, 0, tempPool.BufferPoolCurrentSize)
	tempPool.Buffer = make([]*bytes.Buffer, 0, tempPool.BufferPoolCurrentSize)
	for i := 0; i < tempPool.BufferPoolCurrentSize; i++ {
		tempPool.BusyFlag = append(tempPool.BusyFlag, false)
		tempPool.Buffer = append(
			tempPool.Buffer,
			bytes.NewBuffer(make([]byte, 0, tempPool.BufferSize*1024)))
	}
	return &tempPool
}

/*
Get fetch a single buffer from the buffer pool and returns:
	1. The pointer of a bytes.Buffer;
	2. The buffer index of the buffer pool.
Notice:	Caller has to keep the second return value for purpose of the Release(i)
*/
func (this *BufferPool) Get() (*bytes.Buffer, int) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	for i := 0; i < this.BufferPoolCurrentSize; i++ {
		if !this.BusyFlag[i] {
			this.BusyFlag[i] = true
			return this.Buffer[i], i
		}
	}
	if this.BufferPoolCurrentSize < this.BufferPoolMaxCapacity {
		tempPointer := this.BufferPoolCurrentSize
		this.enlarge()
		this.BusyFlag[tempPointer] = true
		return this.Buffer[tempPointer], tempPointer
	}
	return bytes.NewBuffer(make([]byte, 0, this.BufferSize*1024)), this.BufferPoolMaxCapacity * 2
}

/*
Release put a single buffer back into the buffer pool.
Notice:	Entrance parameter of Release() is the SECOND return value of Get()	not the FIRST one.
*/
func (this *BufferPool) Release(index int) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	// Boundary conditions
	if index < 0 ||
		(index >= this.BufferPoolCurrentSize && index < this.BufferPoolMaxCapacity) {
		return
	} else if index >= this.BufferPoolMaxCapacity {
		return
	}
	// Request to release an empty buffer
	if !this.BusyFlag[index] {
		return
	}
	// Release
	this.BusyFlag[index] = false
	this.Buffer[index].Reset()
}
