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
	Key KeySlice
	Map map[string]interface{}
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
		formerValue, _ := strconv.ParseInt(former, 10, 32)
		latterValue, _ := strconv.ParseInt(latter, 10, 32)
		return formerValue < latterValue
	case float32, float64:
		formerValue, _ := strconv.ParseFloat(former, 64)
		latterValue, _ := strconv.ParseFloat(latter, 64)
		return formerValue < latterValue
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
	var tempMap SortedMap

	tempMap.Key = make([]interface{}, 0, 20)
	tempMap.Map = make(map[string]interface{}, 20)

	return &tempMap
}

func (this *SortedMap) KeySliceToInt() *[]int {
	tempSlice := make([]int, 0, len(this.Key))
	for _, v := range this.Key {
		tempValue, _ := strconv.ParseInt(TypeSwitch(v), 10, 32)
		tempSlice = append(tempSlice, int(tempValue))
	}
	return &tempSlice
}

func (this *SortedMap) KeySliceToFloat() *[]float64 {
	tempSlice := make([]float64, 0, len(this.Key))
	for _, v := range this.Key {
		tempValue, _ := strconv.ParseFloat(TypeSwitch(v), 64)
		tempSlice = append(tempSlice, tempValue)
	}
	return &tempSlice
}

func (this *SortedMap) KeySliceToString() *[]string {
	tempSlice := make([]string, 0, len(this.Key))
	for _, v := range this.Key {
		tempValue := TypeSwitch(v)
		tempSlice = append(tempSlice, tempValue)
	}
	return &tempSlice
}

func (this *SortedMap) Sort() {
	sort.Sort(this.Key)
}

func (this *SortedMap) Insert(key interface{}, value interface{}) bool {
	tempKey := TypeSwitch(key)

	_, ok := this.Map[tempKey]
	if ok {
		return false
	}
	this.Key = append(this.Key, key)
	this.Map[tempKey] = value
	return true
}

func (this *SortedMap) Find(key interface{}) (interface{}, bool) {
	v, ok := this.Map[TypeSwitch(key)]
	return v, ok
}

func (this *SortedMap) Update(key interface{}, value interface{}) bool {
	tempKey := TypeSwitch(key)
	_, ok := this.Map[tempKey]
	if !ok {
		return ok
	}
	this.Map[tempKey] = value
	return ok
}

func (this *SortedMap) Delete(key interface{}) {
	tempKey := TypeSwitch(key)
	if _, ok := this.Map[tempKey]; !ok {
		return
	}
	tempKeyslice := make(KeySlice, 0, 20)
	for _, v := range this.Key {
		if v != key {
			tempKeyslice = append(tempKeyslice, v)
		}
	}
	this.Key = tempKeyslice
	delete(this.Map, tempKey)
}
