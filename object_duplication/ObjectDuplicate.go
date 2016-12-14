package main

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type tiny_struct struct {
	Int_    int
	String_ string
}

type test_struct struct {
	Bool          bool
	Int           int
	Int8          int8
	Int16         int16
	Int32         int32
	Int64         int64
	Uint          uint
	Uint8         uint8
	Uint16        uint16
	Uint32        uint32
	Uint64        uint64
	Uintptr       uintptr
	Float32       float32
	Float64       float64
	Complex64     complex64
	Complex128    complex128
	Array         [10]int
	Chan          chan int
	Interface     interface{}
	Map           map[string]float32
	Ptr           *int
	Slice         []int64
	String        string
	Struct        tiny_struct
	UnsafePointer unsafe.Pointer
	// Func func
}

func main() {
	var (
		input test_struct
	)

	input.Bool = true
	input.Int = 123
	input.Array = [10]int{1, 2, 3, 4, 5, 6}
	input.Chan = make(chan int, 3)
	input.Map = make(map[string]float32, 10)
	input.Map["1"] = 111.11
	input.Map["2"] = 222.22
	input.Map["3"] = 333.33
	input.Ptr = &input.Int
	input.Slice = []int64{123, 456, 789}
	input.String = "test data"
	input.Struct.Int_ = 123123
	input.Struct.String_ = "tiny_struct"
	input.Interface = new(io.Writer)
	result, err := ObjectDuplicate(input.Struct)
	fmt.Println(err)
	fmt.Println("input type:\t ", reflect.TypeOf(input.Struct))
	fmt.Println("input value:\t ", input.Struct)
	// fmt.Println("input address:\t ", &input.Array[0])
	fmt.Println("output type:\t ", reflect.TypeOf(result))
	fmt.Println("output value:\t ", result)
	fmt.Println("output address:\t ", &result)
}

func ObjectDuplicate(input interface{}) (interface{}, error) {
	input_type := reflect.TypeOf(input)
	input_value := reflect.ValueOf(input)
	switch input_type.Kind() {
	case reflect.Invalid:
		return nil,
			errors.New(fmt.Sprintf(
				"ObjectiveDuplicate Error: Invalid input type, IN:[%T]",
				input))
	case reflect.Bool:
		return input_value.Bool(), nil
	case reflect.Int:
		return int(input_value.Int()), nil
	case reflect.Int8:
		return int8(input_value.Int()), nil
	case reflect.Int16:
		return int16(input_value.Int()), nil
	case reflect.Int32:
		return int32(input_value.Int()), nil
	case reflect.Int64:
		return input_value.Int(), nil
	case reflect.Uint:
		return uint(input_value.Uint()), nil
	case reflect.Uint8:
		return uint8(input_value.Uint()), nil
	case reflect.Uint16:
		return uint16(input_value.Uint()), nil
	case reflect.Uint32:
		return uint32(input_value.Uint()), nil
	case reflect.Uint64:
		return input_value.Uint(), nil
	case reflect.Uintptr:
		return input_value.Pointer(), nil
	case reflect.Float32:
		return float32(input_value.Float()), nil
	case reflect.Float64:
		return input_value.Float(), nil
	case reflect.Complex64:
		return complex64(input_value.Complex()), nil
	case reflect.Complex128:
		return input_value.Complex(), nil
	case reflect.Array:
		return input_value, nil
	case reflect.Chan:
		capacity := input_value.Cap()
		temp_chan := reflect.MakeChan(input_type, capacity)
		return temp_chan, nil
	case reflect.Func:
		return nil, nil
	case reflect.Interface:
		return input_value.Interface(), nil
	case reflect.Map:
		keys := input_value.MapKeys()
		temp_map := reflect.MakeMap(input_type)
		for _, v := range keys {
			temp_map.SetMapIndex(v, input_value.MapIndex(v))
		}
		return temp_map, nil
	case reflect.Ptr:
		return input_value.Pointer(), nil
	case reflect.Slice:
		capacity := input_value.Cap()
		temp_slice := reflect.MakeSlice(input_type, 0, capacity)
		temp_slice = reflect.AppendSlice(temp_slice, input_value)
		return temp_slice, nil
	case reflect.String:
		return input_value.String(), nil
	case reflect.Struct:
		temp_struct := reflect.Indirect(reflect.New(input_type))
		num_field := input_type.NumField()
		for i := 0; i < num_field; i++ {
			temp_struct.Field(i).Set(input_value.Field(i))
		}
		return temp_struct, nil
	case reflect.UnsafePointer:
		return input_value.Pointer(), nil
	}
	return nil,
		errors.New(fmt.Sprintf(
			"ObjectiveDuplicate Error: Unknown input type, IN:[%T]",
			input))
}
