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

func main() {

	// Transforms --------------------------------
	// There are several ways of declaring and
	// registering valid Go-native Maltego transforms:

	// 1) You declare a maltego.TransformFunc somewhere,
	// and create a transform directly out of it.
	transform := maltego.NewTransform(
		"Transform Display name",
		ProducerTransform,
	)

	// 2) You can declare a Go type (struct, whatever), and make
	// it implement the maltego.ValidTransform interface. This is
	// handy for several reasons:
	// - A type can be only a valid MaltegoEntity and not a Transform
	// - A type can be a valid transform but not an Entity
	// - A type can be both, if it implements the two interfaces.
	//
	// Here, the credential can be both an Entity and a Transform,
	// but you have to instantiate it first.
	cred := &Credential{}
	credentialTransform := maltego.NewTransform(
		"Transform Display name",
		cred.Do,
	)

	// Here, the MyTransform type can only be a transform,
	// and could not be used as an Entity. This, however,
	// doesn't change anything to the Transform workflow.
	myTransform := UpdaterTransform{}
	transformOnly := maltego.NewTransform(
		"Transform Display name",
		myTransform.Do,
	)

	// Marshalling an Entity to XML
	credential := cred.AsEntity()
	credential.TranslateProperties()
	credential.AddDisplayInformation("My Display Title", "This Content")

	// We will serve all our transforms on one server instance
	server := maltego.NewTransformServer(nil)

	// All transforms are automatically bound to a URL path
	// matching their complete namespace + Name, and this will
	// be equivalent to any Registry Configurations that you
	// might load into a Transform Distribution Server.
	server.RegisterTransform(&credentialTransform)
	server.RegisterTransform(&transformOnly)

	// Start serving the transforms, supposing -here- that we loaded
	// a complete Transform & Registry configuration, ports, TLS, etc.
	server.ListenAndServe()
}
