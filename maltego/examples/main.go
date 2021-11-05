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
	"log"

	"github.com/maxlandon/gondor/maltego"
)

func main() {

	//
	// 1) Distributions & Servers ---------------------------------------------
	//

	// Declare a new Maltego Distribution, which will hold
	// all our Entities, Transforms, Machines and associated content.
	// This distribution can hold information for several servers.
	dist := maltego.NewDistribution()

	// Alternatively, create a single Maltego Transform Server. However:
	// This server contains its Distribution, but its content will be
	// limited to the Transforms you have registered in this server, and
	// their associated input/output Entities.
	// (This is why you cannot pass an existing Distribution to a server,
	// because it would not contain the Transforms/Entities implementations)
	server := maltego.NewTransformServer(nil)

	// However, we can register this Server and its Entities/Transforms
	// into the distribution file, which will automatically merge itself with
	// the Server's one. You're still in charge of starting the Server, obviously.
	dist.RegisterServer(server)

	//
	// 2) Entities ------------------------------------------------------------
	//

	// You can declare a Go type (struct, whatever), and make
	// it implement the maltego.ValidTransform interface. This is
	// handy for several reasons:
	// - A type can be only a valid MaltegoEntity and not a Transform
	// - A type can be a valid transform but not an Entity
	// - A type can be both, if it implements the two interfaces.
	//
	// Here, the credential can be both an Entity and a Transform,
	// but you have to instantiate it first.
	cred := &Credential{}

	// Marshalling an native type back to an Entity, and modifying its fields.
	credential := cred.AsEntity()
	credential.AddLabel("My Display Title", "This Content")

	//
	// 3) Transforms ----------------------------------------------------------
	//

	// There are several ways of declaring and
	// registering valid Go-native Maltego transforms:

	// 1) You declare a maltego.TransformFunc somewhere,
	// and create a transform directly out of it.
	transform := maltego.NewTransform(
		"Transform Display name",
		ProducerTransform,
	)

	// You can also declare a method around a native Go type
	// (which might also happen to be an Entity, that is valid),
	// and register this method as the Transform implementation.
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

	//
	// 4) Registering Transforms & Entities -----------------------------------
	//

	// Again, there are several ways in which you can register your
	// Entities and Transforms for being used by Maltego Clients.
	// These ways partly depend on what choices you've made in the
	// Section 1 of this example:
	// - One Distribution holding potentially more than one server.
	// - One server able to produce a Distribution only for its own content.

	// A - Register to the server -
	// All transforms are automatically bound to a URL
	// path matching their complete namespace + Name.
	// Their Entities are also registered in the Server's distribution.
	server.RegisterTransform(&transformOnly)
	server.RegisterTransform(&transform)

	// B - Register to a distribution -
	// This has a drawback, which is that the transform will be mapped
	// to a default (local) Server contained in the distribution.
	// Thus, you should not have to use this function: prefer declaring
	// a Server, map the transform to it, and map the server to the dist.
	dist.RegisterTransform(credentialTransform)

	// Additionally, when you have implemented an Entity that is not yet
	// used as a Transform Input but that you wish users to access in Maltego,
	// you can register it in a Distribution (a server's one or not).
	server.RegisterEntity(cred)

	// Our server is registered in the Distribution, but the latter
	// takes care of reconciling itself with the server's one, so can
	// do this (double assignment) without any risk.
	dist.RegisterEntity(cred)

	//
	// 5) Creating Distributions & Serving Transforms -------------------------
	//

	// A - Producing / Loading Distribution -
	// Before serving any set of transforms, we need to produce the
	// Distribution file that will be loaded into the Maltego Client
	// (therefore adding Transforms, Entities, Machines, Servers, etc)
	// The Distribution can write itself to an arbitrary path, under
	// the name DistributionName.mtz.
	err := dist.WriteToFile("path/to/directory")
	if err != nil {
		log.Fatal(err)
	}

	// You can do the same for the Server configuration if you want,
	// but in this case it's useless because the Server's in the dist.
	err = server.WriteToFile("path/to/server/config")
	if err != nil {
		log.Fatal(err)
	}

	// B - Starting Transform Servers
	// Start serving the transforms, supposing -here- that we loaded
	// a complete Transform & Registry configuration, ports, TLS, etc.
	server.ListenAndServe()
}
