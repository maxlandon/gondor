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

// MyForeignTransform - Declaring the implementation of a transform that accepts
// a foreign, non-Go-native type as an Entity input. Thus, when using this function
// you can only access the input Entity field through its dedicated methods.
var MyForeignTransform = func(t maltego.Transform) (err error) {

	// However, you have direct access to the input Entity,
	// without the need for further unmarshaling calls:

	return
}

// UpdaterTransform - A Go type that we intend to be a valid Maltego Transform.
type UpdaterTransform struct {
	// You can add any internal logic in here, which accessorily
	// means that you can write a Transform implementation around
	// one of your existing Go types.
}

// Do - This function is a valid maltego.TransformFunc.
// The MyTransform type is now a valid Maltego transform.
// Note that you can declare any number of differently named
// methods around your type: as long as their signature is the
// following, you can register them all as different Transforms.
func (t UpdaterTransform) Do(mt *maltego.Transform) (err error) {

	// You still have access to the transform input Entity:
	mt.Request.Entity.AddProperty(maltego.Field{Display: "New Field"})

	// Add and process any arbitrary Go types in this body.
	// However, you will only be able to return as output Entities
	// those satisfying the maltego.ValidEntity interface.
	// Please refer to the pure function example below.

	// We have added a field to the input entity, it's obviously
	// because our transform is (in part ?) an "updating" transform.
	mt.AddEntity(mt.Request.Entity)

	return
}

// ProducerTransform - Declaring the implementation of a tranform accepting
// a native Go type as an Entity input, with compile-time validity check.
var ProducerTransform = func(t *maltego.Transform) (err error) {

	var target = &Target{}                   // If your type is implements maltego.ValidEntity...
	err = t.Request.Entity.Unmarshal(target) // ...You can make this call, checked at compile-time

	// You can create a new version of your Entity, with all its default settings
	// that you have declared in your constructor, and modify them on the fly,
	// applying only once, for this transform.
	//
	// WARNING: THIS DOES NOT PASS THE INPUT ENTITY STATE
	// (because the target is a NEW valid Entity instance)
	out := target.AsEntity()
	out.Weight = 200
	out.Link.Reverse()
	out.AddOverlay("myOverlayName", maltego.OverlayCenter, maltego.OverlayImage)

	// And finally return this on-the-fly modified entity
	t.AddEntity(out)

	return
}
