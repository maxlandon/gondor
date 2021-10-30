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
