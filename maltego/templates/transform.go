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
	"fmt"

	"github.com/maxlandon/gondor/maltego"
)

// Transform - Enter a description in place of this sentence (leaving the
// "{{Transform}} -" name, so as to automatically include this
// comment as a description to your Transform in the Maltego Client.
type Transform struct {
	*maltego.Transform
	target *Target
}

// Execute - The function implementing your transform functionality.
func (t *Transform) Execute() (err error) {
	t.target = &Target{}
	if err = t.Input.Unmarshal(t.target); err != nil {
		return
	}

	return
}

func (tr *Transform) TranformTestAlt(input maltego.Entity) (output maltego.Entity, err error) {
	// target := input.Data.(*Target)

	return
}

// func TestTransform(t transform.Transform)

var testTransform = func(t maltego.Transform) (err error) {

	var nativeTarget = &Target{}
	err = t.Input.Unmarshal(nativeTarget)
	if err != nil {
		return
	}

	fmt.Println(nativeTarget.IP)

	t.AddEntity(nativeTarget.AsEntity())

	return
}
