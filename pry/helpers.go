package pry

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// InterpretError is an error returned by the interpreter and shouldn't be
// passed to the user or running code.
type InterpretError struct {
	err error
}

func (a *InterpretError) Error() error {
	if a == nil {
		return nil
	}
	return a.err
}

// Append is a runtime replacement for the append function
func Append(arr interface{}, elems ...interface{}) (interface{}, *InterpretError) {
	arrVal := reflect.ValueOf(arr)
	valArr := make([]reflect.Value, len(elems))
	for i, elem := range elems {
		if reflect.TypeOf(arr) != reflect.SliceOf(reflect.TypeOf(elem)) {
			return nil, &InterpretError{fmt.Errorf("%T cannot append to %T", elem, arr)}
		}
		valArr[i] = reflect.ValueOf(elem)
	}
	return reflect.Append(arrVal, valArr...).Interface(), nil
}

// Make is a runtime replacement for the make function
func Make(t interface{}, args ...interface{}) (interface{}, *InterpretError) {
	typ, isType := t.(reflect.Type)
	if !isType {
		return nil, &InterpretError{fmt.Errorf("invalid type %#v", t)}
	}
	switch typ.Kind() {
	case reflect.Slice:
		if len(args) < 1 || len(args) > 2 {
			return nil, &InterpretError{errors.New("invalid number of arguments. Missing len or extra?")}
		}
		length, isInt := args[0].(int)
		if !isInt {
			return nil, &InterpretError{errors.New("len is not int")}
		}
		capacity := length
		if len(args) == 2 {
			capacity, isInt = args[0].(int)
			if !isInt {
				return nil, &InterpretError{errors.New("len is not int")}
			}
		}
		if length < 0 || capacity < 0 {
			return nil, &InterpretError{errors.Errorf("negative length or capacity")}
		}
		slice := reflect.MakeSlice(typ, length, capacity)
		return slice.Interface(), nil

	case reflect.Chan:
		if len(args) > 1 {
			fmt.Printf("CHAN ARGS %#v", args)
			return nil, &InterpretError{errors.New("too many arguments")}
		}
		size := 0
		if len(args) == 1 {
			var isInt bool
			size, isInt = args[0].(int)
			if !isInt {
				return nil, &InterpretError{errors.New("size is not int")}
			}
		}
		if size < 0 {
			return nil, &InterpretError{errors.Errorf("negative buffer size")}
		}
		buffer := reflect.MakeChan(typ, size)
		return buffer.Interface(), nil

	default:
		return nil, &InterpretError{fmt.Errorf("unknown kind type %T", t)}
	}
}

// Close is a runtime replacement for the "close" function.
func Close(t interface{}) (interface{}, *InterpretError) {
	reflect.ValueOf(t).Close()
	return nil, nil
}

// Len is a runtime replacement for the len function
func Len(t interface{}) (interface{}, *InterpretError) {
	return reflect.ValueOf(t).Len(), nil
}
