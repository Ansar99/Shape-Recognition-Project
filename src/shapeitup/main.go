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
	if len(os.Args) < 3 {
		fmt.Println("usage: go run main.go selected_image.jpg resulting_image_save.jpg")
		return
	}
	start := time.Now()

	imageloc := os.Args[1]
	saveResultAs := os.Args[2]

	shapeimg, err := helper.ImageToGrayscaleMat(imageloc)
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
