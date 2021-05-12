package main

// RUN by typing: go run main.go images/image4.jpg shapedImages/shaped_image4.jpg
// or choose your own image: go run main.go selected_image.jpg resulting_image_save.jpg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"shapeitup.com/helper"
	"gocv.io/x/gocv"
)


func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: go run main.go selected_image.jpg resulting_image_save.jpg")
		return
	}
	start := time.Now()

	imageloc := os.Args[1]
	saveResultAs := os.Args[2]

	shapeimg, err := imageToGrayscaleMat(imageloc)
	if err != nil {
		log.Fatalf("Error while creating grayscaled Mat: %v", err)
	} else {
		//fmt.Println("Creation of grayscaled Mat succeeded")
	}

	updatedshapeimg := helper.MarkAndFindShapes(shapeimg)
	shapeimgconversion := gocv.IMWrite(saveResultAs, updatedshapeimg)
	if shapeimgconversion {
		//fmt.Println(saveResultAs + " saved successfully")
	} else {
		log.Fatalf("Error in converting" + saveResultAs)
	}
	elapsed := time.Since(start)
	log.Printf("shapeidentifier took: %s", elapsed)
}


// imageToGrayScaleMat converts a path of an image to a gocv.Mat, grayscales it and blurs it.
// if the path is invalid an error occours.
// returns a gocv.Mat or an error.
func imageToGrayscaleMat(imgname string) (gocv.Mat, error) {
	img := gocv.IMRead(imgname, gocv.IMReadGrayScale)
	if img.Empty() {
		return gocv.Mat{}, errors.New("image img in imageToGrayscaleMat is empty")
	}
	gocv.MedianBlur(img, &img, 11)

	shapeimg := gocv.NewMat()
	gocv.CvtColor(img, &shapeimg, gocv.ColorGrayToBGR)
	if shapeimg.Empty() {
		return gocv.Mat{}, errors.New("image shapeimg in imageToGrayscaleMat is empty")
	}

	return shapeimg, nil
}
