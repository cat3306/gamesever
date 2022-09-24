package aoi

type xzAoi struct {
	aoi          *AOI
	neighbors    map[*xzAoi]struct{}
	xPrev, xNext *xzAoi
	zPrev, zNext *xzAoi
	markVal      int
}

type XZListAOIManager struct {
	aoIDistance float32
	xSweepList  *xAoiList
	zSweepList  *zAoiList
}

func (aoiMan *XZListAOIManager) adjust(aoi *xzAoi) {
	aoiMan.xSweepList.Mark(aoi)
	aoiMan.zSweepList.Mark(aoi)
	// AOI marked twice are neighbors
	for neighbor := range aoi.neighbors {
		if neighbor.markVal == 2 {
			// neighbors kept
			neighbor.markVal = -2 // mark this as neighbor
		} else { // markVal < 2
			// was neighbor, but not any more
			delete(aoi.neighbors, neighbor)
			aoi.aoi.callback.OnLeaveAOI(neighbor.aoi)
			delete(neighbor.neighbors, aoi)
			neighbor.aoi.callback.OnLeaveAOI(aoi.aoi)
		}
	}

	// travel in X list again to find all new neighbors, whose markVal == 2
	aoiMan.xSweepList.GetClearMarkedNeighbors(aoi)
	// travel in Z list again to unmark all
	aoiMan.zSweepList.ClearMark(aoi)
}
func (aoiMan *XZListAOIManager) Moved(aoi *AOI, x, z float32) {
	oldX := aoi.x
	oldY := aoi.z
	aoi.x, aoi.z = x, z
	xzAoi := aoi.implData.(*xzAoi)
	if oldX != x {
		aoiMan.xSweepList.Move(xzAoi, oldX)
	}
	if oldY != z {
		aoiMan.zSweepList.Move(xzAoi, oldY)
	}
	aoiMan.adjust(xzAoi)
}

func (aoiMan *XZListAOIManager) Leave(aoi *AOI) {
	xzAoi := aoi.implData.(*xzAoi)
	aoiMan.xSweepList.Remove(xzAoi)
	aoiMan.zSweepList.Remove(xzAoi)
	aoiMan.adjust(xzAoi)
}

func (aoiMan *XZListAOIManager) Enter(aoi *AOI, x, z float32) {
	aoi.dist = aoiMan.aoIDistance

	xzAoi := &xzAoi{
		aoi:       aoi,
		neighbors: map[*xzAoi]struct{}{},
	}
	aoi.x, aoi.z = x, z
	aoi.implData = xzAoi
	aoiMan.xSweepList.Insert(xzAoi)
	aoiMan.zSweepList.Insert(xzAoi)
	aoiMan.adjust(xzAoi)
}
func NewXZListAOIManager(aoIDistance float32) Manager {
	return &XZListAOIManager{
		aoIDistance: aoIDistance,
		xSweepList:  newXAoiList(aoIDistance),
		zSweepList:  newZAoiList(aoIDistance),
	}
}
func NewXZListAOIManagerV2(aoIDistance float32) *XZListAOIManager {
	return &XZListAOIManager{
		aoIDistance: aoIDistance,
		xSweepList:  newXAoiList(aoIDistance),
		zSweepList:  newZAoiList(aoIDistance),
	}
}
