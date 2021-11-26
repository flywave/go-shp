package shp

/*
#cgo linux LDFLAGS: -L ./  -Wl,--start-group  -lm -Wl,--end-group
#cgo darwin LDFLAGS: -L /  -lm
#cgo darwin,arm LDFLAGS: -L / -lm
#cgo windows LDFLAGS: -L ./  -lm -fPIC
#include <shapefil.h>
#include <string.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

type (
	SHPHandle C.SHPHandle
	DBFHandle C.DBFHandle
	SHPObject struct {
		ShapeType    C.int
		ShapeId      C.int
		NParts       C.int
		PanPartStart *C.int
		PanPartType  *C.int
		NVertices    C.int
		PadfX        *C.double
		PadfY        *C.double
		PadfZ        *C.double
		PadfM        *C.double
		XMin         C.double
		YMin         C.double
		ZMin         C.double
		MMin         C.double
		XMax         C.double
		YMax         C.double
		ZMax         C.double
		MMax         C.double
	}
)

func goSHPOpen(filename, mode string) SHPHandle {
	filename_, mode_ := C.CString(filename), C.CString(mode)
	defer C.free(unsafe.Pointer(filename_))
	defer C.free(unsafe.Pointer(mode_))
	return SHPHandle(C.goSHPOpen(filename_, mode_))
}

func goSHPClose(h SHPHandle) {
	C.goSHPClose(h)
}

func goSHPGetInfo(hSHP SHPHandle) (shapeType int, nEntries int, minBound, maxBound [4]C.double) {
	var shapeType_, nEntries_ C.int
	C.goSHPGetInfo(hSHP,
		&nEntries_,
		&shapeType_,
		&minBound[0],
		&maxBound[0])
	shapeType, nEntries = int(shapeType_), int(nEntries_)
	return
}

func goSHPReadObject(hSHP SHPHandle, iShape int) *SHPObject {
	return (*SHPObject)(unsafe.Pointer(C.goSHPReadObject(hSHP, C.int(iShape))))
}

func goSHPDestroyObject(o *SHPObject) {
	C.goSHPDestroyObject((*C.SHPObject)(unsafe.Pointer(o)))
}

func goDBFOpen(filename, mode string) DBFHandle {
	filename_, mode_ := C.CString(filename), C.CString(mode)
	defer C.free(unsafe.Pointer(filename_))
	defer C.free(unsafe.Pointer(mode_))
	return DBFHandle(C.goDBFOpen(filename_, mode_))
}

func goDBFClose(h DBFHandle) {
	C.goDBFClose(h)
}

func goDBFGetFieldCount(h DBFHandle) int {
	return int(C.goDBFGetFieldCount(h))
}

func goDBFGetRecordCount(h DBFHandle) int {
	return int(C.goDBFGetRecordCount(h))
}

func goDBFGetFieldIndex(h DBFHandle, fieldName string) int {
	fieldName_ := C.CString(fieldName)
	defer C.free(unsafe.Pointer(fieldName_))
	return int(C.goDBFGetFieldIndex(h, fieldName_))
}

func goDBFGetFieldInfo(h DBFHandle, fieldIndex int) (fieldName string, fieldType int, nWidth, nDecimals int) {
	var fieldName_ [12]C.char
	var nWidth_, nDecimals_ C.int
	fieldType = int(C.goDBFGetFieldInfo(h, C.int(fieldIndex), &fieldName_[0], &nWidth_, &nDecimals_))
	fieldName = C.GoString(&fieldName_[0])
	nWidth = int(nWidth_)
	nDecimals = int(nDecimals_)
	return
}

func goDBFReadIntegerAttribute(h DBFHandle, shapeIndex, fieldIndex int) int {
	return int(C.goDBFReadIntegerAttribute(h, C.int(shapeIndex), C.int(fieldIndex)))
}

func goDBFReadDoubleAttribute(h DBFHandle, shapeIndex, fieldIndex int) float64 {
	return float64(C.goDBFReadDoubleAttribute(h, C.int(shapeIndex), C.int(fieldIndex)))
}

func goDBFReadStringAttribute(h DBFHandle, shapeIndex, fieldIndex int) []byte {
	cstr := C.goDBFReadStringAttribute(h, C.int(shapeIndex), C.int(fieldIndex))
	slen := int(C.strlen(cstr))
	bytes := make([]byte, slen)
	if slen > 0 {
		C.memcpy(unsafe.Pointer(&bytes[0]), unsafe.Pointer(cstr), C.size_t(slen))
	}
	return bytes
}
