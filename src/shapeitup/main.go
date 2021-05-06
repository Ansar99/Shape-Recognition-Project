package main

// RUN by typing: go run main.go images/image4.jpg shapedImages/shaped_image4.jpg
// or choose your own image: go run main.go selected_image.jpg resulting_image_save.jpg

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"gocv.io/x/gocv"
)

// A Result represents a shape and a point where text should be written.
type Result struct {
	Shape     string      // The shape
	Textpoint image.Point // The point where text is placed
}

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
		fmt.Println("Creation of grayscaled Mat succeeded")
	}

	updatedshapeimg := markAndFindShapes(shapeimg)
	shapeimgconversion := gocv.IMWrite(saveResultAs, updatedshapeimg)
	if shapeimgconversion {
		fmt.Println(saveResultAs + " saved successfully")
	} else {
		log.Fatalf("Error in converting" + saveResultAs)
	}
	elapsed := time.Since(start)
	log.Printf("shapeidentifier took: %s", elapsed)
}

// markAndFindShapes creates a worker for each shape found in shapeimg.
// for each shape found, a text is added to the image.
// returns a gocv.Mat containing marked shapes.
func markAndFindShapes(shapeimg gocv.Mat) gocv.Mat {
	canny := gocv.NewMat()
	gocv.Canny(shapeimg, &canny, 10, 10)

	contours := gocv.FindContours(canny, gocv.RetrievalExternal, gocv.ChainApproxTC89L1)
	imgpoints := contours.ToPoints()
	amtOfJobs := contours.Size()

	jobs := make(chan int, amtOfJobs)
	result := make(chan Result, amtOfJobs)

	for amountOfRoutines := 0; amountOfRoutines < runtime.NumCPU()-1; amountOfRoutines++ {
		go worker(shapeimg, contours, imgpoints, jobs, result)
	}

	for i := 0; i < amtOfJobs; i++ {
		jobs <- i
	}
	close(jobs)

	for j := 0; j < amtOfJobs; j++ {
		shaperesult := <-result
		red := color.RGBA{255, 0, 0, 0}
		gocv.PutText(&shapeimg, shaperesult.Shape, shaperesult.Textpoint, 2, 0.75, red, 1)
		fmt.Printf("\t %s\n", shaperesult.Shape)
	}

	gocv.DrawContours(&shapeimg, contours, -1, color.RGBA{0, 0, 255, 0}, 1)
	return shapeimg
}

// worker calls detectshape with the image points of a single shape.
// the amount of jobs depends on the number of shapes found in markAndFindShapes.
// the detected shape is sent to the result channel.
// returns null.
func worker(shapeimg gocv.Mat, contours gocv.PointsVector, imgpoints [][]image.Point, jobs <-chan int, result chan<- Result) {
	for i := range jobs {
		shapeimgpointvector := imgpoints[i]
		shapevector := gocv.NewPointVectorFromPoints(shapeimgpointvector)

		shapedetectresult := detectshape(shapevector)
		result <- shapedetectresult
	}
}

// detectshape calculates the number of corners from a PointVector, containing points for a shape.
// If no shape is found, "unidentified shape" is put into Result.
// returns a Result containg the shape and a point for text.
func detectshape(pvr gocv.PointVector) Result {
	shape := "unidentified shape"
	shapeperimeter := gocv.ArcLength(pvr, true)
	shapeguess := gocv.ApproxPolyDP(pvr, 0.03*shapeperimeter, true)

	shapeguessRightType := shapeguess.ToPoints()

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
		if (p1p2-p2p3) < 10 && (p2p3-p3p4) < 10 && (p4p1-p1p2) < 10 {
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
	} // If more shapes needs to be detected, add them here.

	var result Result
	result.Shape = shape
	result.Textpoint = textpoint

	return result
}

// calculateDistanceBetweenTwoPoints calculates the distance between two image.Point.
// returns the distance as a float64
func calculateDistanceBetweenTwoPoints(point1 image.Point, point2 image.Point) float64 {
	return math.Sqrt(float64((point2.X-point1.X)*(point2.X-point1.X)) + float64((point2.Y-point1.Y)*(point2.Y-point1.Y)))
}

// isOctagon checks wether a shape is a circle or an octagon
// returns true if its an octagon, else false.
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
