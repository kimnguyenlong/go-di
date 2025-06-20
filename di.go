// Simple dependency injection framework for Go applications.
package di

import (
	"errors"
	"fmt"
	"reflect"
)

// Container is a map from object names to their instances.
type Container map[string]any

// NewContainer return an empty container.
func NewContainer() *Container {
	return &Container{}
}

// Plug registers an object in the container with a given name.
// The name is used to identify the object when injecting dependencies.
// The value can be any type, but it's concrete type must be a pointer to a struct for dependency injection to work.
func (c *Container) Plug(name string, value any) {
	if c == nil {
		return
	}

	(*c)[name] = value
}

// Wire performs dependency injection by resolving all registered objects in the container.
// It iterates over each object, checking for fields tagged with "di" and injecting the corresponding dependencies.
func (c *Container) Wire() error {
	if c == nil {
		return errors.New("container is nil")
	}
	for name, obj := range *c {
		if err := c.inject(obj); err != nil {
			return fmt.Errorf("di.Wire() failed for %s: %w", name, err)
		}
	}
	return nil
}

// inject injects dependencies into the fields of the given object.
// It checks for fields tagged with "di" and assigns the corresponding dependencies from the container.
// If the field is not settable, it returns an error.
// If the dependency type does not match the field type, it returns an error.
// If the injection is successful, it sets the field with the dependency value.
// If the field is a pointer to a struct, it recursively do the injection for that field.
func (c *Container) inject(obj any) error {
	if c == nil {
		return errors.New("container is nil")
	}

	if obj == nil {
		return errors.New("object is nil")
	}

	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("object is not a pointer to a struct")
	}

	v := reflect.ValueOf(obj).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Pointer && v.Field(i).Elem().Kind() == reflect.Struct {
			if err := c.inject(v.Field(i).Interface()); err != nil { // Recursively inject if it's a pointer to a struct
				return err
			}
			continue
		}

		field := v.Type().Field(i)
		tag := field.Tag.Get("di")
		if len(tag) == 0 {
			continue // No injection tag
		}

		if !v.Field(i).CanSet() { // Check if the field is settable
			return fmt.Errorf("cannot set field '%s'", field.Name)
		}

		dependency, exists := (*c)[tag]
		if !exists { // Check if the dependency exists in the container
			return fmt.Errorf("dependency '%s' not found for field '%s'", tag, field.Name)
		}

		depValue := reflect.ValueOf(dependency)
		if !depValue.Type().AssignableTo(v.Field(i).Type()) { // Check if the dependency type matches the field type
			return fmt.Errorf("cannot assign dependency '%s' to field '%s' (type mismatch)", tag, field.Name)
		}

		// Set the field with the dependency
		v.Field(i).Set(reflect.ValueOf(dependency))
	}

	return nil
}
