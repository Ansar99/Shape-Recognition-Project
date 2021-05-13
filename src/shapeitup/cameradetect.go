package main

import (
	"gocv.io/x/gocv"
	"shapeitup.com/helper"
)

func main() {
	webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Shape Detect")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		shapeimg, err := helper.BlurMat(img)
		if err != nil {
			break
		}
		updatedshapeimg := helper.MarkAndFindShapes(shapeimg)
		window.IMShow(updatedshapeimg)
		//window.WaitKey(1)

		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
