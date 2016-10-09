package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

/*
	Generic sorted map map[interface{}]interface{}.

*/
type KeySlice []interface{}

type SortedMap struct {
	Key_ KeySlice
	map_ map[string]interface{}
}

func TypeSwitch(input interface{}) string {
	switch input := input.(type) {
	case nil:
		return "NULL"
	case bool:
		if input {
			return "TRUE"
		}
		return "FALSE"
	case int:
		return fmt.Sprintf("%d", input)
	case float32, float64:
		return fmt.Sprintf("%g", input)
	case string:
		return input
	default:
		panic(fmt.Sprintf("unexpected type %T: %v", input, input))
	}
}

func (this KeySlice) Len() int {
	return len(this)
}

func (this KeySlice) Less(i, j int) bool {
	former := TypeSwitch(this[i])
	latter := TypeSwitch(this[j])

	switch this[i].(type) {
	case int:
		former_value, _ := strconv.ParseInt(former, 10, 32)
		latter_value, _ := strconv.ParseInt(latter, 10, 32)
		return former_value < latter_value
	case float32, float64:
		former_value, _ := strconv.ParseFloat(former, 64)
		latter_value, _ := strconv.ParseFloat(latter, 64)
		return former_value < latter_value
	case string:
		if strings.Compare(former, latter) == -1 {
			return true
		}
		return false
	}
	return false
}

func (this KeySlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func NewSortedMap() *SortedMap {
	var temp_map SortedMap

	temp_map.Key_ = make([]interface{}, 0, 20)
	temp_map.map_ = make(map[string]interface{}, 20)

	return &temp_map
}

func (this *SortedMap) KeySliceToInt() *[]int {
	temp_slice := make([]int, 0, len(this.Key_))
	for _, v := range this.Key_ {
		temp_value, _ := strconv.ParseInt(TypeSwitch(v), 10, 32)
		temp_slice = append(temp_slice, int(temp_value))
	}
	return &temp_slice
}

func (this *SortedMap) KeySliceToFloat() *[]float64 {
	temp_slice := make([]float64, 0, len(this.Key_))
	for _, v := range this.Key_ {
		temp_value, _ := strconv.ParseFloat(TypeSwitch(v), 64)
		temp_slice = append(temp_slice, temp_value)
	}
	return &temp_slice
}

func (this *SortedMap) KeySliceToString() *[]string {
	temp_slice := make([]string, 0, len(this.Key_))
	for _, v := range this.Key_ {
		temp_value := TypeSwitch(v)
		temp_slice = append(temp_slice, temp_value)
	}
	return &temp_slice
}

func (this *SortedMap) Sort() {
	sort.Sort(this.Key_)
}

func (this *SortedMap) Insert(key interface{}, value interface{}) bool {
	temp_key := TypeSwitch(key)

	_, ok := this.map_[temp_key]
	if ok {
		return false
	}
	this.Key_ = append(this.Key_, key)
	this.map_[temp_key] = value
	return true
}

func (this *SortedMap) Find(key interface{}) (interface{}, bool) {
	v, ok := this.map_[TypeSwitch(key)]
	return v, ok
}

func (this *SortedMap) Update(key interface{}, value interface{}) bool {
	temp_key := TypeSwitch(key)
	_, ok := this.map_[temp_key]
	if !ok {
		return ok
	}
	this.map_[temp_key] = value
	return ok
}

func (this *SortedMap) Delete(key interface{}) {
	temp_key := TypeSwitch(key)
	if _, ok := this.map_[temp_key]; !ok {
		return
	}
	temp_key_slice := make(KeySlice, 0, 20)
	for _, v := range this.Key_ {
		if v != key {
			temp_key_slice = append(temp_key_slice, v)
		}
	}
	this.Key_ = temp_key_slice
	delete(this.map_, temp_key)
}
