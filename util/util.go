package util

import "reflect"

func All(bs ...bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

func Assert(should bool, callback func()) {
	if !should {
		callback()
	}
}

// Contains .
func Contains(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
