package shp

import "testing"

func TestReadPoint(t *testing.T) {
	shp := Open("./test_files/point.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}

func TestReadPointM(t *testing.T) {
	shp := Open("./test_files/pointm.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}

func TestReadPointZ(t *testing.T) {
	shp := Open("./test_files/pointz.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}

func TestReadMultiPoint(t *testing.T) {
	shp := Open("./test_files/multipoint.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}

func TestReadMultiPointZ(t *testing.T) {
	shp := Open("./test_files/multipointz.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}

func TestReadPolyLine(t *testing.T) {
	shp := Open("./test_files/polyline.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}

func TestReadPolyLineZ(t *testing.T) {
	shp := Open("./test_files/polylinez.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}

func TestReadPolygon(t *testing.T) {
	shp := Open("./test_files/polygon.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}

func TestReadPolygonZ(t *testing.T) {
	shp := Open("./test_files/polygonz.shp")

	if shp == nil || shp.ShapeCount == 0 {
		t.FailNow()
	}

	feat := shp.Feature(0)

	if feat == nil {
		t.FailNow()
	}

}
