# Descritpion

This project implements a Cloud Native Video Processing application - streamer. It's core is a thin shim written in Go, that calls ffmpeg binary for the real video manipulation.

## Quick start

### Local Streamer pipeline

To run a static Streamer pipeline locally:

```
./scripts/local.sh
```

This will expose a set of video processing functions on the following ports:
```
http://localhost:2000 - the video source
http://localhost:2001 - the video source + the logo
http://localhost:2002 - the video source + the logo + the timestamp
http://localhost:2003 - the video source + the logo + the timestamp + resized to 1366x768
```

### Kubernetes Streamer pipeline
If you are running a local Kubernetes cluster, you can simple build and images in the local Docker and deploy:

```
make build
make k8s-deploy
```

Please consult the Docker documentation if you need to export the images and upload them to your Kubernetes cluster image rository.

### Istio Streamer pipeline

Similarly to Kubernetes, a basic local installation is done by executing:

```
make build
make istio-deploy
```

More complex setups will require advanced knowledge, which is beyond the scope of this Quick Start.

## Configuration
Streamer inteprets its arguments and gnereates the relevant ffmpeg command line options. Currently it supports the following flags:

```
  -font string
    	Specify the overlay font to use
  -height int
    	Output video height (default -1)
  -listen
    	Make ffmpeg serve the http connections
  -logo string
    	Overlay the logo file in the top left corner of the video
  -port int
    	The listen port (default 10100)
  -preset string
    	WxH preset. Possible values are:
    	 144p, 240p, 360p, 480p, 720p, HD, WXGA+, HD+, FHD
  -source string
    	The video source
  -text string
    	Specify text to show
  -textfile string
    	Specify the textfile to show
  -timestamp
    	Enable inserting a localtime timestamp
  -width int
    	Output video width (default -1)
```

## Design

The generic video processing setup includes the Streamer listneing for incoming HTTP requests, spawning a chilf ffmpeg process to output in the stdout pipe. The content of that pipe is stramed back as and HTTP video stream, in reposnse to the initial HTTP request.

```
                 +--------------+
       HTTP GET  |              | spawn
     +----------->   Streamer   +--+
                 |              |  |
                 +------^-------+  |
                        |          |
                        |      +---v----+
                        | pipe |        |  HTTP GET
                        +------+ ffmpeg +------------>
                               |        |
                               +--------+

```

Depending on the command line arguments, the spawned ffmpeg can feed its processing pipeline from a remote HTTP connection e.g. another Streamer. Optionally it can read a local video stream, a file or other video source.

# Copyrights

This project is licnesed under the Apache 2.0 license. Please see the main LICENSE.txt for more details.

***video.mp4*** "*Rory and the snow*" by Stephen McPolin is
Licensed under Public Domain.

***bedstead.otf*** from https://fontlibrary.org/en/font/bedstead
Licensed under Public Domain.
