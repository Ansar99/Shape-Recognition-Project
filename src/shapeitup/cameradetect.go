package main

import (
	"fmt"
	"image"
	"shapeitup.com/helper"
	"gocv.io/x/gocv"
)

// A Result represents a shape and a point where text should be written.
type Result struct {
	Shape     string      // The shape
	Textpoint image.Point // The point where text is placed
}

func main() {
	webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Shape Detect")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		shapeimg, err := imageToGrayscaleMat(img)
		if err != nil {
			fmt.Println("hej")
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


func imageToGrayscaleMat(imgname gocv.Mat) (gocv.Mat, error) {
	gocv.MedianBlur(imgname, &imgname, 11)
	shapeimg := gocv.NewMat()
	gocv.BilateralFilter(imgname, &shapeimg, 10, float64(100), float64(100))
	//gocv.CvtColor(imgname, &shapeimg, gocv.ColorGrayToBGR)

	return shapeimg, nil
}
