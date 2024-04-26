package shp

import (
	"math"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
)

func shpToGeomFeature(shp *Shape) *geom.Feature {
	g := shpToGeometry(shp)
	box := shpToGeomBoundingBox(shp)
	return &geom.Feature{ID: uint64(shp.Id), BoundingBox: &box, Geometry: g, Properties: shp.Attrs}
}

func shpToGeomBoundingBox(shp *Shape) geom.BoundingBox {
	switch shp.Type {
	case ShapePoint, ShapePointM, ShapeMultiPoint, ShapeMultiPointM, ShapePolyLine, ShapePolyLineM, ShapePolygon, ShapePolygonM:
		return geom.BoundingBox{{shp.Box.Min.X, shp.Box.Min.Y, -math.MaxFloat64}, {shp.Box.Max.X, shp.Box.Max.Y, math.MaxFloat64}}
	case ShapePointZ, ShapeMultiPointZ, ShapePolyLineZ, ShapePolygonZ:
		return geom.BoundingBox{{shp.Box.Min.X, shp.Box.Min.Y, shp.Box.Min.Z}, {shp.Box.Max.X, shp.Box.Max.Y, shp.Box.Max.Z}}
	}
	return geom.BoundingBox{}
}

func shpToGeometry(shp *Shape) geom.Geometry {
	switch shp.Type {
	case ShapePoint, ShapePointM:
		return shpToPoint(shp)
	case ShapePointZ:
		return shpToPoint3(shp)
	case ShapeMultiPoint, ShapeMultiPointM:
		return shpToMultiPoint(shp)
	case ShapeMultiPointZ:
		return shpToMultiPoint3(shp)
	case ShapePolyLine, ShapePolyLineM:
		return shpToLineString(shp)
	case ShapePolyLineZ:
		return shpToLineString3(shp)
	case ShapePolygon, ShapePolygonM:
		return shpToPolygon(shp)
	case ShapePolygonZ:
		return shpToPolygon3(shp)
	}
	return nil
}

func shpToPoint(shp *Shape) geom.Point {
	return general.NewPoint([]float64{shp.Parts[0].Vertices[0].X, shp.Parts[0].Vertices[0].Y})
}

func shpToPoint3(shp *Shape) geom.Point3 {
	return general.NewPoint3([]float64{shp.Parts[0].Vertices[0].X, shp.Parts[0].Vertices[0].Y, shp.Parts[0].Vertices[0].Z})
}

func shpToMultiPoint(shp *Shape) geom.MultiPoint {
	pts := [][]float64{}
	for _, par := range shp.Parts[0].Vertices {
		pts = append(pts, []float64{par.X, par.Y})
	}
	return general.NewMultiPoint(pts)
}

func shpToMultiPoint3(shp *Shape) geom.MultiPoint3 {
	pts := [][]float64{}
	for _, par := range shp.Parts[0].Vertices {
		pts = append(pts, []float64{par.X, par.Y, par.Z})
	}
	return general.NewMultiPoint3(pts)
}

func shpToLineString(shp *Shape) geom.LineString {
	pts := [][]float64{}
	for _, par := range shp.Parts[0].Vertices {
		pts = append(pts, []float64{par.X, par.Y})
	}
	return general.NewLineString(pts)
}

func shpToLineString3(shp *Shape) geom.LineString3 {
	pts := [][]float64{}
	for _, par := range shp.Parts[0].Vertices {
		pts = append(pts, []float64{par.X, par.Y, par.Z})
	}
	return general.NewLineString3(pts)
}

func shpToPolygon(shp *Shape) geom.Polygon {
	plines := [][][]float64{}
	for _, par := range shp.Parts {
		if par.Type == PartOuterRing || par.Type == PartFirstRing {
			lines := [][]float64{}
			for _, vers := range par.Vertices {
				lines = append(lines, []float64{vers.X, vers.Y})
			}
			plines = append([][][]float64{lines}, plines...)
		} else {
			lines := [][]float64{}
			for _, vers := range par.Vertices {
				lines = append(lines, []float64{vers.X, vers.Y})
			}
			plines = append(plines, lines)
		}
	}
	return general.NewPolygon(plines)
}

func shpToPolygon3(shp *Shape) geom.Polygon3 {
	plines := [][][]float64{}
	for _, par := range shp.Parts {
		if par.Type == PartOuterRing || par.Type == PartFirstRing {
			lines := [][]float64{}
			for _, vers := range par.Vertices {
				lines = append(lines, []float64{vers.X, vers.Y, vers.Z})
			}
			plines = append([][][]float64{lines}, plines...)
		} else {
			lines := [][]float64{}
			for _, vers := range par.Vertices {
				lines = append(lines, []float64{vers.X, vers.Y, vers.Z})
			}
			plines = append(plines, lines)
		}
	}
	return general.NewPolygon3(plines)
}
