package protocol

import (
	"fmt"
	"unsafe"
)

type rawString struct {
	CodeType CodeType
}

func (r *rawString) ToString() string {
	return "String"
}
func (r *rawString) Unmarshal(b []byte, v interface{}) error {
	if vv, ok := v.(*string); ok {
		*vv = *(*string)(unsafe.Pointer(&b))
	}
	return nil
}
func (r *rawString) Marshal(v interface{}) ([]byte, error) {
	if vv, ok := v.(string); ok {
		return []byte(vv), nil
	}
	return nil, fmt.Errorf("v type not string")
}
