package screenshot

import (
	// #cgo LDFLAGS: -framework CoreGraphics
	// #cgo LDFLAGS: -framework CoreFoundation

	// #include <CoreGraphics/CoreGraphics.h>
	// #include <CoreFoundation/CoreFoundation.h>
	// #cgo CFLAGS: -x objective-c
	// #cgo LDFLAGS: -framework Foundation
	// #include <Foundation/Foundation.h>
	// #cgo LDFLAGS: -framework Cocoa
	// #include <Cocoa/Cocoa.h>
	/*
		void hello() {
				NSLog(@"Hello World");
		}
		float getScaleFactor(){
			//NSString * text = @"片仮名, カタカナ ABCD💣💣";
			//int stringLength = [text length];
			//NSLog(@"%d",stringLength);
			CGFloat scale = [[NSScreen mainScreen] backingScaleFactor];
			//NSLog(@"%f",scale);
			return scale;
		}
	*/
	"C"
	"image"
	"reflect"
	"unsafe"
)

//https://coderwall.com/p/l9jr5a/accessing-cocoa-objective-c-from-go-with-cgo
// func PrintTest() {
// 	//flt := C.printScaleFactor()
// 	//println(flt)
// }
var scalingFactor = C.getScaleFactor()

func ScreenRect() (image.Rectangle, error) {

	//a := C.NSString("qwe")
	displayID := C.CGMainDisplayID()
	width := int(C.float(C.CGDisplayPixelsWide(displayID)) * scalingFactor)
	//width := int(math.Ceil(float64(C.CGDisplayPixelsWide(displayID))/16) * 16)
	height := int(C.float(C.CGDisplayPixelsHigh(displayID)) * scalingFactor)
	return image.Rect(0, 0, width, height), nil
}

func CaptureScreen() (*image.RGBA, error) {
	rect, err := ScreenRect()
	if err != nil {
		return nil, err
	}
	return CaptureRect(rect)
}

func CaptureRect(rect image.Rectangle) (*image.RGBA, error) {
	displayID := C.CGMainDisplayID()
	//scalingFactor := C.getScaleFactor()
	width := int(C.CGDisplayPixelsWide(displayID))

	imagecap := C.CGDisplayCreateImage(displayID)
	provider := C.CGImageGetDataProvider(imagecap)
	rawData := C.CGDataProviderCopyData(provider)

	//src_bytes_per_row := C.CGImageGetBytesPerRow(imagecap)
	//src_bytes_per_pixel := C.CGImageGetBitsPerPixel(imagecap) / 8

	length := int(C.CFDataGetLength(rawData))
	ptr := unsafe.Pointer(C.CFDataGetBytePtr(rawData))

	var slice []byte
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdrp.Data = uintptr(ptr)
	hdrp.Len = length
	hdrp.Cap = length

	imageBytes := make([]byte, length)

	for i := 0; i < length; i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
	}

	C.CFRelease(C.CFTypeRef(rawData))

	img := &image.RGBA{Pix: imageBytes, Stride: int(C.float(4*width) * scalingFactor), Rect: rect}
	return img, nil
}
