package di

import (
	"fmt"
	"reflect"
)

var graph = map[string]*Object{}

type Object struct {
	Name  string
	Value interface{}
}

func Plug(objects ...*Object) {
	for _, obj := range objects {
		graph[obj.Name] = obj
	}
}

func Wire() error {
	for _, obj := range graph {
		if err := inject(obj); err != nil {
			return err
		}
	}
	return nil
}

func inject(obj *Object) error {
	if obj.Value == nil {
		return nil // Nothing to inject
	}

	// Use reflection to find fields with `di` tags and inject dependencies
	val := reflect.ValueOf(obj.Value).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("di")
		if tag == "" {
			continue // No injection tag
		}

		dependency, exists := graph[tag]
		if !exists {
			return fmt.Errorf("dependency '%s' not found for field '%s'", tag, field.Name)
		}

		val.Field(i).Set(reflect.ValueOf(dependency.Value))
	}

	return nil
}
