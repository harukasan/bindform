// Package bindform implements mapping URL query string values or post form values
// into Go struct values.
package bindform

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

// BindForm binds form values into the each fields of the v. BindoForm also use
// URL query string as form value.
func BindForm(r *http.Request, v interface{}) error {
	r.ParseForm()

	structDef := reflect.TypeOf(v).Elem()
	structField := reflect.ValueOf(v).Elem()
	for i := 0; i < structDef.NumField(); i++ {
		field := structField.Field(i)
		tag := structDef.Field(i).Tag.Get("form")
		err := bindValue(field, r.FormValue(tag))
		if err != nil {
			return err
		}
	}

	return nil
}

// BindPostForm binds post form values into the each fields of the v.
func BindPostForm(r *http.Request, v interface{}) error {
	r.ParseForm()

	structDef := reflect.TypeOf(v).Elem()
	structField := reflect.ValueOf(v).Elem()
	for i := 0; i < structDef.NumField(); i++ {
		field := structField.Field(i)
		tag := structDef.Field(i).Tag.Get("form")
		err := bindValue(field, r.PostFormValue(tag))
		if err != nil {
			return err
		}
	}

	return nil
}

func bindValue(field reflect.Value, v string) error {
	switch field.Kind() {
	case reflect.Bool:
		return bindBool(field, v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return bindInt(field, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return bindUInt(field, v)
	case reflect.Float32, reflect.Float64:
		return bindFloat(field, v)
	case reflect.String:
		field.SetString(v)
		return nil
	}
	return errors.New("can't bind value")
}

func bindBool(field reflect.Value, str string) error {
	if str == "" {
		return nil
	}

	v, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}
	field.SetBool(v)
	return nil
}

func bindInt(field reflect.Value, str string) error {
	if str == "" {
		return nil
	}

	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	field.SetInt(v)
	return nil
}

func bindUInt(field reflect.Value, str string) error {
	if str == "" {
		return nil
	}

	v, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}
	field.SetUint(v)
	return nil
}

func bindFloat(field reflect.Value, str string) error {
	if str == "" {
		return nil
	}

	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	field.SetFloat(v)
	return nil
}
