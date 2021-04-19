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
	"time"

	"gocv.io/x/gocv"
)

type Result struct {
	Shape     string
	Textpoint image.Point
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go ./images/selected_image.jpg")
		return
	}
	start := time.Now()

	imageloc := os.Args[1]

	img, shapeimg, err := imageToGrayscaleMat(imageloc)
	defer img.Close()
	defer shapeimg.Close()
	if err != nil {
		log.Fatalf("Error while creating grayscaled Mat: %v", err)
	} else {
		fmt.Println("Creation of grayscaled Mat succeeded")
	}

	updatedshapeimg := markAndFindShapes(shapeimg, img)
	shapeimgconversion := gocv.IMWrite("./shapedImages/shapedimage2.jpg", updatedshapeimg)
	if shapeimgconversion {
		fmt.Println("Image2 saved successfully")
	} else {
		log.Fatalf("Error in converting image2") // error handling.
	}
	elapsed := time.Since(start)
	log.Printf("shapeidentifier took: %s", elapsed)
}

func markAndFindShapes(shapeimg gocv.Mat, img gocv.Mat) gocv.Mat {
	canny := gocv.NewMat()
	defer canny.Close()
	gocv.Canny(shapeimg, &canny, 10, 10)

	contours := gocv.FindContours(canny, gocv.RetrievalList, gocv.ChainApproxTC89L1)
	imgpoints := contours.ToPoints() // type: [][]image.Point
	amtOfJobs := contours.Size()

	jobs := make(chan int, amtOfJobs)
	result := make(chan Result, amtOfJobs)

	for amountOfRoutines := 0; amountOfRoutines < amtOfJobs; amountOfRoutines++ {
		go worker(shapeimg, contours, imgpoints, jobs, result)
		//fmt.Println("created new worker")
	}
	//shapes := make(chan gocv.Mat, amtOfJobs)

	for i := 0; i < amtOfJobs; i = i + 2 {
		jobs <- i
	}
	close(jobs)

	for j := 0; j < amtOfJobs; j = j + 2 {
		shaperesult := <-result
		red := color.RGBA{255, 0, 0, 0}
		gocv.PutText(&shapeimg, shaperesult.Shape, shaperesult.Textpoint, 2, 1, red, 1)
		fmt.Printf("\t %s\n", shaperesult.Shape)
	}

	gocv.DrawContours(&shapeimg, contours, -1, color.RGBA{0, 0, 255, 0}, 1)
	return shapeimg
}

func worker(shapeimg gocv.Mat, contours gocv.PointsVector, imgpoints [][]image.Point, jobs <-chan int, result chan<- Result) {
	for i := range jobs {
		shapeimgpointvector := imgpoints[i]
		shapevector := gocv.NewPointVectorFromPoints(shapeimgpointvector)

		shapedetectresult := detectshape(shapevector)
		result <- shapedetectresult
	}
}

func detectshape(pvr gocv.PointVector) Result {
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

	var result Result
	result.Shape = shape
	result.Textpoint = textpoint

	return result
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
	return false
}

func imageToGrayscaleMat(imgname string) (gocv.Mat, gocv.Mat, error) {
	img := gocv.IMRead(imgname, gocv.IMReadGrayScale) // Laddar in bilden satt av Args[1] och lägger den i variabeln img som är en grayscalade Mat.
	if img.Empty() {
		return gocv.Mat{}, gocv.Mat{}, errors.New("Image img in imageToGrayscaleMat is empty")
	}
	gocv.MedianBlur(img, &img, 11) // blurrar bilden, tar bort noise

	shapeimg := gocv.NewMat()                          // skapar en ny mat
	gocv.CvtColor(img, &shapeimg, gocv.ColorGrayToBGR) // Kopierar över bilden från img till shapeimg samt circleimg
	if shapeimg.Empty() {
		return gocv.Mat{}, gocv.Mat{}, errors.New("Image shapeimg in imageToGrayscaleMat is empty")
	}

	return img, shapeimg, nil
}
