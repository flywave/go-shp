package shp

import (
	"fmt"

	"github.com/flywave/go-geom"
)

type ShapeType int

const (
	ShapeNull        ShapeType = 0
	ShapePoint       ShapeType = 1
	ShapePolyLine    ShapeType = 3
	ShapePolygon     ShapeType = 5
	ShapeMultiPoint  ShapeType = 8
	ShapePointZ      ShapeType = 11
	ShapePolyLineZ   ShapeType = 13
	ShapePolygonZ    ShapeType = 15
	ShapeMultiPointZ ShapeType = 18
	ShapePointM      ShapeType = 21
	ShapePolyLineM   ShapeType = 23
	ShapePolygonM    ShapeType = 25
	ShapeMultiPointM ShapeType = 28
	ShapeMultiPatch  ShapeType = 31
)

type PartType int

const (
	PartTriangleStrip PartType = 0
	PartTriangleFan   PartType = 1
	PartOuterRing     PartType = 2
	PartInnerRing     PartType = 3
	PartFirstRing     PartType = 4
	PartRing          PartType = 5
)

type FieldType int

const (
	String FieldType = iota
	Integer
	Double
	Logical
	Invalid
)

type Part struct {
	Type     PartType
	Vertices []Point
}

func (p Part) XY() (x, y []float64) {
	x, y = make([]float64, len(p.Vertices)), make([]float64, len(p.Vertices))
	for i := range p.Vertices {
		x[i], y[i] = p.Vertices[i].X, p.Vertices[i].Y
	}
	return
}

type Shape struct {
	Type  ShapeType
	Id    int
	Parts []Part
	Box   BoundingBox
	Attrs map[string]interface{}
}

type ShapeFile struct {
	hShape     SHPHandle
	Type       ShapeType
	ShapeCount int
	Box        BoundingBox

	hDb        DBFHandle
	FieldCount int
}

type Point struct {
	X, Y, Z, M float64
}

type BoundingBox struct {
	Min Point
	Max Point
}

func Open(file string) *ShapeFile {
	hShape := goSHPOpen(file, "rb")
	if hShape == nil {
		panic("Cannot open shape file " + file)
	}

	hDb := goDBFOpen(file, "rb")
	if hDb == nil {
		panic("Cannot open db file " + file)
	}

	shapeType, nEntries, minBound, maxBound := goSHPGetInfo(hShape)
	box := BoundingBox{
		Min: Point{
			float64(minBound[0]),
			float64(minBound[1]),
			float64(minBound[2]),
			float64(minBound[3]),
		},
		Max: Point{
			float64(maxBound[0]),
			float64(maxBound[1]),
			float64(maxBound[2]),
			float64(maxBound[3]),
		},
	}

	fieldCount := goDBFGetFieldCount(hDb)
	recordCount := goDBFGetRecordCount(hDb)
	if recordCount != nEntries {
		panic("Shape count and db record count does not match.")
	}
	return &ShapeFile{hShape, ShapeType(shapeType), nEntries, box,
		hDb, fieldCount}
}

func (f *ShapeFile) Close() {
	goSHPClose(f.hShape)
	goDBFClose(f.hDb)
}

func (f *ShapeFile) Shape(shapeIndex int) *Shape {
	if shapeIndex >= f.ShapeCount {
		return nil
	}
	s := goSHPReadObject(f.hShape, shapeIndex)
	defer goSHPDestroyObject(s)

	return &Shape{
		Type: ShapeType(s.ShapeType),
		Id:   int(s.ShapeId),
		Box: BoundingBox{
			Min: Point{float64(s.XMin), float64(s.YMin), float64(s.ZMin), float64(s.MMin)},
			Max: Point{float64(s.XMax), float64(s.YMax), float64(s.ZMax), float64(s.MMax)},
		},
		Parts: func() []Part {
			vertices := func() []Point {
				vertices := make([]Point, s.NVertices)
				for i := range vertices {
					vertices[i].X = GetFloat(s.PadfX, i)
					vertices[i].Y = GetFloat(s.PadfY, i)
					vertices[i].Z = GetFloat(s.PadfZ, i)
					vertices[i].M = GetFloat(s.PadfM, i)
				}
				return vertices
			}()

			if s.NParts == 0 {
				return []Part{{Vertices: vertices}}
			}

			parts := make([]Part, s.NParts)
			for i := 1; i < len(parts); i++ {
				parts[i-1].Type = PartType(GetInt(s.PanPartType, i-1))
				start := GetInt(s.PanPartStart, i-1)
				end := GetInt(s.PanPartStart, i)
				if end > len(vertices) {
					end = len(vertices)
				}
				if start < len(vertices) {
					parts[i-1].Vertices = vertices[start:end]
				} else {
					return parts
				}
			}
			last := len(parts) - 1
			if last >= 0 {
				parts[last].Type = PartType(GetInt(s.PanPartStart, last))
				start := GetInt(s.PanPartStart, last)
				if start < len(vertices) {
					parts[last].Vertices = vertices[start:]
				} else {
					return parts
				}
			}
			return parts
		}(),
		Attrs: func() map[string]interface{} {
			attrs := map[string]interface{}{}
			for j := 0; j < f.FieldCount; j++ {
				name, type_, _, _ := goDBFGetFieldInfo(f.hDb, j)
				switch FieldType(type_) {
				case String:
					attrs[name] = string(goDBFReadStringAttribute(f.hDb, shapeIndex, j))
				case Integer, Logical:
					attrs[name] = goDBFReadIntegerAttribute(f.hDb, shapeIndex, j)
				case Double:
					attrs[name] = goDBFReadDoubleAttribute(f.hDb, shapeIndex, j)
				default:
					panic("Unkown field type.")
				}
			}
			return attrs
		}(),
	}
}

func (f *ShapeFile) Feature(shapeIndex int) *geom.Feature {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic: %s\n", p)
		}
	}()
	shp := f.Shape(shapeIndex)
	if shp == nil {
		return nil
	}
	feat := shpToGeomFeature(shp)
	return feat
}
