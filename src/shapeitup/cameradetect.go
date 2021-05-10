package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"

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
		updatedshapeimg := markAndFindShapes(shapeimg)
		window.IMShow(updatedshapeimg)
		//window.WaitKey(1)

		if window.WaitKey(1) >= 0 {
            break
        }
	}
}

// markAndFindShapes creates a worker for each shape found in shapeimg.
// for each shape found, a text is added to the image.
// returns a gocv.Mat containing marked shapes.
func markAndFindShapes(shapeimg gocv.Mat) gocv.Mat {
	canny := gocv.NewMat()
	gocv.Canny(shapeimg, &canny, 10, 10)

	contours := gocv.FindContours(canny, gocv.RetrievalList, gocv.ChainApproxTC89L1)
	imgpoints := contours.ToPoints()
	amtOfJobs := contours.Size()

	jobs := make(chan int, amtOfJobs)
	result := make(chan Result, amtOfJobs)
	//fmt.Println(runtime.NumCPU())
	for amountOfRoutines := 0; amountOfRoutines < runtime.NumCPU()-1; amountOfRoutines++ {
		go worker(shapeimg, contours, imgpoints, jobs, result)
	}

	for i := 0; i < amtOfJobs; i = i + 2 {
		jobs <- i
	}
	close(jobs)

	for j := 0; j < amtOfJobs; j = j + 2 {
		shaperesult := <-result
		red := color.RGBA{255, 0, 0, 0}
		if shaperesult.Shape == "bad" {
			continue
		}
		gocv.PutText(&shapeimg, shaperesult.Shape, shaperesult.Textpoint, 2, 0.75, red, 1)
		gocv.DrawContours(&shapeimg, contours, -1, color.RGBA{0, 255, 0, 0}, 1)
		//fmt.Printf("\t %s\n", shaperesult.Shape)
	}

	gocv.DrawContours(&shapeimg, contours, -1, color.RGBA{0, 255, 0, 0}, 1)
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

	if shapeperimeter < 500 {
		var resultbad Result
		resultbad.Shape = "bad"
		resultbad.Textpoint = shapeguessRightType[0]

		return resultbad
	}

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
		if (p1p2-p2p3) < 30 && (p2p3-p3p4) < 30 && (p4p1-p1p2) < 30 {
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
		shape = ""
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

func imageToGrayscaleMat(imgname gocv.Mat) (gocv.Mat, error) {
	gocv.MedianBlur(imgname, &imgname, 11)
	shapeimg := gocv.NewMat()
	gocv.BilateralFilter(imgname, &shapeimg, 10, float64(100), float64(100))
	//gocv.CvtColor(imgname, &shapeimg, gocv.ColorGrayToBGR)

	return shapeimg, nil
}
