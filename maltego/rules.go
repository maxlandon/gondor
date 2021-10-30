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

// MatchingRule - Matching rules are used to specify how an entity will be
// merged in the Maltego user interface. Strict matching specifies that an
// entity will only be merged with another if all it's fields (including the
// value) are equal. Loose matching specifies that two entities will be merged
// if only the entity values are equal.
type MatchingRule string

const (
	// MatchStrict - Strict matching specifies that an entity will only be
	// merged with another if all it's fields (including the value) are equal.
	MatchStrict MatchingRule = "strict"

	// MatchLoose - Loose matching specifies that two entities
	// will be merged if only the entity values are equal.
	MatchLoose MatchingRule = "loose"
)

// Label - Used to convey extra information associated with an Entity in the Maltego
// client GUI. Unlike entity fields, labels are only transmitted in response messages
// and cannot be passed from transform to transform as a source of input.
type Label struct {
	Value string // add xml tags or replace with xml type
	Type  string // add xml tags or replace with xml type (default: text/text)
	name  string // add xml tags or replace with xml type
}
