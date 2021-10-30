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

// TransformFunc - This type defines what is the valid implementation of a Transform
// in Go code. The transform passed as argument is a "self-reference", which gives you
// access to all the methods for querying, modifying and adding input/output Entities,
// as well as some of the core Transform settings.
// This is another way to register a new valid Maltego transform, without
// wrapping it around a native Go type implementing maltego.ValidTransform.
type TransformFunc func(t Transform) (err error)

// Transform - The base Go implementation of a Maltego transform.
// This type holds all the information necessary to the correct registration
// and functioning of an equivalent Maltego Client Transform.
type Transform struct {
	Description string             // Defaults to the Go-doc comment of the user-provided TransformFunc
	Settings    []TransformSetting // All settings for this Transform (global/local, and all their properties)
	Input       Entity             // To be replaced by the Entity interface (we also need the Go native type in there)
	run         TransformFunc      // The transform function implementation, declared and passed by the user
}

// NewTransform - Instantiate a new Transform by passing a valid Transform function
// implementation. This leaves you the choice on where you want to declare this function
// whether it is a type method or a pure function (depends on your needs and code), etc.
func NewTransform(name string, run TransformFunc) Transform {
	t := Transform{
		Description: getTransformDescription(run),
		run:         run,
	}
	return t
}

func (t *Transform) AddEntity(e ValidEntity) (err error) {
	return
}

// validateInput - The transform checks that all Entity fields that need
// to satisfy some requirements/presence actually do that, and other checks.
func (t *Transform) validateInput() (err error) {
	return
}

// sendResponse - Once the core transform implementation has ran, package
// the resulting objects in a Maltego-compliant format, verify required fields
// and settings are correctly set, and send the result back to the Server.
func (t *Transform) sendResponse() (err error) {
	return
}

// TransformSetting - An individual Transform Setting, which can be customized
// by a user in control of a Transform type (through its .Settings field).
type TransformSetting struct {
	Name     string
	Display  string
	Type     string
	Default  string
	Optional bool
	Popup    bool
	Global   bool
}
