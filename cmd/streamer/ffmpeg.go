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

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// FfmpegStreamer holds the configuration of ffmpeg
type FfmpegStreamer struct {
	// the parsed flags
	flags struct {
		source    string
		listen    bool
		port      int
		width     int
		height    int
		logo      string
		font      string
		timestamp bool
		text      string
		textfile  string
	}
	isSourceHTTP  bool
	isTextOverlay bool
	// the ffmpeg generated args
	args []string
}

func (ff *FfmpegStreamer) processFlags() {

	source := flag.String("source", "", "The video source")
	listen := flag.Bool("listen", false, "Make ffmpeg serve the http connections")
	port := flag.Int("port", 10100, "The listen port")
	width := flag.Int("width", -1, "Output video width")
	height := flag.Int("height", -1, "Output video height")
	preset := flag.String("preset", "",
		"WxH preset. Possible values are:\n 144p, 240p, 360p, 480p, 720p, HD, WXGA+, HD+, FHD")
	logo := flag.String("logo", "", "Overlay the logo file in the top left corner of the video")
	font := flag.String("font", "", "Specify the overlay font to use")
	timestamp := flag.Bool("timestamp", false, "Enable inserting a localtime timestamp")
	text := flag.String("text", "", "Specify text to show")
	textfile := flag.String("textfile", "", "Specify the textfile to show")

	flag.Parse()

	// Flag validation
	if len(strings.TrimSpace(*source)) == 0 {
		log.Fatalln("Please supply a video source")
	}

	ff.isSourceHTTP = strings.HasPrefix(strings.TrimSpace(*source), "http")

	if p, ok := scalePresets[*preset]; ok {
		*width = p.width
		*height = p.height
	}

	ff.isTextOverlay = len(*text) > 0 || len(*textfile) > 0 || *timestamp

	if len(strings.TrimSpace(*font)) == 0 && ff.isTextOverlay {
		log.Fatalln("Please specify the font to use")
	}

	ff.flags.source = *source
	ff.flags.listen = *listen
	ff.flags.port = *port
	ff.flags.width = *width
	ff.flags.height = *height
	ff.flags.logo = *logo
	ff.flags.font = *font
	ff.flags.timestamp = *timestamp
	ff.flags.text = *text
	ff.flags.textfile = *textfile
}

func (ff *FfmpegStreamer) generateArgs() {
	f := ff.flags

	process := (f.width > 0) || (f.height > 0)
	process = process || (len(f.logo) > 0)
	process = process || ff.isSourceHTTP
	process = process || f.timestamp || (len(f.text) > 0) || (len(f.textfile) > 0)

	args := []string{}

	if process {
	} else {
		args = append(args, "-fflags", "+genpts", "-stream_loop", "-1")
	}

	args = append(args, "-i", f.source)

	if process {
		if (f.width > 0) || (f.height > 0) {
			// scale
			args = append(args, "-vf",
				"scale="+strconv.Itoa(f.width)+":"+strconv.Itoa(f.height))
		} else if len(f.logo) > 0 {
			// logo overlay
			args = append(args, "-i", f.logo,
				"-filter_complex", "[0:v][1:v] overlay=20:0")
		} else if ff.isTextOverlay {
			// text overlay
			text := ":text='%{localtime\\:%T}'" // default to timestamp
			if len(f.text) > 0 {
				text = ":text=" + f.text
			} else if len(f.textfile) > 0 {
				text = ":textfile=" + f.textfile
			}

			args = append(args, "-vf",
				"drawtext=fontsize=50:fontfile="+f.font+":fontcolor=White"+
					text+
					":box=1:boxcolor=Black"+
					":x=w-tw-50:y=h-th-30")
		}

		// do not recode audio
		args = append(args, "-c:a", "copy")
		// use x264 at fastest speed
		args = append(args, "-c:v", "libx264",
			"-preset", "ultrafast", "-tune", "zerolatency",
			"-crf", "15")
	} else {
		args = append(args, "-c", "copy")
	}

	args = append(args, "-f", "mpegts")
	if f.listen {
		args = append(args,
			"-listen", "1",
			"http://0.0.0.0:"+strconv.Itoa(f.port))
	} else {
		// output the result to stdout
		args = append(args, "-")
	}

	ff.args = args
}

func (ff *FfmpegStreamer) runFfmpeg() {
	log.Println("Starting ffmpeg with args: ", ff.args)
	cmd := exec.Command("ffmpeg", ff.args...)

	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		io.Copy(os.Stdout, outPipe)
	}()

	go func() {
		io.Copy(os.Stderr, errPipe)
	}()

	err = cmd.Wait()
	if err != nil {
		log.Println("Exited with :", err)
	}
	log.Println("-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-")
}

// GetListen returns the listen config flag
func (ff *FfmpegStreamer) GetListen() bool {
	return ff.flags.listen
}

// GetListenHost returns the ffmpeg listen host
func (ff *FfmpegStreamer) GetListenHost() string {
	return ":" + strconv.Itoa(ff.flags.port)
}

// IsSourceHTTP checks if the source is HTTP
func (ff *FfmpegStreamer) IsSourceHTTP() bool {
	return ff.isSourceHTTP
}

// GetArgs returns args
func (ff *FfmpegStreamer) GetArgs() []string {
	return ff.args
}

// NewFfmpegStreamer creates the new ffmpeg state object
func NewFfmpegStreamer() *FfmpegStreamer {
	new := &FfmpegStreamer{}
	new.processFlags()
	new.generateArgs()

	return new
}
