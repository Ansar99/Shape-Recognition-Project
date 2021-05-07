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
	Shape     string        // The shape
	Textpoint image.Point   // The point where text is placed
	Vertices  []image.Point //ANTON
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: go run main.go selected_image.jpg resulting_image_save.jpg")
		return
	}
	//runtime.GOMAXPROCS(4) //ANTON
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
	gocv.IMWrite("./shapedImages/4canny1010.jpg", canny) //ANTON
	contours := gocv.FindContours(canny, gocv.RetrievalList, gocv.ChainApproxTC89L1)
	imgpoints := contours.ToPoints()
	amtOfJobs := contours.Size()

	jobs := make(chan int, amtOfJobs)
	result := make(chan Result, amtOfJobs)
	//fmt.Println(runtime.NumCPU())
	for amountOfRoutines := 0; amountOfRoutines < runtime.NumCPU()-1; amountOfRoutines++ {
		go worker(shapeimg, contours, imgpoints, jobs, result)
	}

	for i := 0; i < amtOfJobs; i++ { // i = i + 2 TODO
		jobs <- i
	}
	close(jobs)

	for j := 0; j < amtOfJobs; j++ { // j = j + 2 TODO
		shaperesult := <-result
		if shaperesult.Shape == "bad" { //ANTON
			continue
		}
		red := color.RGBA{255, 0, 0, 0}
		gocv.PutText(&shapeimg, shaperesult.Shape, shaperesult.Textpoint, 2, 0.75, red, 1)
		for i, x := range shaperesult.Vertices { //ANTON
			//gocv.Circle(&shapeimg, x, 5, red, 10)
			//gocv.PutText(&shapeimg, fmt.Sprint(i), x, 2, 0.75, red, 1)
			if i < len(shaperesult.Vertices)-1 {
				gocv.Line(&shapeimg, x, shaperesult.Vertices[(i+4)%len(shaperesult.Vertices)], red, 2)
			}
		}

		//fmt.Printf("\t %s\n", shaperesult.Shape)
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
	//Här borde man kunna sortera ut borde göra i förhållande bildstorleken, ANTON
	shapeperimeter := gocv.ArcLength(pvr, true)
	if shapeperimeter < 200 { //ANTON
		var resultbad Result
		resultbad.Shape = "bad"
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
		shape = fourVertices(shapeguessRightType) //ANTON
	} else if vertices == 5 {
		shape = "pentagon"
	} else if vertices == 6 {
		shape = "hexagon"
	} else if vertices == 7 {
		shape = "heptagon"
	} else if vertices == 8 {
		shape = isOctagon(shapeguessRightType, shapeperimeter)
		fmt.Println(shape)
	} else if vertices == 9 {
		shape = "nonagon"
	} else {
		shape = "unknown"
	} // If more shapes needs to be detected, add them here.

	var result Result
	result.Shape = shape
	result.Textpoint = textpoint
	result.Vertices = shapeguessRightType
	return result
}

// calculateDistanceBetweenTwoPoints calculates the distance between two image.Point.
// returns the distance as a float64
func calculateDistanceBetweenTwoPoints(point1 image.Point, point2 image.Point) float64 {
	return math.Sqrt(float64((point2.X-point1.X)*(point2.X-point1.X)) + float64((point2.Y-point1.Y)*(point2.Y-point1.Y)))
}

// func threePointAngle(pointA image.Point, pointB image.Point, pointC image.Point) float64 {
// 	return math.Atan2(float64(pointC.X-pointA.X), float64(pointC.Y-pointA.Y)) - math.Atan2(float64(pointB.X-pointA.X), float64(pointB.Y-pointA.Y))
// }

func dotAngle(pointA image.Point, pointB image.Point, pointC image.Point) float64 {
	a := [2]float64{float64(pointA.X - pointB.X), float64(pointA.Y - pointB.Y)} //FULt ändra ANTON
	b := [2]float64{float64(pointA.X - pointC.X), float64(pointA.Y - pointC.Y)}
	return math.Acos((a[0]*b[0] + a[1]*b[1]) / (calculateDistanceBetweenTwoPoints(pointA, pointB) * calculateDistanceBetweenTwoPoints(pointA, pointC)))
}

func parallell(pointA image.Point, pointB image.Point, pointC image.Point, pointD image.Point) bool {
	kA := float64(pointA.Y-pointB.Y) / float64(pointA.X-pointB.X)
	kB := float64(pointC.Y-pointD.Y) / float64(pointC.X-pointD.X)
	return math.Abs(kA-kB) < 0.1 //FIXME: Vettigt nr
}

func fourVertices(points []image.Point) string { //ANTON
	// Points (X,Y) of the qudriangle
	p1 := points[0]
	p2 := points[1]
	p3 := points[2]
	p4 := points[3]

	// Distance between points in quadriangle
	p1p2 := calculateDistanceBetweenTwoPoints(p1, p2)
	p2p3 := calculateDistanceBetweenTwoPoints(p2, p3)
	p3p4 := calculateDistanceBetweenTwoPoints(p3, p4)
	p4p1 := calculateDistanceBetweenTwoPoints(p4, p1)

	// Angles
	p1A := dotAngle(p1, p2, p4)
	//p2A := dotAngle(p2, p1, p3)
	p3A := dotAngle(p3, p2, p4)
	//p4A := dotAngle(p4, p3, p1)

	if math.Abs(p1p2-p2p3) < 10 && math.Abs(p2p3-p3p4) < 10 && math.Abs(p4p1-p1p2) < 10 {
		return "square" //ANTON När returnerna romb ??
	} else if math.Abs(p1p2-p3p4) < 10 || math.Abs(p2p3-p4p1) < 10 {
		if math.Abs(p1A-p3A) < math.Pi/30 { //equal angle on opposite sides
			if math.Abs(p1A-math.Pi/2) < math.Pi/30 { //angle ~= 90
				return "rectangle"
			} else {
				return "parallelogram"
			}
		} else {
			return "paralleltrapets" //Hur fånga ?
		}
	} else if parallell(p1, p2, p3, p4) || parallell(p2, p3, p1, p4) {
		return "paralleltrapts"
	} else {
		return "quadriangle"
	}
}

// isOctagon checks wether a shape is a circle or an octagon
// returns true if its an octagon, else false.
func isOctagon(points []image.Point, shapeperimeter float64) string {
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

	length := len(points)
	var circumference float64 = 0
	for i, v := range points {
		if i < length-1 {
			circumference += calculateDistanceBetweenTwoPoints(v, points[i+1])
		} else {
			circumference += calculateDistanceBetweenTwoPoints(v, points[0])
		}
	}

	if math.Abs(shapeperimeter-circumference) < shapeperimeter*0.01 { //FIXME: 0.02 needs to be tested 0.032 ?
		fmt.Println("shapeperimeter for octagon: ", shapeperimeter)
		fmt.Println("circumference for octagon: ", circumference)
		return "octagon"
	} else if math.Abs(p1p5-p2p6) > 10 && math.Abs(p2p6-p3p7) > 10 && math.Abs(p3p7-p4p8) > 10 {
		fmt.Println("shapeperimeter for ovale: ", shapeperimeter)
		fmt.Println("circumference for ovale: ", circumference)
		return "ovale"
	} else {
		fmt.Println("shapeperimeter for circle: ", shapeperimeter)
		fmt.Println("circumference for circle: ", circumference)
		return "circle"
	}
}

// imageToGrayScaleMat converts a path of an image to a gocv.Mat, grayscales it and blurs it.
// if the path is invalid an error occours.
// returns a gocv.Mat or an error.
func imageToGrayscaleMat(imgname string) (gocv.Mat, error) {
	img := gocv.IMRead(imgname, gocv.IMReadGrayScale)
	if img.Empty() {
		return gocv.Mat{}, errors.New("image img in imageToGrayscaleMat is empty")
	}
	gocv.IMWrite("./shapedImages/1greyimage.jpg", img) //ANTON
	gocv.MedianBlur(img, &img, 11)
	gocv.IMWrite("./shapedImages/2blurimage.jpg", img) //ANTON
	shapeimg := gocv.NewMat()
	gocv.CvtColor(img, &shapeimg, gocv.ColorGrayToBGR)
	gocv.IMWrite("./shapedImages/3greyblurimage.jpg", shapeimg) //ANTON
	if shapeimg.Empty() {
		return gocv.Mat{}, errors.New("image shapeimg in imageToGrayscaleMat is empty")
	}

	return shapeimg, nil
}

///go run main.go ./images/image.jpg ./shapedImages/shapedimage.jpg ANTON
