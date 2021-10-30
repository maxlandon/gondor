package maltego

/*
   Gondor - Go Maltego Transform Framework
   Copyright (C) 2021 Maxime Landon

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"fmt"
	"reflect"
	"sync"
)

// ValidEntity - A rather special interface that is paradoxically not intended
// for use by the writer of an Entity. This is only to make sure that the core
// Go type declared by the user has a function in which he returns a base Entity
// type, so that we can use this core Go type as an input in a transform, along
// with the base Entity type holding all Maltego-related stuff and logic.
type ValidEntity interface {
	AsEntity() Entity
}

// Entity - A Go representation of a Maltego Entity type.
// Because the Maltego client might pass Entities inputs that are not Go native types,
// (or Go types not known to this program), this Entity type contains all properties and
// details for ANY Entity, no matter its origin, and users have access to many methods
// to add and query its properties, as well as to set its various Maltego details.
type Entity struct {
	// Base properties
	Namespace string // The Maltego namespace of this entity (Maltego entities always fit within a tree)
	Alias     string // The alias under which the Entity can be searched for/ grabbed.
	eType     string // The string representation of the Entity type (determined through reflection)
	Category  string // The category of entities to which this category belongs (eg: a DNS server => services)
	Value     string // The value of the Entity, used by the Maltego client
	Weight    int    // The weight attributed to this entity on the graph

	// Advanced properties
	// These properties are objects of their own because they bear
	// lots of metadata and Maltego-specific information with them.
	properties []Field // Additional fields, are the actual Entity properties, as a list to preserve order.
	Notes      string
	Bookmark   BookmarkColor
	labels     []Label

	// Display properties
	// These properties are all the other properties related to
	// how the entity is displayed, with its various overlays.
	IconURL  string                      // An optional URL to the Entity main icon
	overlays map[OverlayPosition]Overlay // Access the various overlays by their position.
	Link     Link                        // All settings of the link to/from this entity

	// Operating
	mutex *sync.RWMutex // Concurrency management
	data  interface{}   // Underlying native Go struct, holds base fields with struct tags, might be nil
}

// AsEntity - Self implementation of the Entity interface type.
// This function is very important for a few reasons:
// 1) You must return this function from within your own custom
//    implementation of this function with your Go native type.
// 2) When you write a transform that accepts an non-Go native
//    type as an Input Entity, the transform will automatically
//    process this Input Entity into a base type, before handing
//    it to you for query and usage within your transform func.
func (e Entity) AsEntity() Entity {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e
}

// SetData - Set the actual Go type that is to be treated as a Maltego entity in your transforms.
// The data passed MUST BE A POINTER TO A STRUCT. There are several situations where you want to
// use this function, but the most common one is:
// Somewhere in a Go package a native Go type has been declared with its struct tags, and somewhere
// (else or not), there is a call to create an entity with this type (used to set up its name, category,
// namespace, create its configuration template, etc). Because this entity will be probably be a singleton,
// you want to be able to replace its underlying data each time before being marshalled to Maltego XML.
func (e *Entity) SetData(data interface{}) (err error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	e.data = data
	return
}

func (e *Entity) GetData() interface{} {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.data
}

// Field - Returns the string value of a Property field (regardless of its true,
// underlying type), given the name (key) of the field as argument. If not found,
// the function returns an empty string.
func (e *Entity) Field(name string) string {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for _, p := range e.properties {
		if p.Name == name {
			return fmt.Sprintf("%v", p.Value)
		}
	}
	return ""
}

// AddField - Add a field to an Entity base type. You can use this
// function when you want to add a property to it because the input/output
// entity is either not a native Go type, or because you don't have access
// to it, or because you want this field to be added into the Maltego UI but
// not as a persistent struct field in your Go code.
//
// Note that you can't directly set a field as an overlay when declaring it
// through this function. You need to reference it again in Entity.AddOverlay().
func (e *Entity) AddField(f Field) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	e.properties = append(e.properties, f)
}

// AddOverlay - Set one of the Entity's overlay items, specifying three things:
// - Its value, which is MOST OF THE TIME the name of one of the Entity's fields,
// - Its position, which is a Go enum so that you can't pass an invalid one.
// - Its type, also as a Go enum to avoid invalid ones.
//
// Note that you can also specify entity fields as overlays when tagging a native
// Go type fields with the appropriate tags (overlay:"W,text", overlay:"N,image", etc).
// Please refer to the NewEntity() function documentation for info on these tags.
func (e *Entity) AddOverlay(value string, pos OverlayPosition, oType OverlayType) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	overlay := Overlay{
		Value:    value,
		Position: pos,
		Type:     oType,
	}
	e.overlays[pos] = overlay
}

// NewEntity - Instantiate a new Entity type. The interface data passed as parameter
// MUST BE A POINTER TO A STRUCT having some fields tagged with the following list.
// If this function is successfull, a Transform function will be able to use the
// interface type as a valid input.
// The returned Entity type allows you to set all and every Maltego-related settings
// for the Entity, but not its struct fields: this is because you should not ever need
// them after this function. Transforms only care about Go native types.
//
// ----- Type Compliance Documentation ------
//
// The following is an exhaustive list of all valid and/or required struct
// field tags for a type to be considered a valid & working Maltego Entity.
// We take the example of a field named IP, of type string:
//
// display:"IP Address"   - The display name of the field in Maltego (default: IP).
// name:"IP"              - The programmatic, Java/MaltegoXML name of the field
//                          (default is set through reflect).
// strict:"yes"           - If non nil, the Matching Rule of this field is "strict",
//                          otherwise it's "loose".
//                          ("loose"/"strict", default:"loose")
// alias:"ipaddress"      - The Maltego alias for this field.
// overlay:"W,image"      - Use the field as an overlay: notation is <Position>,<type>.
//                          Valid positions: W, N, S, C, NW, SW
//                          Valid types: text, image, colour/color
//                          If color is used, must be a valid RGB format (eg. #45e06f)
//
// ------------------------------------------
func NewEntity(namespace, category, alias string, data interface{}) Entity {
	e := Entity{
		mutex: &sync.RWMutex{},
	}

	// Compute default description from the core Go type

	return e
}

// TODO: check link sync.Mutex not nil when instantiating

// New - Instantiate a working Entity type to be embedded
func NewEntityDefault(data interface{}) Entity {
	e := Entity{
		mutex: &sync.RWMutex{},
	}

	// Get the namespace + Name from the Go runtime package + type
	// Set the Display name to the type name with spaces and caps
	// Set the alias

	// Compute default description from the core Go type

	return e
}

// Unmarshal - A Maltego entity is being passed a Go native type
// in which to unmarshal its properties. This function is useful
// when you want to cast an input entity into your native input
// type, while retaining the possibility of using the Entity.
func (e *Entity) Unmarshal(eType ValidEntity) (err error) {

	// Here, apply custom XML unmarshaling logic, from e to eType

	// If embedded structs, for each marked field in them, do:
	// structName.FieldName.
	// If struct is anonymous, ignore structName.

	return
}

func newEntityToXML(e interface{}, fieldDecoratorName string) {
	osField := reflect.Indirect(reflect.ValueOf(e)).FieldByName(fieldDecoratorName)
	if osField.Kind() == reflect.String {

	}

	numFields := reflect.TypeOf(e).Elem().NumField()

	// For each struct field, process its tags
	for i := 0; i < numFields; i++ {
		//
		field := reflect.TypeOf(e).Elem().Field(i)
		display, ok := field.Tag.Lookup("display")
		if !ok {
		} else {
			fmt.Println(display)
		}
	}
	// field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldDecoratorName)
	// display, ok := field.Tag.Lookup("display")
}
