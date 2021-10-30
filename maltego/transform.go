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
	"github.com/maxlandon/gondor/maltego/entity"
)

type TransformFunc func(t Transform)

// TransformAlt - Any type implementing the Transform interface is automatically
// a compliant Maltego transform, able to process arbitrary input Entities and
// return arbitrary output Entities. The interface also allows the user to set
// meta-level details on the Transform.
type TransformAlt interface {
	ValidateInput() error // The transform checks the input Entity compliance and casting to Go type
	Execute(eData interface{}) error
	AddEntity(entity.Entity) error
	SendResponse() error
}

// Transform - The base Go implementation of a Maltego transform.
// This type holds all the information necessary to the correct registration
// and functioning of an equivalent Maltego Client Transform.
type Transform struct {
	Description string             // Defaults to the Go-doc comment of the user-provided TransformFunc
	Settings    []TransformSetting // All settings for this Transform (global/local, and all their properties)
	Input       entity.Entity      // To be replaced by the Entity interface (we also need the Go native type in there)
}

// func NewTransform(name string, tr TransformAlt, input entity.Entity) (t Transform) {
func NewTransform(name string, run TransformFunc, input entity.Entity) (t Transform) {
	t = Transform{
		Description: getTransformDescription(run),
	}

	return
}

func (t *Transform) AddEntity(e entity.Entity) (err error) {
	return
}

// ValidateInput - The transform checks that all Entity fields that need
// to satisfy some requirements/presence actually do that, and other checks.
func (t *Transform) ValidateInput() (err error) {
	return
}

// SendResponse - Once the core transform implementation has ran, package
// the resulting objects in a Maltego-compliant format, verify required fields
// and settings are correctly set, and send the result back to the Server.
func (t *Transform) SendResponse() (err error) {
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

type globalConfig struct {
	Settings []TransformSetting
}

// GlobalConfigFromFile - Reads the Maltego Transform Configuration file located
// at path. If not found, returns a default, empty (but non-nil) configuration, and
// an error to indicate the user that some action might be required for perfect work.
func GlobalConfigFromFile(path string) (conf *globalConfig, err error) {
	return
}

// GlobalConfigFromBytes - Unmarshal a Maltego Transform Configuration as bytes
// If unmarshaling fails, returns a default, empty (but non-nil) configuration, and
// an error to indicate the user that some action might be required for perfect work.
func GlobalConfigFromBytes(data []byte) (conf *globalConfig, err error) {
	return
}
