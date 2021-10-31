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

//
// Maltego Transforms - Specification & Instantiation ------------------------------------------
//

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
)

// TransformFunc - This type defines what is the valid implementation of a Transform
// in Go code. The transform passed as argument is a "self-reference", which gives you
// access to all the methods for querying, modifying and adding input/output Entities,
// as well as some of the core Transform settings.
// This is another way to register a new valid Maltego transform, without
// wrapping it around a native Go type implementing maltego.ValidTransform.
//
// Any error returned from the function will be translated into a Maltego Transform exception.
// You can return an error at any time within your Tranform function implementation.
type TransformFunc func(t *Transform) (err error)

// Transform - The base Go implementation of a Maltego transform.
// This type holds all the information necessary to the correct registration
// and functioning of an equivalent Maltego Client Transform, and exactly as
// in the Python code, is NOT restricted to any type of output Entity.
type Transform struct {
	Description string             // Defaults to the Go-doc comment of the user-provided TransformFunc
	Settings    []TransformSetting // All settings for this Transform (global/local, and all their properties)
	Request     Request            // Other information and settings passed by the Maltego Client
	Input       Entity             // To be replaced by the Entity interface (we also need the Go native type in there)
	run         TransformFunc      // The transform function implementation, declared and passed by the user
	entities    []Entity           // All entities to be returned as the Transform output.
	messages    []MessageUI        // Transform log messages
	exceptions  []Exception        // All errors throwed during execution.
	mutex       *sync.RWMutex      // Concurrency
}

// NewTransform - Instantiate a new Transform by passing a valid Transform function
// implementation. This leaves you the choice on where you want to declare this function
// whether it is a type method or a pure function (depends on your needs and code), etc.
//
// The function allows you to pass an optional list of transform settings that you want to
// apply to this transform AND ONLY THIS ONE. If you want global settings (applying to all
// transforms served by an HTTP server), pass this settings to the server itself.
// You can also add settings to the Transform later, through its AddSetting() method. In all
// cases, you should always register them BEFORE serving the Transforms to their client.
func NewTransform(name string, run TransformFunc, settings ...TransformSetting) Transform {
	t := Transform{
		Description: getTransformDescription(run),
		run:         run,
		Settings:    settings,
		mutex:       &sync.RWMutex{},
	}
	return t
}

//
// Maltego Transforms - User API -------------------------------------------------------------
//

// AddEntity - Add an Entity to the list of entities to be sent in the Transform response.
// Generally, you want to call it with either yourGoType.AsEntity() function, or directly
// passing a maltego.Entity type when you can't/don't want to use a native Go type in the Transform.
func (t *Transform) AddEntity(e ValidEntity) (err error) {
	// Do not append the entity if the we topped
	// the maximum allowed number of output entities.
	if t.Request.Slider == len(t.entities) {
		return
	}
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	t.entities = append(t.entities)
	return
}

// Debugf - Log an debug-level message in the Maltego transform window.
func (t *Transform) Debugf(format string, args ...interface{}) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	msg := fmt.Sprintf(format, args)
	t.messages = append(t.messages, MessageUI{Text: msg, Type: "Debug"})
}

// Infof - Log an info-level message in the Maltego transform window.
func (t *Transform) Infof(format string, args ...interface{}) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	msg := fmt.Sprintf(format, args)
	t.messages = append(t.messages, MessageUI{Text: msg, Type: "Inform"})
}

// Warnf - Log an warning-level message in the Maltego transform window.
func (t *Transform) Warnf(format string, args ...interface{}) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	msg := fmt.Sprintf(format, args)
	t.messages = append(t.messages, MessageUI{Text: msg, Type: "Warning"})
}

// Errorf - Log an error-level message in the Maltego transform window.
// This function returns the error, so that if you want to terminate the
// transform because of it, you can "return err" from anywhere.
func (t *Transform) Errorf(format string, args ...interface{}) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	msg := fmt.Sprintf(format, args)
	t.exceptions = append(t.exceptions, Exception(msg))
	return errors.New(msg)
}

//
// Transform Internal Implementation -----------------------------------------------
//

// unmarshalXML - Unmarshal a Transform Request passed as XML (sent over HTTP).
func (t *Transform) unmarshalXML() (err error) {
	return
}

// unmarshalArgs - Unmarshal a local Transform Request, invoked via command-line args.
func (t *Transform) unmarshalArgs() (err error) {
	return
}

// validateInput - The transform checks that all Entity fields that need
// to satisfy some requirements/presence actually do that, and other checks.
func (t *Transform) validateInput() (err error) {
	return
}

// marshalOutput - The transform packages the output Entities within an XML string.
func (t *Transform) marshalOutput(runErr error) (out []byte, err error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Message container
	message := Message{
		x: xml.Name{Local: "MaltegoMessage"},
	}

	// We have either failed (and the error is already stored)
	if runErr != nil {
		message.Exception = TransformExceptionMessage{
			Exceptions: t.exceptions,
		}
	}

	// Or succeeded, with output entities and UI messages
	if runErr == nil {
		message.Response = TransformResponseMessage{
			Entities: t.entities,
			Messages: t.messages,
		}
	}

	// Marshal the overall message and its content.
	return xml.Marshal(message)
}

// sendResponse - Once the core transform implementation has ran, package
// the resulting objects in a Maltego-compliant format, verify required fields
// and settings are correctly set, and send the result back to the Server.
func (t *Transform) sendResponse() (err error) {
	return
}
