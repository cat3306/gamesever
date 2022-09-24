package aoi

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const (
	MinX = -500
	MaxX = 500
	MinY = -500
	MaxY = 500
)

type TestObj struct {
	aoi            AOI
	Id             int
	neighbors      map[*TestObj]struct{}
	totalNeighbors int64
	nCalc          int64
}

func (obj *TestObj) getObj(aoi *AOI) *TestObj {
	return aoi.Data.(*TestObj)
}
func (obj *TestObj) OnLeaveAOI(otherAoi *AOI) {
	other := obj.getObj(otherAoi)
	if obj == other {
		panic("should not leave self")
	}
	if _, ok := obj.neighbors[other]; !ok {
		panic("duplicate leave aoi")
	}
	delete(obj.neighbors, other)
	obj.totalNeighbors += int64(len(obj.neighbors))
	obj.nCalc += 1
}
func (obj *TestObj) OnEnterAOI(otherAoi *AOI) {
	other := obj.getObj(otherAoi)
	if obj == other {
		panic("should not enter self")
	}
	if _, ok := obj.neighbors[other]; ok {
		panic("should not enter self")
	}
	obj.neighbors[other] = struct{}{}
	obj.totalNeighbors += int64(len(obj.neighbors))
	obj.nCalc += 1

}
func randCoordinate(min, max float32) float32 {
	return min + float32(rand.Intn(int(max)-int(min)))
}
func TestXZListAoiManager(t *testing.T) {
	testAoi(1000, NewXZListAOIManager(100), t)
}
func TestEnter(t *testing.T) {
	manger := NewXZListAOIManagerV2(100)
	obj := &TestObj{Id: 1, neighbors: map[*TestObj]struct{}{}}
	InitAOI(&obj.aoi, 100, obj, obj)
	manger.Enter(&obj.aoi, 0, 0)
	obj1 := &TestObj{Id: 2, neighbors: map[*TestObj]struct{}{}}
	InitAOI(&obj1.aoi, 100, obj1, obj1)
	manger.Enter(&obj1.aoi, 1, 1)
	fmt.Println(manger.xSweepList.head.aoi.Data.(*TestObj).Id)
	fmt.Println(manger.xSweepList.head.xNext.aoi.Data.(*TestObj).Id)
}
func testAoi(numAOI int, aoiMan Manager, t *testing.T) {
	var objs []*TestObj
	for i := 0; i < numAOI; i++ {
		obj := &TestObj{Id: i + 1, neighbors: map[*TestObj]struct{}{}}
		InitAOI(&obj.aoi, 100, obj, obj)
		objs = append(objs, obj)
		aoiMan.Enter(&obj.aoi, randCoordinate(MinX, MaxX), randCoordinate(MinY, MaxY))
	}
	for i := 0; i < 10; i++ {
		t0 := time.Now()
		for _, obj := range objs {
			aoiMan.Moved(&obj.aoi, obj.aoi.x+randCoordinate(-10, 10), obj.aoi.z+randCoordinate(-10, 10))
			aoiMan.Leave(&obj.aoi)
			aoiMan.Enter(&obj.aoi, obj.aoi.x+randCoordinate(-10, 10), obj.aoi.z+randCoordinate(-10, 10))
		}
		dt := time.Now().Sub(t0)
		t.Logf("%s tick %d objects takes %s", "XZListAOI", numAOI, dt)
	}

	for _, obj := range objs {
		aoiMan.Leave(&obj.aoi)
	}
	totalCalc := int64(0)
	for _, obj := range objs {
		totalCalc += obj.nCalc
	}
	println("Average calculate count:", totalCalc/int64(len(objs)))

}
