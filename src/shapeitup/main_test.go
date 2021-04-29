package main

import (
	"image"
	"regexp"
	"testing"

	"gocv.io/x/gocv"
)

func TestTriangleShape(t *testing.T) {
	pvector := gocv.NewPointVector()
	pvector.Append(image.Point{100, 200})
	pvector.Append(image.Point{500, 300})
	pvector.Append(image.Point{44, 1212})
	name := "triangle"
	want := regexp.MustCompile(`\b` + name + `\b`)
	result := detectshape(pvector)

	if !want.MatchString(result.Shape) {
		t.Fatalf(`detectShape("pvector") = %q, want match for %#q, nil`, result.Shape, want)
	}
}
func TestSquareShape(t *testing.T) {
	pvector := gocv.NewPointVector()
	pvector.Append(image.Point{0, 0})
	pvector.Append(image.Point{1, 0})
	pvector.Append(image.Point{0, 1})
	pvector.Append(image.Point{1, 1})
	name := "square"
	want := regexp.MustCompile(`\b` + name + `\b`)
	result := detectshape(pvector)

	if !want.MatchString(result.Shape) {
		t.Fatalf(`detectShape("pvector") = %q, want match for %#q, nil`, result.Shape, want)
	}
}

func TestRectangleShape(t *testing.T) {
	pvector := gocv.NewPointVector()
	pvector.Append(image.Point{0, 0})
	pvector.Append(image.Point{200, 0})
	pvector.Append(image.Point{0, 100})
	pvector.Append(image.Point{200, 100})
	name := "rectangle"
	want := regexp.MustCompile(`\b` + name + `\b`)
	result := detectshape(pvector)

	if !want.MatchString(result.Shape) {
		t.Fatalf(`detectShape("pvector") = %q, want match for %#q, nil`, result.Shape, want)
	}
}

func TestPentagonShape(t *testing.T) {
	pvector := gocv.NewPointVector()
	pvector.Append(image.Point{0, 0})
	pvector.Append(image.Point{2, 0})
	pvector.Append(image.Point{2, 2})
	pvector.Append(image.Point{0, 2})
	pvector.Append(image.Point{1, 1})
	name := "pentagon"
	want := regexp.MustCompile(`\b` + name + `\b`)
	result := detectshape(pvector)

	if !want.MatchString(result.Shape) {
		t.Fatalf(`detectShape("pvector") = %q, want match for %#q, nil`, result.Shape, want)
	}
}

func TestHexagonShape(t *testing.T) {
	pvector := gocv.NewPointVector()
	pvector.Append(image.Point{0, 1})
	pvector.Append(image.Point{1, 0})
	pvector.Append(image.Point{2, 0})
	pvector.Append(image.Point{3, 1})
	pvector.Append(image.Point{2, 2})
	pvector.Append(image.Point{1, 2})
	name := "hexagon"
	want := regexp.MustCompile(`\b` + name + `\b`)
	result := detectshape(pvector)

	if !want.MatchString(result.Shape) {
		t.Fatalf(`detectShape("pvector") = %q, want match for %#q, nil`, result.Shape, want)
	}
}

func TestHeptaGonShape(t *testing.T) {
	pvector := gocv.NewPointVector()
	pvector.Append(image.Point{0, 1})
	pvector.Append(image.Point{1, 0})
	pvector.Append(image.Point{2, 0})
	pvector.Append(image.Point{3, 1})
	pvector.Append(image.Point{2, 2})
	pvector.Append(image.Point{1, 2})
	pvector.Append(image.Point{2, 3})
	name := "heptagon"
	want := regexp.MustCompile(`\b` + name + `\b`)
	result := detectshape(pvector)

	if !want.MatchString(result.Shape) {
		t.Fatalf(`detectShape("pvector") = %q, want match for %#q, nil`, result.Shape, want)
	}
}

func TestOctagonShape(t *testing.T) {
	pvector := gocv.NewPointVector()
	pvector.Append(image.Point{0, 1})
	pvector.Append(image.Point{1, 0})
	pvector.Append(image.Point{2, 0})
	pvector.Append(image.Point{3, 1})
	pvector.Append(image.Point{2, 2})
	pvector.Append(image.Point{1, 2})
	pvector.Append(image.Point{2, 3})
	pvector.Append(image.Point{3, 3})
	name := "octagon"
	want := regexp.MustCompile(`\b` + name + `\b`)
	result := detectshape(pvector)

	if !want.MatchString(result.Shape) {
		t.Fatalf(`detectShape("pvector") = %q, want match for %#q, nil`, result.Shape, want)
	}
}

func TestNonagonShape(t *testing.T) {
	pvector := gocv.NewPointVector()
	pvector.Append(image.Point{0, 1})
	pvector.Append(image.Point{1, 0})
	pvector.Append(image.Point{2, 0})
	pvector.Append(image.Point{3, 1})
	pvector.Append(image.Point{2, 2})
	pvector.Append(image.Point{1, 2})
	pvector.Append(image.Point{2, 3})
	pvector.Append(image.Point{3, 3})
	pvector.Append(image.Point{1, 3})
	name := "nonagon"
	want := regexp.MustCompile(`\b` + name + `\b`)
	result := detectshape(pvector)

	if !want.MatchString(result.Shape) {
		t.Fatalf(`detectShape("pvector") = %q, want match for %#q, nil`, result.Shape, want)
	}
}
