package bufferpool

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type testUnit struct {
	reqNum      int
	bufferIndex int
	buffer      *bytes.Buffer
	writeLength int
	writeErr    error
	md5         [16]byte
}

func TestBufferPool(t *testing.T) {

	var chanCounter int
	testChannel := make(chan testUnit, 30)
	defer close(testChannel)

	convey.Convey("Buffer Pool Test: ", t, func() {

		testPool := NewBufferPool()

		convey.Convey("Buffer pool serial test:", func() {

			// For purpose of boundary condition test,
			// have to allocate the first half at the beginning
			for i := 0; i < testPool.BufferPoolMaxCapacity/2; i++ {
				testPool.Get()
			}
			testPool.Release(-1)
			testPool.Release(testPool.BufferPoolMaxCapacity - 1)
			for i := testPool.BufferPoolMaxCapacity / 2; i < testPool.BufferPoolMaxCapacity+5; i++ {
				testPool.Get()
			}
			for i := 0; i < testPool.BufferPoolMaxCapacity+5; i++ {
				testPool.Release(i)
			}
			testPool.Release(0)
			testPool.Release(testPool.BufferPoolMaxCapacity - 1)
		})

		convey.Convey("Buffer pool release test: should pass.", func() {

			for i := 0; i < 10; i++ {
				intToBytes := make([]byte, 4)
				tempBuffer := make([]byte, 0, testPool.BufferSize*1024)
				binary.BigEndian.PutUint32(intToBytes, uint32(i))

				for x := 0; x < testPool.BufferSize*256; x++ {
					tempBuffer = append(tempBuffer, intToBytes...)
				}
				md5temp := md5.Sum(tempBuffer)

				buffer, bufferIndex := testPool.Get()
				buffer.Write(tempBuffer)

				tempBufferRead := make([]byte, testPool.BufferSize*1024)

				n, err := buffer.Read(tempBufferRead)
				convey.So(err, convey.ShouldBeNil)
				convey.So(n, convey.ShouldEqual, testPool.BufferSize*1024)
				testPool.Release(bufferIndex)
				convey.So(
					md5temp,
					convey.ShouldEqual,
					md5.Sum(tempBufferRead))
			}
		})

		convey.Convey("Buffer pool concurrent test: should pass.", func() {

			for i := 0; i < testPool.BufferPoolMaxCapacity; i++ {
				go func(count int) {
					var sendUnit testUnit

					intToBytes := make([]byte, 4)
					tempBufferWrite := make(
						[]byte,
						0,
						testPool.BufferSize*1024)
					binary.BigEndian.PutUint32(intToBytes, uint32(count))

					for x := 0; x < testPool.BufferSize*256; x++ {
						tempBufferWrite = append(
							tempBufferWrite,
							intToBytes...)
					}
					sendUnit.md5 = md5.Sum(tempBufferWrite)
					sendUnit.reqNum = count
					sendUnit.buffer, sendUnit.bufferIndex =
						testPool.Get()
					sendUnit.writeLength, sendUnit.writeErr =
						sendUnit.buffer.Write(tempBufferWrite)
					testChannel <- sendUnit
				}(i)
			}

			for receiveUnit := range testChannel {
				tempBufferRead := make([]byte, testPool.BufferSize*1024)

				convey.So(receiveUnit.writeErr, convey.ShouldBeNil)
				convey.So(
					receiveUnit.writeLength,
					convey.ShouldEqual,
					testPool.BufferSize*1024)

				n, err := receiveUnit.buffer.Read(tempBufferRead)
				convey.So(err, convey.ShouldBeNil)
				convey.So(n, convey.ShouldEqual, testPool.BufferSize*1024)

				md5temp := md5.Sum(tempBufferRead)
				testPool.Release(receiveUnit.bufferIndex)
				convey.So(md5temp, convey.ShouldEqual, receiveUnit.md5)

				chanCounter++
				if chanCounter == testPool.BufferPoolMaxCapacity {
					break
				}
			}
		})
	})
}

// For purpose of -benchmem
func BenchmarkBufferPool(b *testing.B) {
	testPool := NewBufferPool()

	for i := 0; i < b.N; i++ {
		testBuffer, index := testPool.Get()
		intToBytes := make([]byte, 4)
		tempBufferWrite := make([]byte, 0, testPool.BufferSize*1024)
		binary.BigEndian.PutUint32(intToBytes, uint32(i))

		for x := 0; x < testPool.BufferSize*256; x++ {
			tempBufferWrite = append(tempBufferWrite, intToBytes...)
		}
		testBuffer.Write(tempBufferWrite)
		testPool.Release(index)
	}
}
