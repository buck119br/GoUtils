package utils

import (
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSortedSet(t *testing.T) {

	test_key := make([]float64, 100)
	for i := 0; i < 100; i++ {
		test_key[i] = rand.Float64()
	}

	convey.Convey("Sorted Set test :", t, func() {
		NormalLog(test_key)

		convey.Convey("Sorted Set INT test:", func() {

			var ok bool
			test_map := NewFloatSortedMap()
			//test_map := NewSortedMap()

			// Insert
			for i := 0; i < 50; i++ {
				ok = test_map.Insert(test_key[i], float64(i))
				if !ok {
					break
				}
			}
			convey.So(ok, convey.ShouldBeTrue)

			// Sort
			test_map.Sort()

			// Find
			NormalLog(test_map.Key_)

			for _, v := range test_map.Key_ {
				_, ok = test_map.Find(v)
				if !ok {
					break
				}
			}
			convey.So(ok, convey.ShouldBeTrue)

			ok = test_map.Insert(test_key[9], float64(9))
			convey.So(ok, convey.ShouldBeFalse)

			// Update
			ok = test_map.Update(123.0, 123.0)
			convey.So(ok, convey.ShouldBeFalse)

			ok = true
			for i := 0; i < 50; i++ {
				temp_value_before, _ := test_map.Find(test_key[i])
				test_map.Update(test_key[i], test_key[i+50])
				temp_value_after, _ := test_map.Find(test_key[i])
				if temp_value_before == temp_value_after {
					ok = false
				}
			}
			convey.So(ok, convey.ShouldBeTrue)

			// Delete
			test_map.Delete(123.0)

			for _, v := range test_map.Key_ {
				test_map.Delete(v)
			}
		})
	})
}
