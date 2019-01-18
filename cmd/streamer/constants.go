// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

type resolution struct {
	width  int
	height int
}

var scalePresets = map[string]resolution{
	// Youtube 144p
	"144p": {
		width:  256,
		height: 144,
	},
	// Youtube 240p
	"240p": {
		width:  426,
		height: 240,
	},
	// Youtube 360p
	"360p": {
		width:  480,
		height: 360,
	},
	// Youtube 480p
	"480p": {
		width:  852,
		height: 480,
	},
	// 720p
	"720p": {
		width:  1280,
		height: 720,
	},
	// HD - 27.16% Web users 13.33% Steam users 08/2018
	"HD": {
		width:  1366,
		height: 768,
	},
	//  WXGA+ - 6.61% Web users 3.37% Steam users 08/2018
	"WXGA+": {
		width:  1440,
		height: 900,
	},
	// HD+ - 5.58% Web users 3.55% Steam users 08/2018
	"HD+": {
		width:  1600,
		height: 900,
	},
	// FHD - 19.57% Web users 63.72% Steam users 08/2018
	"FHD": {
		width:  1920,
		height: 1080,
	},
}
