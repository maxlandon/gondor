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

// Field - A property field for a Maltego entity. You can use this
// type from within a transform, when you want to add a property to
// it because the input/output entity is either not a native Go type,
// or because you don't have access to it, or because you want this
// field to be added into the Maltego UI but not in the persistent Go code.
//
// Note that you can't directly set a field as an overlay when declaring it
// through this function. You need to reference it again in Entity.AddOverlay().
type Field struct {
	Name         string       // The programmatic name, required.
	Display      string       // The display name of this field
	Alias        string       // An alias for the field, default to .Name
	Value        interface{}  // Its value, automatically passed as an XML string
	MatchingRule MatchingRule // The individual match rule for this field
}

// String - The most basic type of Entity field that is available in Gondor,
// also used by more advanced field types provided by this package.
// type String struct {
//         name        string
//         isValue     bool
//         displayName string
//         alias       string
//         matchRule   MatchingRule // Default is strict, ensured when marshalling XML
//         err         error
//         decorator   func() string
// }
