# ScreenRec
Simple cross-platform pure Go screen recorder. (tested on linux&windows&osx)
* screen capture routines now work on retina displayes as well - added scaling factor support.

This project is actually an experiment to see how go integrates with native code & with piping data to ffmpeg via stdin.
That said, it is fully functional 
it produces about 4.5 fps on my macPro, which is not much - but this is not a real screen recorder, which typically have OS drivers.

<br/>

## Installation
```go
go get github.com/amitbet/screenrec
```

<br/>

## Basic Usage
just run ./screenrec
you can also mess around with the main.go file, and try on different encodings (all but mjpeg use ffmpeg):
* MJPeg avi
* x264 mp4
* dv9 webm
* dv8 webm

<br/>

## Dependencies
* **All** - FFMpeg file in the same directory or path
* **Windows** - None
* **Linux/FreeBSD** - https://github.com/BurntSushi/xgb
* **OSX** - cgo (CoreGraphics,CoreFoundation,Cocoa that should not be a problem)

<br/>

## Examples
Look at `example/` folder.
