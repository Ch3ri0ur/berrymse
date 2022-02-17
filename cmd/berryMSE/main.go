// BerryMSE: Low-latency live video via Media Source Extensions (MSE)
// Copyright (C) 2020 Chris Hiszpanski
// https://github.com/thinkski/berrymse
// Modified Work Copyright (C) 2022 Christian Wappler and Maximilian Koch
// https://github.com/thinkski/berrymse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (

	//CMD Prints
	"fmt"


	//For static website stuff in binary executable github.com/markbates/pkger
	"github.com/markbates/pkger"

	// OWN STUFF ##########################################################################

	//FLAG & Configuration
	Config "github.com/ch3ri0ur/berrymse/src/config"

	//BerryMSE Main Pkg
	BerryMSE "github.com/ch3ri0ur/berrymse/src/BerryMSE"
)

//Init methode
//Defining Flags and Default values
func init() {
	Config.DefaultFlagInit()
}

func main() {

	fmt.Println("Starting BerryMSE Streaming")

	var configuration = Config.SetupConfigFlags()
	BerryMSE.BerryMSE(configuration, pkger.Dir("/cmd/berryMSE/web/static"))

}
