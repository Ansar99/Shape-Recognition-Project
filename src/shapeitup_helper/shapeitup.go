package helper

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"

	"gocv.io/x/gocv"
)

// A Result represents a Shape string, a Textpoint where Shape should be placed and a slice of points representing the Shape's Contour.
type Result struct {
	Shape     string        // The shape
	Textpoint image.Point   // The point where text is placed
	Contour   []image.Point //The points for the countor of the detected shape
}

// markAndFindShapes creates a worker for each shape found in shapeimg.
// for each shape found, a text is added to the image.
// returns a gocv.Mat containing marked shapes.
func MarkAndFindShapes(shapeimg gocv.Mat) gocv.Mat {
	canny := gocv.NewMat()
	gocv.Canny(shapeimg, &canny, 10, 10)

	contours := gocv.FindContours(canny, gocv.RetrievalExternal, gocv.ChainApproxTC89L1)
	imgpoints := contours.ToPoints()
	amtOfJobs := contours.Size()

	jobs := make(chan int, amtOfJobs)
	result := make(chan Result, amtOfJobs)

	for amountOfRoutines := 0; amountOfRoutines < runtime.NumCPU()-1; amountOfRoutines++ {
		go worker(imgpoints, jobs, result)
	}

	for i := 0; i < amtOfJobs; i++ {
		jobs <- i
	}
	close(jobs)

	for j := 0; j < amtOfJobs; j++ {
		shaperesult := <-result
		if shaperesult.Shape == "unidentified" {
			continue
		}
		red := color.RGBA{255, 0, 0, 0}
		gocv.PutText(&shapeimg, shaperesult.Shape, shaperesult.Textpoint, 2, 0.75, red, 1)
		contour := [][]image.Point{shaperesult.Contour}
		gocv.DrawContours(&shapeimg, gocv.NewPointsVectorFromPoints(contour), -1, color.RGBA{0, 0, 255, 0}, 1)
		fmt.Printf("%s\n", shaperesult.Shape)
	}
	return shapeimg
}

// worker calls detectshape with the image points of a single shape.
// the amount of jobs depends on the number of shapes found in markAndFindShapes.
// the detected shape is sent to the result channel.
// returns null.
func worker(imgpoints [][]image.Point, jobs <-chan int, result chan<- Result) {
	for i := range jobs {
		shapeimgpointvector := imgpoints[i]
		shapevector := gocv.NewPointVectorFromPoints(shapeimgpointvector)
		shapedetectresult := detectshape(shapevector, shapeimgpointvector)
		result <- shapedetectresult
	}
}

// detectshape calculates the number of corners from a PointVector, containing points for a shape.
// If no shape is found, "unidentified shape" is put into Result.
// returns a Result containg the shape and a point for text.
func detectshape(pvr gocv.PointVector, shapeimgpointvector []image.Point) Result {
	shape := "unidentified shape"
	shapeperimeter := gocv.ArcLength(pvr, true)
	if shapeperimeter < 200 {
		var resultbad Result
		resultbad.Shape = "unidentified"
		return resultbad
	}
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
		shape = isOctagon(shapeguessRightType, shapeperimeter, pvr)
	} else if vertices == 9 {
		shape = "nonagon"
	} else {
		shape = "unidentified"
	} // If more shapes needs to be detected, add them here.

	var result Result
	result.Shape = shape
	result.Textpoint = textpoint
	result.Contour = shapeimgpointvector

	return result
}

// calculateDistanceBetweenTwoPoints calculates the distance between two image.Point.
// returns the distance as a float64
func calculateDistanceBetweenTwoPoints(point1 image.Point, point2 image.Point) float64 {
	return math.Sqrt(float64((point2.X-point1.X)*(point2.X-point1.X)) + float64((point2.Y-point1.Y)*(point2.Y-point1.Y)))
}

// isOctagon checks wether a shape is a circle, oval or octagon
// returns a string of the shape detected
func isOctagon(points []image.Point, shapeperimeter float64, pvr gocv.PointVector) string {
	p1 := points[0]
	p2 := points[1]
	p3 := points[2]
	p4 := points[3]
	p5 := points[4]
	p6 := points[5]
	p7 := points[6]
	p8 := points[7]

	p1p5 := calculateDistanceBetweenTwoPoints(p1, p5)
	p2p6 := calculateDistanceBetweenTwoPoints(p2, p6)
	p3p7 := calculateDistanceBetweenTwoPoints(p3, p7)
	p4p8 := calculateDistanceBetweenTwoPoints(p4, p8)

	if len(gocv.ApproxPolyDP(pvr, 0.01*shapeperimeter, true).ToPoints()) > 8 {
		if math.Abs(p1p5-p2p6) > 10 && math.Abs(p2p6-p3p7) > 10 && math.Abs(p3p7-p4p8) > 10 {
			return "ovale"
		} else {
			return "circle"
		}
	} else {
		return "octagon"
	}
}
