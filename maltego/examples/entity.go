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
func (tgt *Target) AsEntity() (e maltego.Entity) {

	// This call results in an Entity type whose basic/needed operational
	// settings have been set with, in part, information on the Go type
	// (here Credential), like the Go module package path, names, etc.
	//
	// For instance, the Maltego namespace of the Credential entity is,
	// by default, the complete Go-module path+name of the Credential type.
	// Please see the Credential type below for an example where we modify it.
	return maltego.NewEntity(tgt)
}

// Credential - A native Go type that has some struct fields declared as properties,
// and some others ignored by default. You will be able, however to add them on the
// fly from within a Transform implementation, with entity.AddField().
type Credential struct {
	Login      string `display:"Login"`                   // Declaring a new Maltego field
	PublicKey  []byte `display:"Public Key" strict:"yes"` // This PublicKey is unique among all graph entities
	PrivateKey []byte // For whatever reason you might want NOT to push the PrivateKey to Maltego
}

// AsEntity - This type is a valid Maltego entity.
func (cred *Credential) AsEntity() (e maltego.Entity) {

	// This call results in an Entity type whose basic/needed operational
	// settings have been set with, in part, information on the Go type
	// (here Credential), like the Go module package path, names, etc.
	//
	// For instance, the Maltego namespace of the Credential entity is,
	// by default, the complete Go-module path+name of the Credential type.
	e = maltego.NewEntity(cred)

	// You can still modify the settings if you want
	e.Link.Reverse()         // This link will be an output to input one.
	e.Link.Color = "#43eb36" // Must be a valid RGB color code.

	// Add dynamic fields that you don't have in your native
	// Go type fields, for whatever reasons. Know however that
	// this is "quite" pointless, as this method acts as constructor
	// each time your type is used an Entity. Still, you can do it.
	//
	// If this case, we might like to read from a file containing
	// the Public key, which you might not want to do from within
	// the transform each time.
	e.AddProperty(maltego.Field{
		Name:         "PublicKey", // Will override the c.PublicKey field.
		Display:      "Public Key",
		Value:        "your Public Key bytes here",
		MatchingRule: maltego.MatchStrict,
	})

	e.AddOverlay("PublicKey", maltego.OverlayNorth, maltego.OverlayText)
	e.IconURL = "http://domain.dom/image.png"

	return e
}

// Do - The credential type, by implementing the Do() func,
// also satisfies the maltego.ValidTransform interface.
func (cred *Credential) Do(mt *maltego.Transform) (err error) {

	// You can make this call, checked at compile-time
	err = mt.Request.Entity.Unmarshal(cred)

	// Completely overwrite the input Entity settings
	mt.AddEntity(cred)

	return
}
