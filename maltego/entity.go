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
	"runtime/debug"
	"strings"
	"sync"
)

//
// Maltego Entities - Specification & Instantiation ------------------------------------------
//

// ValidEntity - A rather special interface that is paradoxically not intended
// for use by the writer of an Entity. This is only to make sure that the core
// Go type declared by the user has a function in which he returns a base Entity
// type, so that we can use this core Go type as an input in a transform, along
// with the base Entity type holding all Maltego-related stuff and logic.
type ValidEntity interface {
	AsEntity() Entity // The type is able to wrap itself into a maltego.Entity
}

// Entity - A Go representation of a Maltego Entity type.
// Because the Maltego client might pass Entities inputs that are not Go native types,
// (or Go types not known to this program), this Entity type contains all properties and
// details for ANY Entity, no matter its origin, and users have access to many methods
// to add and query its properties, as well as to set its various Maltego details.
type Entity struct {
	// Base properties
	Namespace   string `xml:"-"` // The Maltego namespace of this entity (Maltego entities always fit within a tree)
	DisplayName string // Defaults to the camelCase-split Entity type if Go native.
	Alias       string `xml:"-"`         // The alias under which the Entity can be searched for/ grabbed.
	Type        string `xml:"Type,attr"` // The string representation of the Entity type (determined through reflection)
	Category    string `xml:"-"`         // The category of entities to which this category belongs (eg: a DNS server => services)
	Value       string `xml:",cdata"`    // The value of the Entity, used by the Maltego client
	Weight      int    `xml:"Weight"`    // The weight attributed to this entity on the graph

	// Display properties
	// These properties are all the other properties related to
	// how the entity is displayed, with its various overlays.
	Link     Link          `xml:"-"`                        // Wraps itself into Properties later.
	IconURL  string        `xml:"IconURL,omitempty"`        // An optional URL to the Entity main icon
	Bookmark BookmarkColor `xml:"-"`                        // Wraps itself into Properties later.
	Overlays Overlays      `xml:"Overlays"`                 // Access the various overlays by their position.
	Labels   []Label       `xml:"DisplayInformation>Label"` // Additional display information

	// The actual Entity properties, as a list to preserve order.
	// When this Entity is an Input to a Transform, the underlying
	// Go type is not automatically populated: you have to call the
	// Entity.Unmarshal(&YourType{}) method to get it populated.
	// You'll still be able to access the corresponding fields, but
	// the returned value type will always be a string.
	Properties Properties `xml:"AdditionalFields"`

	// Operating
	mutex *sync.RWMutex `xml:"-"` // Concurrency management
	data  interface{}   `xml:"-"` // Underlying native Go struct, holds base fields with struct tags, might be nil
}

// NewEntity - Instantiate a new Entity type. The interface data passed as parameter
// MUST BE A POINTER TO A STRUCT having some fields tagged with the following list.
// If this function is successfull, a Transform function will be able to use the
// interface type as a valid input.
// The returned Entity type allows you to set all and every Maltego-related settings
// for the Entity, but not its struct fields: this is because you should not ever need
// them after this function. Transforms only care about Go native types.
//
// Struct Tags & Type Compliance:
//
// The following is an exhaustive list of all valid and/or required struct
// field tags for a type to be considered a valid & working Maltego Entity.
// We take the example of a field named IP, of type string:
//
// display:"IP Address"   - Required. The display name of the field in Maltego
// strict:"yes"           - If non nil, the Matching Rule of this field is "strict",
//                          otherwise it's "loose".
//                          ("loose"/"strict", default:"loose")
// alias:"ipaddress"      - The Maltego alias for this field.
// overlay:"W,image"      - Use the field as an overlay: notation is <Position>,<type>.
//                          Valid positions: W, N, S, C, NW, SW
//                          Valid types: text, image, colour/color
//                          If color is used, must be a valid RGB format (eg. #45e06f)
// hidden:"yes"           - If not nil, the field is hidden in the Properties Window.
// sample:"127.0.0.1"     - A value used when the Entity is created manually in Maltego.
// default:"0.0.0.0"      - A value that is always populated by default.
//
func NewEntity(data interface{}) Entity {
	e := Entity{
		Overlays:   Overlays{},
		Properties: Properties{},
		mutex:      &sync.RWMutex{},
		data:       data,
	}

	// Get the namespace + Name from the Go runtime package + type
	e.Namespace = reflect.TypeOf(data).Elem().PkgPath()
	e.Type = reflect.TypeOf(data).Elem().Name()

	// Set the Display name to the type name with spaces and caps
	e.DisplayName = e.Type
	bi, _ := debug.ReadBuildInfo()
	fmt.Println(bi.Main.Path)
	// name := "DNSToIp"
	// re := regexp.MustCompile(`([0-9]+)`)
	// name = re.ReplaceAllString(name, "$1")
	// fmt.Println(name)
	// re2 := regexp.MustCompile(`([a-z])([A-Z]+)`)
	// name = re2.ReplaceAllString(name, "$1$1")
	// fmt.Println(name)
	// re3 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	// name = re3.ReplaceAllString(name, "$1$1")
	// fmt.Println(name)
	// re := regexp.MustCompile(`[A-Z][a-z]+|[A-Z]+(![a-z])`)
	// e.DisplayName = strings.Join(re.FindAllString(e.Type, -1), " ")

	return e
}

//
// Maltego Entities - User API -------------------------------------------------------------
//

// AsEntity - Self implementation of the Entity interface type.
// This function is very important for a few reasons:
// 1) You always implicitly return this function from within your own
//    custom implementation of this function with your Go native type.
// 2) When you write a transform that accepts an non-Go native
//    type as an Input Entity, the transform will automatically
//    process this Input Entity into a base type, before handing
//    it to you for query and usage within your transform func.
func (e Entity) AsEntity() Entity {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e
}

// Property - Returns the string value of a Property field (regardless of its true,
// underlying type), given the name (key) of the field as argument. If not found,
// the function returns an empty string.
func (e *Entity) Property(name string) string {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for _, p := range e.Properties {
		if p.Name == name {
			return fmt.Sprintf("%v", p.Value)
		}
	}
	return ""
}

// AddProperty - Add a field to an Entity base type. You can use this
// function when you want to add a property to it because the input/output
// entity is either not a native Go type, or because you don't have access
// to it, or because you want this field to be added into the Maltego UI but
// not as a persistent struct field in your Go code.
//
// Note that you can't directly set a field as an overlay when declaring it
// through this function. You need to reference it again in Entity.AddOverlay().
func (e *Entity) AddProperty(p Field) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	e.Properties[p.Name] = p
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
		PropertyName: value,
		Position:     pos,
		Type:         oType,
	}
	e.Overlays[pos] = overlay
}

// AddLabel - Add a specific Display information to this Entity.
// If the title argument is nil (""), it will default to "Info".
func (e *Entity) AddLabel(title, content string) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	if title == "" {
		title = "Info"
	}
	e.Labels = append(e.Labels, Label{
		Name:    title,
		Content: content,
		Type:    "text/html", // Fixed for all labels
	})
}

// SetNote - Set the note for this Entity.
func (e *Entity) SetNote(note string) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	e.AddProperty(Field{
		Name:    "notes#",
		Display: "Notes",
		Value:   note,
	})
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

//
// Maltego Entity - Internal Implementation -----------------------------------------------
//

// getBaseProperties - Creates Entity properties from some builtin
// types of the Go Entity, like Links, Bookmarks, etc.
func (e *Entity) getBaseProperties() (err error) {

	// The link should add all its content to the list of properties
	e.AddProperty(Field{
		Name:    "link#maltego.link.color",
		Display: "LinkColor",
		Value:   e.Link.Color,
	})
	e.AddProperty(Field{
		Name:    "link#maltego.link.style",
		Display: "LinkStyle",
		Value:   e.Link.Style,
	})
	e.AddProperty(Field{
		Name:    "link#maltego.link.thickness",
		Display: "Thickness",
		Value:   e.Link.Thickness,
	})
	e.AddProperty(Field{
		Name:    "link#maltego.link.label",
		Display: "Label",
		Value:   e.Link.Label,
	})
	e.AddProperty(Field{
		Name:         "link#maltego.link.direction",
		Display:      "link#maltego.link.direction", // ??
		MatchingRule: MatchLoose,
		Value:        e.Link.Direction,
	})
	for _, property := range e.Link.properties {
		e.AddProperty(property)
	}

	// The bookmark as a property
	e.AddProperty(Field{
		Name:    "bookmark#",
		Display: "Bookmark",
		Value:   e.Bookmark,
	})

	return
}

// getGoProperties - This function uses reflection to package all valid fields in the struct
// (as interface) stored by the Entity, as properties. We overwrite them directly each time.
func (e *Entity) getGoProperties() (err error) {
	if e.data == nil {
		return
	}

	// For everything but structs, we directly package and return.
	switch e.data.(type) {
	case struct{}, *struct{}:
		// But we send structs in a recursive loop, for any embedded structs.
		ptrval := reflect.ValueOf(e.data)
		realval := reflect.Indirect(ptrval)
		e.marshalStruct("", realval, nil)
	default:
		// Simply add the field with fmt.Sprintf representation of the data.
		// This might be big, so people better know what they are passing.
		e.AddProperty(Field{
			Name:         "Go Type: " + reflect.TypeOf(e.data).String(),
			MatchingRule: MatchLoose,
			Value:        e.data,
		})
	}

	return
}

// marshalStruct - Marshal a struct with an arbitrary level of nesting, and package its content as Properties.
func (e *Entity) marshalStruct(namespace string, realval reflect.Value, sField *reflect.StructField) {

	// Always add a "root" separation property, with the type as Key and the name as value
	e.AddProperty(Field{
		Name:         getNamespace(namespace, realval.Type().Name()),
		Display:      realval.String(),
		MatchingRule: MatchLoose,
		Value:        "Go type",
	})

	// Compute the current namespace for this struct
	nested := getNamespace(namespace, sField.Name)

	// For each field, check tags
	numFields := realval.Type().NumField()
	for fieldCount := 0; fieldCount < numFields; fieldCount++ {
		field := realval.Type().Field(fieldCount)
		fieldKind := field.Type.Kind()
		valueKind := field.Type.Elem().Kind()

		fieldVal := realval.Field(fieldCount) // Can be nil

		// We can't read unexported fields, nor
		if !field.IsExported() {
			continue
		}

		// If the field is itself a struct, create a new
		// namespace level and call this func recursively.
		if field.Type.Kind() == reflect.Struct {
			e.marshalStruct(nested, fieldVal, &field)
			continue
		}

		// If the field is a pointer to a struct, we must
		// check that it's non nil here, cannot loop before.
		if fieldKind == reflect.Ptr && valueKind == reflect.Struct {

			// Initialize an empty value in the struct, then
			// go into a new recursive loop on this one.
			if fieldVal.IsNil() {
				fieldVal = reflect.New(fieldVal.Type().Elem())
			}
			e.marshalStruct(nested, reflect.Indirect(fieldVal), &field)
			continue
		}

		// The only required is display:"", not nil
		if _, ok := field.Tag.Lookup("display"); !ok {
			continue
		}

		// Process MatchRules and Aliases
		var match = MatchLoose
		value, ok := field.Tag.Lookup("match")
		if ok && value != "" {
			match = MatchStrict
		}
		aliasTag, ok := field.Tag.Lookup("alias")
		if !ok || aliasTag == "" {
			aliasTag = strings.ToLower(field.Name)
		}

		// Else, pick the tags and populate field
		f := Field{
			Name:         namespace,
			Value:        field.Name,
			Display:      field.Type.Name(),
			MatchingRule: match,
			Alias:        aliasTag,
		}
		e.AddProperty(f)

		// Finally, if this field is marked as an overlay, create it.
		overlayTag, yes := field.Tag.Lookup("overlay")
		if !yes {
			continue
		}
		e.addFieldAsOverlay(f, overlayTag)
	}

	return
}

// getNamespace - Compute the namespace for a field (or a series of them)
func getNamespace(namespace, name string) string {
	full := strings.Join([]string{namespace, strings.ToLower(name)}, ".")
	return strings.Trim(full, ".")
}

// addFieldAsOverlay - A struct field has been tagged as overlay,
// so validate it, create it and register it to the entity.
func (e *Entity) addFieldAsOverlay(f Field, tag string) {
	infos := strings.Split(tag, ",")
	if len(infos) == 1 && infos[0] == "" {
		return
	}

	// If we have only the position, we're fine,
	// we default the type as text.
	if len(infos) == 1 && isOverlayPosition(infos[0]) {
		e.AddOverlay(f.Name, OverlayPosition(infos[0]), OverlayText)
		return
	}

	// If we have two items, we will be fine in one case, not in the other
	if len(infos) == 2 {
		// If none is good, return
		if !isOverlayPosition(infos[0]) && !isOverlayType(infos[1]) {
			return
		}
		// Type is invalid, use text as default
		if isOverlayPosition(infos[0]) && !isOverlayType(infos[1]) {
			e.AddOverlay(f.Name, OverlayPosition(infos[0]), OverlayText)
			return
		}
		// Both types are valid, populate both
		if isOverlayPosition(infos[0]) && isOverlayType(infos[1]) {
			e.AddOverlay(f.Name, OverlayPosition(infos[0]), OverlayType(infos[1]))
		}
	}
}

// writeConfig - The Entity creates a file in path/Entities/EntityName,
// and writes itself as an XML message into it.
func (e Entity) writeConfig(path string) (err error) {
	return
}
