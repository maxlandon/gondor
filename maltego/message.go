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

import "encoding/xml"

// Message - A type containing all the output elements of a Transform.
type Message struct {
	x         xml.Name
	Response  TransformResponseMessage  `xml:"MaltegoTransformResponseMessage,omitempty"`
	Exception TransformExceptionMessage `xml:"MaltegoTransformExceptionMessage,omitempty"`
}

// TransformResponseMessage - A type containing all the output elements of a Transform.
type TransformResponseMessage struct {
	Entities []Entity    `xml:"Entities"`   // All entities to be returned as the Transform output.
	Messages []MessageUI `xml:"UIMessages"` // Transform log messages
}

// TransformExceptionMessage - A type containing all the exceptions (errors) that
// occured during the execution of a Transform. While you can return an error at
// any point in your code, thereby terminating execution, you can also simply log
// them with Transform.AddError(), and they will be passed along any other output.
type TransformExceptionMessage struct {
	Exceptions []Exception
}

// Exception - Term for an error in a Transform. Can be terminating, or not.
type Exception string

// MessageUI - A log message passed along a Transform
// output for display in the Maltego transform window.
type MessageUI struct {
	Text string `xml:"text"`
	Type string `xml:"type"`
}
