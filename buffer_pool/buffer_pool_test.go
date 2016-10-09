package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"flag"
	"testing"

	"github.com/astaxie/beego/config"
	"github.com/smartystreets/goconvey/convey"
)

type testUnit struct {
	req_num_      int
	buffer_index_ int
	buffer_       *bytes.Buffer
	write_length_ int
	write_err_    error
	md5_          [16]byte
}

func TestBufferPool(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")
	// flag.Lookup("log_dir").Value.Set("./")
	flag.Lookup("v").Value.Set("600")
	flag.Parse()
	// defer glog.Flush()

	var chan_counter int
	test_chan := make(chan testUnit, 30)
	defer close(test_chan)

	conf, err := config.NewConfig(
		CONFIG_FILE_PROVIDER,
		ConfDir()+"app.conf")
	if err != nil {
		FatalLog(err)
	}

	convey.Convey("Buffer pool test: ", t, func() {

		test_pool := NewBufferPool(conf)

		convey.Convey("Buffer pool serial test:", func() {

			// For purpose of boundary condition test,
			// have to allocate the first half at the beginning
			for i := 0; i < test_pool.buffer_pool_max_capacity_/2; i++ {
				test_pool.Get()
			}
			test_pool.Release(-1)
			test_pool.Release(test_pool.buffer_pool_max_capacity_ - 1)
			for i := test_pool.buffer_pool_max_capacity_ / 2; i < test_pool.buffer_pool_max_capacity_+5; i++ {
				test_pool.Get()
			}
			for i := 0; i < test_pool.buffer_pool_max_capacity_+5; i++ {
				test_pool.Release(i)
			}
			test_pool.Release(0)
			test_pool.Release(test_pool.buffer_pool_max_capacity_ - 1)
		})

		convey.Convey("Buffer pool release test: should pass.", func() {

			for i := 0; i < 10; i++ {
				intToBytes := make([]byte, 4)
				temp_buffer := make([]byte, 0, test_pool.buffer_size_*1024)
				binary.BigEndian.PutUint32(intToBytes, uint32(i))

				for x := 0; x < test_pool.buffer_size_*256; x++ {
					temp_buffer = append(temp_buffer, intToBytes...)
				}
				md5_temp := md5.Sum(temp_buffer)

				buffer_, buffer_index_ := test_pool.Get()
				buffer_.Write(temp_buffer)

				temp_buffer_read := make([]byte, test_pool.buffer_size_*1024)

				n, err := buffer_.Read(temp_buffer_read)
				convey.So(err, convey.ShouldBeNil)
				convey.So(n, convey.ShouldEqual, test_pool.buffer_size_*1024)
				test_pool.Release(buffer_index_)
				convey.So(
					md5_temp,
					convey.ShouldEqual,
					md5.Sum(temp_buffer_read))
			}
		})

		convey.Convey("Buffer pool concurrent test: should pass.", func() {

			for i := 0; i < test_pool.buffer_pool_max_capacity_; i++ {
				go func(count int) {
					var send_unit testUnit

					intToBytes := make([]byte, 4)
					temp_buffer_write := make(
						[]byte,
						0,
						test_pool.buffer_size_*1024)
					binary.BigEndian.PutUint32(intToBytes, uint32(count))

					for x := 0; x < test_pool.buffer_size_*256; x++ {
						temp_buffer_write = append(
							temp_buffer_write,
							intToBytes...)
					}
					send_unit.md5_ = md5.Sum(temp_buffer_write)
					send_unit.req_num_ = count
					send_unit.buffer_, send_unit.buffer_index_ =
						test_pool.Get()
					send_unit.write_length_, send_unit.write_err_ =
						send_unit.buffer_.Write(temp_buffer_write)
					test_chan <- send_unit
				}(i)
			}

			for receive_unit := range test_chan {
				temp_buffer_read := make([]byte, test_pool.buffer_size_*1024)

				convey.So(receive_unit.write_err_, convey.ShouldBeNil)
				convey.So(
					receive_unit.write_length_,
					convey.ShouldEqual,
					test_pool.buffer_size_*1024)

				n, err := receive_unit.buffer_.Read(temp_buffer_read)
				convey.So(err, convey.ShouldBeNil)
				convey.So(n, convey.ShouldEqual, test_pool.buffer_size_*1024)

				md5_temp := md5.Sum(temp_buffer_read)
				test_pool.Release(receive_unit.buffer_index_)
				convey.So(md5_temp, convey.ShouldEqual, receive_unit.md5_)

				chan_counter++
				if chan_counter == test_pool.buffer_pool_max_capacity_ {
					break
				}
			}
		})
	})
}

// For purpose of -benchmem
func BenchmarkBufferPool(b *testing.B) {
	conf, err := config.NewConfig(
		CONFIG_FILE_PROVIDER,
		ConfDir()+"app.conf")
	if err != nil {
		ErrorLog(err)
	}
	test_pool := NewBufferPool(conf)

	for i := 0; i < b.N; i++ {
		test_buffer, index := test_pool.Get()
		intToBytes := make([]byte, 4)
		temp_buffer_write := make([]byte, 0, test_pool.buffer_size_*1024)
		binary.BigEndian.PutUint32(intToBytes, uint32(i))

		for x := 0; x < test_pool.buffer_size_*256; x++ {
			temp_buffer_write = append(temp_buffer_write, intToBytes...)
		}
		test_buffer.Write(temp_buffer_write)
		test_pool.Release(index)
	}
}
