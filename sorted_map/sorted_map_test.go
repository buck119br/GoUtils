package utils

import (
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSortedSet(t *testing.T) {

	testKey := make([]float64, 100)
	for i := 0; i < 100; i++ {
		testKey[i] = rand.Float64()
	}

	convey.Convey("Sorted Set test :", t, func() {
		convey.Convey("Sorted Set INT test:", func() {

			var ok bool
			testMap := NewSortedMap()

			// Insert
			for i := 0; i < 50; i++ {
				ok = testMap.Insert(testKey[i], float64(i))
				if !ok {
					break
				}
			}
			convey.So(ok, convey.ShouldBeTrue)

			// Sort
			testMap.Sort()

			// Find
			for _, v := range testMap.Key {
				_, ok = testMap.Find(v)
				if !ok {
					break
				}
			}
			convey.So(ok, convey.ShouldBeTrue)

			ok = testMap.Insert(testKey[9], float64(9))
			convey.So(ok, convey.ShouldBeFalse)

			// Update
			ok = testMap.Update(123.0, 123.0)
			convey.So(ok, convey.ShouldBeFalse)

			ok = true
			for i := 0; i < 50; i++ {
				tempValueBefore, _ := testMap.Find(testKey[i])
				testMap.Update(testKey[i], testKey[i+50])
				tempValueAfter, _ := testMap.Find(testKey[i])
				if tempValueBefore == tempValueAfter {
					ok = false
				}
			}
			convey.So(ok, convey.ShouldBeTrue)

			// Delete
			testMap.Delete(123.0)

			for _, v := range testMap.Key {
				testMap.Delete(v)
			}
		})
	})
}
