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

// Overlay - An overlay is a piece of information that is displayed
// at some position relative (close) to the Entity. An overlay can
// be a piece of text, an image or a color. You specify its position
// and its value which, for example if the Overlay type is an Image
// can hold the URL to this image, or in other cases, the name of
// the Entity property that must be used to find a value.
type Overlay struct {
	Value    string          // Either an Entity property name, a URL to an image, etc
	Position OverlayPosition // The relative position of the item relative to the Entity
	Type     OverlayType     // The type of overlay that we want to show.
}

// OverlayPosition - The position of a Maltego Entity Overlay element.
type OverlayPosition string

const (
	OverlayNorth     OverlayPosition = "N"
	OverlaySouth     OverlayPosition = "S"
	OverlayWest      OverlayPosition = "W"
	OverlayNorthWest OverlayPosition = "NW"
	OverlaySouthWest OverlayPosition = "SW"
	OverlayCenter    OverlayPosition = "C"
)

// OverlayType - The type of a Maltego Entity Overlay element.
type OverlayType string

const (
	OverlayImage  OverlayType = "image"
	OverlayColour OverlayType = "color"
	OverlayText   OverlayType = "text"
)
