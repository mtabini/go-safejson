// Package safejson provides a facility to selectively export data from Go to JSON. It's
// designed for scenarios in which it is desirable to only allow certain data to be
// exported during a marshaling operation in a failsafe manner.
//
// safejson introduces a new tag, called `safejson`, that must be explicitly present
// in order for an exported field to be marshaled to JSON. Fields that do not have
// this tag are simply ignored.
//
// safejson works independently from the encoding/json package, allowing you to
// define different marshaling and unmarshaling rules for your data.
package safejson

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"
)

// Marshaler allows an object to provide its own serialized version of itself
type Marshaler interface {
	SafeMarshal() (interface{}, error)
}

var marshalerInterface = reflect.TypeOf((*Marshaler)(nil)).Elem()

func filter(v reflect.Value) (interface{}, error) {
	k := v.Kind()

	if k == reflect.Invalid {
		return nil, nil
	}

	log.Printf("%#v", v.Type().Kind() == reflect.Struct)

	for k == reflect.Ptr || k == reflect.Interface {
		if v.IsNil() {
			return nil, nil
		}

		if v.Type().Implements(marshalerInterface) {
			return v.Interface().(Marshaler).SafeMarshal()
		}

		v = v.Elem()
		k = v.Kind()
	}

	switch k {
	case reflect.Struct, reflect.Interface:
		result := map[string]interface{}{}

		t := v.Type()
		length := t.NumField()
		var err error

		for index := 0; index < length; index++ {
			f := t.Field(index)

			if f.PkgPath == "" {
				tag := f.Tag.Get("safejson")

				if tag != "" {
					val := v.Field(index)

					if tt, ok := val.Interface().(time.Time); ok {
						result[tag] = tt
					} else if result[tag], err = filter(v.Field(index)); err != nil {
						return nil, err
					}
				}
			}
		}

		return result, nil

	case reflect.Slice, reflect.Array:
		length := v.Len()
		result := make([]interface{}, length)
		var err error

		for index := 0; index < length; index++ {
			if result[index], err = filter(v.Index(index)); err != nil {
				return nil, err
			}
		}

		return result, nil
	}

	return v.Interface(), nil
}

// Filter marshals a value but does not convert it to []byte
func Filter(i interface{}) (interface{}, error) {
	v := reflect.ValueOf(i)

	vv := v

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Invalid:
	default:
		return nil, errors.New("Values passed to Marshal must be instances of struct, slice, or array, or conform to the Marshaler interface")
	}

	return filter(vv)
}

// Marshal converts the provided value to its “safe” JSON representation.
//
// The root value passed to the function must be a struct, slice, or array, or conform to the
// Marshaler interface, or the marshaling operation fails and an error is returned.
func Marshal(i interface{}) ([]byte, error) {
	if result, err := Filter(i); err == nil {
		return json.Marshal(result)
	} else {
		return nil, err
	}
}
