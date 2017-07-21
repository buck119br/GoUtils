package buffer

import (
	"bytes"
)

type BufferQueue struct {
	queue []*bytes.Buffer
	len   int
	cap   int
	head  int
	tail  int
}

func NewBufferQueue(cap int) *BufferQueue {
	bq := new(BufferQueue)
	bq.queue = make([]*bytes.Buffer, cap)
	bq.len = 0
	bq.cap = cap
	bq.head = 0
	bq.tail = 0
	return bq
}

func (bq *BufferQueue) Len() int  { return bq.len }
func (bq *BufferQueue) Cap() int  { return bq.cap }
func (bq *BufferQueue) Head() int { return bq.head }
func (bq *BufferQueue) Tail() int { return bq.tail }

func (bq *BufferQueue) Enqueue(buffer *bytes.Buffer) bool {
	if bq.len == bq.cap {
		return false
	}
	bq.queue[bq.tail] = buffer
	if bq.tail == bq.cap-1 {
		bq.tail = 0
	} else {
		bq.tail++
	}
	bq.len++
	return true
}

func (bq *BufferQueue) Dequeue() (buffer *bytes.Buffer) {
	if bq.len == 0 {
		return nil
	}
	buffer = bq.queue[bq.head]
	if bq.head == bq.cap-1 {
		bq.head = 0
	} else {
		bq.head++
	}
	bq.len--
	return buffer
}
