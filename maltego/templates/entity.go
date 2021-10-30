package main

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

	"github.com/maxlandon/gondor/maltego"
)

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

// Target - Enter a description in place of this sentence (leaving the
// "{{Target}} -" name, so as to automatically include this
// comment as a description to your Target in the Maltego Client.
type Target struct {
	OS       string `display:"Operating System" name:"os" strict:"yes" alias:"alias"`
	Hostname string `name:"hostname" strict:"yes" alias:"host"`
	IP       string `display:"IP Address" name:"ip" alias:"address" overlay:"W,image"`
}

// AsEntity - This function makes the Target type a valid Maltego entity.
// It works by the very fact that the base maltego.Entity implements the
// valid MaltegoEntity itself, so that by returning this entity with your
// data stored in it, you creating your core Go type as a valid Entity.
func (t *Target) AsEntity() (e maltego.ValidEntity) {
	e = maltego.NewEntity(
		"{{.Template.Namespace}}", // Should default to runtime package path + type Name
		"Target Environment",      // Default to Name with space and Caps
		"target",                  // default to name
		t,                         // Can be nil, but there is not point doing it
	)
	// Any additional settings for the Maltego Entity
	// can be done here, by accessing the e variable.
	return e
}

// NewTarget - Instantiate a new {{Entity}} type. This type is compliant with
// the library's entity.Entity interface, and can be used in a Go transform.
func main() {

	// Register the entity to a Transform
	// transform := transform.NewTransform(
	//         "Transform Display name",
	//         TranformTest,
	//         target,
	// )

	// Create a new entity, with a native Go type as underlying data
	target := maltego.NewEntity(
		"{{.Template.Namespace}}",
		"Target Environment",
		"target",
		&Target{},
	)

	// Setting additional details
	target.Link.Direction = maltego.Bidirectional

	// Set the underlying Go type on the fly: when
	// the transform crafts the response, it will
	// use this type to produce additional properties.
	target.SetData(&Target{})

	// Add a field to the entity, regardless of its underlying type.
	// Field are always added AFTER the fields contained by the Entity
	// underlying Go type, if there is one.
	target.AddField(maltego.Field{
		Name:         "activeDirectoryDomain",
		Display:      "Active Directory Domain",
		Alias:        "addomain",
		Value:        "C://Users",
		MatchingRule: maltego.MatchLoose,
	})

	fmt.Println(target)
}
