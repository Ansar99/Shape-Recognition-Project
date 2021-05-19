package main

// RUN by typing: go run main.go images/image4.jpg shapedImages/shaped_image4.jpg
// or choose your own image: go run main.go selected_image.jpg resulting_image_save.jpg

import (
	"fmt"
	"log"
	"os"
	"time"

	"gocv.io/x/gocv"
	"shapeitup.com/helper"
)

func main() {
	if len(os.Args) == 1 {
		runCamera()
	} else if len(os.Args) < 3 {
		fmt.Println("usage: go run main.go selected_image.jpg resulting_image_save.jpg")
		return
	}
	teststr := run(os.Args[1], os.Args[2])

	fmt.Println(teststr)
}

func run(imagepath string, outputimagepath string) string {
	start := time.Now()

	imageloc := imagepath
	saveResultAs := outputimagepath

	shapeimg, err := helper.ImageToGrayscaleMat(imageloc)
	if err != nil {
		log.Fatalf("Error while creating grayscaled Mat: %v", err)
	}

	updatedshapeimg, teststr := helper.MarkAndFindShapes(shapeimg)
	shapeimgconversion := gocv.IMWrite(saveResultAs, updatedshapeimg)
	if !shapeimgconversion {
		log.Fatalf("Error in converting" + saveResultAs)
	}

	elapsed := time.Since(start)
	log.Printf("shapeidentifier took: %s", elapsed)

	return teststr
}

func runCamera() {
	webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Shape Detect")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		shapeimg, err := helper.BlurMat(img)
		if err != nil {
			break
		}
		updatedshapeimg, teststr := helper.MarkAndFindShapes(shapeimg)
		window.IMShow(updatedshapeimg)

		if window.WaitKey(1) >= 0 {
			break
		}
		fmt.Println(teststr)
	}
}
