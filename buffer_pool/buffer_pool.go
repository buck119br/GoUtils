package utils

import (
	"bytes"
	"sync"

	"github.com/astaxie/beego/config"
)

type BufferPool struct {
	mutex_                      sync.Mutex
	busy_flag_                  []bool
	buffer_                     []*bytes.Buffer
	buffer_size_                int
	buffer_pool_max_capacity_   int
	buffer_pool_enlarge_factor_ int
	buffer_pool_current_size_   int
}

// Enlarge the buffer pool by buffer_pool_enlarge_factor
func (this *BufferPool) enlarge() {
	CommonLog("Buffer pool needs to be enlarged...")

	// Have to make sure that the capacity of the BufferPool WOULD NOT LARGER
	// than the max capacity.
	var enlarge_factor = this.buffer_pool_enlarge_factor_
	if this.buffer_pool_current_size_+this.buffer_pool_enlarge_factor_ >
		this.buffer_pool_max_capacity_ {
		enlarge_factor =
			this.buffer_pool_max_capacity_ - this.buffer_pool_current_size_
	}

	for n := 0; n < enlarge_factor; n++ {
		this.busy_flag_ = append(this.busy_flag_, false)
		this.buffer_ = append(
			this.buffer_,
			bytes.NewBuffer(make([]byte, 0, this.buffer_size_*1024)))
	}

	this.buffer_pool_current_size_ = len(this.busy_flag_)

	CommonLog(
		"Buffer pool enlarged and length: ",
		this.buffer_pool_current_size_)
}

func NewBufferPool(conf config.Configer) *BufferPool {

	var temp_pool BufferPool

	// Loading parameters
	temp_pool.buffer_size_ = 40
	temp_pool.buffer_pool_current_size_ = 100
	temp_pool.buffer_pool_max_capacity_ = 2000
	temp_pool.buffer_pool_enlarge_factor_ = 100

	// Buffer pool initialization
	temp_pool.busy_flag_ = make([]bool, 0, temp_pool.buffer_pool_current_size_)
	temp_pool.buffer_ = make(
		[]*bytes.Buffer,
		0,
		temp_pool.buffer_pool_current_size_)

	for i := 0; i < temp_pool.buffer_pool_current_size_; i++ {
		temp_pool.busy_flag_ = append(temp_pool.busy_flag_, false)
		temp_pool.buffer_ = append(
			temp_pool.buffer_,
			bytes.NewBuffer(make([]byte, 0, temp_pool.buffer_size_*1024)))
	}

	CommonLog(
		"Buffer Pool initialization finished with capacity: ",
		len(temp_pool.busy_flag_))

	return &temp_pool
}

/*
Get fetch a single buffer from the buffer pool and returns:
	1. The pointer of a bytes.Buffer;
	2. The buffer index of the buffer pool.
Notice:
	Caller has to keep the second return value for purpose of the Release(i)
*/
func (this *BufferPool) Get() (*bytes.Buffer, int) {

	this.mutex_.Lock()
	defer this.mutex_.Unlock()

	DebugLog("Buffer GET: request.")
	for i := 0; i < this.buffer_pool_current_size_; i++ {
		if !this.busy_flag_[i] {
			this.busy_flag_[i] = true
			DebugfLog("Buffer GET: NO.%d succeeded.", i)
			return this.buffer_[i], i
		}
	}

	if this.buffer_pool_current_size_ < this.buffer_pool_max_capacity_ {
		temp_pointer := this.buffer_pool_current_size_
		this.enlarge()

		this.busy_flag_[temp_pointer] = true
		DebugfLog("Buffer GET: NO.%d succeeded.", temp_pointer)
		return this.buffer_[temp_pointer], temp_pointer
	}

	NormalLog("Warning !!! Buffer pool is full!")
	DebugLog("Buffer GET: temp buffer succeeded.")
	return bytes.NewBuffer(make([]byte, 0, this.buffer_size_*1024)),
		this.buffer_pool_max_capacity_ * 2
}

/*
Release put a single buffer back into the buffer pool.

Notice:
	Entrance parameter of Release() is the SECOND return value of Get()
	not the FIRST one.
*/
func (this *BufferPool) Release(index int) {

	this.mutex_.Lock()
	defer this.mutex_.Unlock()
	DebugLog("Buffer RELEASE: request.")

	// Boundary conditions
	if index < 0 ||
		(index >= this.buffer_pool_current_size_ &&
			index < this.buffer_pool_max_capacity_) {
		ErrorfLog("Buffer RELEASE: index{%d} outrange.", index)
		return
	} else if index >= this.buffer_pool_max_capacity_ {
		NormalfLog("Buffer RELEASE: temp buffer release : %d", index)
		return
	}

	// Request to release an empty buffer
	if !this.busy_flag_[index] {
		NormalfLog("Buffer RELEASE: buffer{%d} already free.", index)
		return
	}

	// Release
	this.busy_flag_[index] = false
	this.buffer_[index].Reset()
	DebugfLog("Buffer RELEASE: NO.%d succeeded.", index)
	return
}
