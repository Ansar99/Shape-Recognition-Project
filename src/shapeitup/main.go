package main

// RUN BY TYPING: go run main.go ./images/circles.jpg
// or choose your own image: go run main.go ./images/selected_image.jpg

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: go run main.go ./images/selected_image.jpg selected_image.jpg") //ANTON
		return
	}

	imageloc := os.Args[1]
	imagename := os.Args[2]

	img, circleimg, shapeimg, err := imageToGrayscaleMat(imageloc)
	defer img.Close()
	defer circleimg.Close()
	defer shapeimg.Close()
	if err != nil {
		log.Fatalf("Error while creating grayscaled Mat: %v", err)
	} else {
		fmt.Println("Creation of grayscaled Mat succeeded")
	}

	updatedshapeimg := markAndFindShaped(shapeimg, img)
	shapeimgconversion := gocv.IMWrite("./shapedImages/shaped_"+imagename+".jpg", updatedshapeimg) //ANTON
	if shapeimgconversion {
		fmt.Println("Image2 saved successfully")
	} else {
		log.Fatalf("Error in converting image2") // error handling.
	}
}

func markAndFindShaped(shapeimg gocv.Mat, img gocv.Mat) gocv.Mat {
	canny := gocv.NewMat()
	defer canny.Close()
	gocv.Canny(shapeimg, &canny, 10, 10)

	contours := gocv.FindContours(canny, gocv.RetrievalList, gocv.ChainApproxTC89L1)
	imgpoints := contours.ToPoints() // type: [][]image.Point

	for i := 0; i < contours.Size(); i = i + 2 {
		shapeimgpointvector := imgpoints[i]
		shapevector := gocv.NewPointVectorFromPoints(shapeimgpointvector)

		/*shapeperimetertest := gocv.ArcLength(shapevector, true)
		shapeguesstest := gocv.ApproxPolyDP(shapevector, 0.03*shapeperimetertest, true)
		shapeguessRightTypetest := shapeguesstest.ToPoints()
		fmt.Println(shapeguessRightTypetest)
		greentest := color.RGBA{0, 255, 0, 0}
		if len(shapeguessRightTypetest) == 8 {
			for j := 0; j < 8; j++ {
				gocv.Circle(&shapeimg, shapeguessRightTypetest[j], 2, greentest, 3)
			}
		}*/ // Bra att ha för att se punkterna på en cirkel, om nu isOctagon behöver omdefinieras.

		shape, textpoint := detectshape(shapevector)
		red := color.RGBA{255, 0, 0, 0}
		gocv.PutText(&shapeimg, shape, textpoint, 2, 1, red, 1)
		fmt.Printf("\t %s\n", shape)
	}

	gocv.DrawContours(&shapeimg, contours, -1, color.RGBA{0, 0, 255, 0}, 1)

	return shapeimg
}

func detectshape(pvr gocv.PointVector) (string, image.Point) {
	shape := "unidentified shape"
	shapeperimeter := gocv.ArcLength(pvr, true)                     // Calculates the perimeter (omkrets) of the current shape
	shapeguess := gocv.ApproxPolyDP(pvr, 0.03*shapeperimeter, true) // A polygonal curve - A curve that is entirely made up of line segments. ( no arcs) A closed curve - A curve that begins and ends in the same location.
	// 0.03 är ren trial and error, den verkar fungera som den ska.

	shapeguessRightType := shapeguess.ToPoints() // Done for valuable performance gain

	textpoint := shapeguessRightType[0]
	textpoint.X = textpoint.X - 5
	vertices := len(shapeguessRightType)
	if vertices == 3 {
		shape = "triangle"
	} else if vertices == 4 {
		// Points (X,Y) of the rectangle/square
		p1 := shapeguessRightType[0]
		p2 := shapeguessRightType[1]
		p3 := shapeguessRightType[2]
		p4 := shapeguessRightType[3]

		// If difference in distances is < 10, then its a square, otherwise rectangle
		p1p2 := calculateDistanceBetweenTwoPoints(p1, p2)
		p2p3 := calculateDistanceBetweenTwoPoints(p2, p3)
		p3p4 := calculateDistanceBetweenTwoPoints(p3, p4)
		p4p1 := calculateDistanceBetweenTwoPoints(p4, p1)
		if (p1p2-p2p3) < 10 && (p2p3-p3p4) < 10 && (p4p1-p1p2) < 10 { // Creating a perfect square is hard, therefore a margin of 10, contour also somewhat rounds corners so it's needed either way.
			shape = "square"
		} else {
			shape = "rectangle"
		}
	} else if vertices == 5 {
		shape = "pentagon"
	} else if vertices == 6 {
		shape = "hexagon"
	} else if vertices == 7 {
		shape = "heptagon"
	} else if vertices == 8 {
		if isOctagon(shapeguessRightType) {
			shape = "octagon"
		} else {
			shape = "circle"
		}
	} else if vertices == 9 {
		shape = "nonagon"
	} else {
		shape = "circle"
	} // Lägg till mer former här! :-D

	return shape, textpoint
}

func calculateDistanceBetweenTwoPoints(point1 image.Point, point2 image.Point) float64 {
	return math.Sqrt(float64((point2.X-point1.X)*(point2.X-point1.X)) + float64((point2.Y-point1.Y)*(point2.Y-point1.Y)))
}

func isOctagon(points []image.Point) bool {
	var pointsArr [8]int
	for i := 0; i < 8; i++ {
		for _, v := range pointsArr {
			if v == points[i].X {
				return true
			}
		}
		pointsArr[i] = points[i].X
	}
	/*
		lägg till i array
		nästa värde - kolla om array innehåller värdet, om det gör det return true, annars fortsätt
		inget värde samma som något annat, returna false.
	*/
	return false
}

func imageToGrayscaleMat(imgname string) (gocv.Mat, gocv.Mat, gocv.Mat, error) {
	img := gocv.IMRead(imgname, gocv.IMReadGrayScale) // Laddar in bilden satt av Args[1] och lägger den i variabeln img som är en grayscalade Mat.
	if img.Empty() {
		return gocv.Mat{}, gocv.Mat{}, gocv.Mat{}, errors.New("Image img in imageToGrayscaleMat is empty")
	}
	gocv.MedianBlur(img, &img, 11) // blurrar bilden, tar bort noise

	shapeimg := gocv.NewMat() // skapar en ny mat
	circleimg := gocv.NewMat()
	gocv.CvtColor(img, &shapeimg, gocv.ColorGrayToBGR) // Kopierar över bilden från img till shapeimg samt circleimg
	gocv.CvtColor(img, &circleimg, gocv.ColorGrayToBGR)
	if shapeimg.Empty() {
		return gocv.Mat{}, gocv.Mat{}, gocv.Mat{}, errors.New("Image shapeimg in imageToGrayscaleMat is empty")
	}

	return img, circleimg, shapeimg, nil
}
