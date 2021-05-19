package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestTriangle(t *testing.T) {
	name := [8]string{"triangle", "triangle", "triangle", "triangle", "triangle", "triangle", "triangle", "triangle"}
	result := run("testImages/8triangle.jpg", "shapedTestImages/8triangle.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestTriangle failed")
		}
	}
}

func TestSquare(t *testing.T) {
	name := [5]string{"square", "square", "square", "square", "square"}
	result := run("testImages/5square.jpg", "shapedTestImages/5square.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestSquare failed")
		}
	}
}

func TestRectangle(t *testing.T) {
	name := [7]string{"rectangle", "rectangle", "rectangle", "rectangle", "rectangle", "rectangle", "rectangle"}
	result := run("testImages/7rectangle.jpg", "shapedTestImages/7rectangle.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestRectangle failed")
		}
	}
}

func TestPentagon(t *testing.T) {
	name := [5]string{"pentagon", "pentagon", "pentagon", "pentagon", "pentagon"}
	result := run("testImages/5pentagon.jpg", "shapedTestImages/5pentagon.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestPentagon failed")
		}
	}
}

func TestHexagon(t *testing.T) {
	name := [5]string{"hexagon", "hexagon", "hexagon", "hexagon", "hexagon"}
	result := run("testImages/5hexagon.jpg", "shapedTestImages/5hexagon.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestHexagon failed")
		}
	}
}

func TestHeptagon(t *testing.T) {
	name := [5]string{"heptagon", "heptagon", "heptagon", "heptagon", "heptagon"}
	result := run("testImages/5heptagon.jpg", "shapedTestImages/5heptagon.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestHeptagon failed")
		}
	}
}

func TestOctagon(t *testing.T) {
	name := [6]string{"octagon", "octagon", "octagon", "octagon", "octagon", "octagon"}
	result := run("testImages/6octagon.jpg", "shapedTestImages/6octagon.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestOctagon failed")
		}
	}
}

func TestNonagon(t *testing.T) {
	name := [4]string{"nonagon", "nonagon", "nonagon", "nonagon"}
	result := run("testImages/4nonagon.jpg", "shapedTestImages/4nonagon.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestNonagon failed")
		}
	}
}

func TestCircle(t *testing.T) {
	name := [4]string{"circle", "circle", "circle", "circle"}
	result := run("testImages/4circle.jpg", "shapedTestImages/4circle.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestCircle failed")
		}
	}
}

func TestOvale(t *testing.T) {
	name := [5]string{"ovale", "ovale", "ovale", "ovale", "ovale"}
	result := run("testImages/5ovale.jpg", "shapedTestImages/5ovale.jpg")
	resultslice := strings.Split(result, ", ")

	for _, v := range name {
		containsBool, containsIndex := contains(resultslice, v)
		if containsBool {
			resultslice = removeIndex(resultslice, containsIndex)
			fmt.Println(resultslice)
			continue
		} else {
			t.Fatalf("run(imagepath, imageoutputpath) in TestOvale failed")
		}
	}
}

// HELPER FUNCTION
//
// Returns true if the selected string is an element of the slice.
func contains(s []string, str string) (bool, int) {
	for i, v := range s {
		if v == str {
			return true, i
		}
	}

	return false, -1
}

// HELPER FUNCTION
//
// Removes an element at index index in slice s
func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
