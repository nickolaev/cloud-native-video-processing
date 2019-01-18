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
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func realHandler(ff *FfmpegStreamer, w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

	var headers []string

	log.Println("Request received with headers:", r.Header)

	for header, values := range r.Header {
		for _, v := range values {
			if ff.IsSourceHTTP() && (strings.ToLower(header) == "process-video") {
				headers = append(headers, "-headers", header+": "+v)
			}
			// if config.is_text_overlay && header == "User" {
			// 	config.flags.text = "Registered to user " + v //TODO: unsafe, sanitize the string
			// 	generate_args()                               //force args regeneration
			// }
		}
	}

	argsWithHeaders := append(headers, ff.GetArgs()...)

	log.Println("Starting ffmpeg with args: ", argsWithHeaders)
	cmd := exec.CommandContext(r.Context(), "ffmpeg", argsWithHeaders...)
	log.Println("cmd.Args ", cmd.Args)

	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		io.Copy(w, outPipe)
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

func main() {
	ffmpeg := NewFfmpegStreamer()

	for {
		if ffmpeg.GetListen() {
			ffmpeg.runFfmpeg()
		} else {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				realHandler(ffmpeg, w, r)
			})
			log.Printf("Listening on %v", ffmpeg.GetListenHost())
			err := http.ListenAndServe(ffmpeg.GetListenHost(), nil)
			if err != nil {
				log.Panicf("ListenAndServe returned: %v", err)
			}
		}
	}
}
