package aoi

type zAoiList struct {
	aoiDistance float32
	head        *xzAoi
	tail        *xzAoi
}

func newZAoiList(aoiDistance float32) *zAoiList {
	return &zAoiList{aoiDistance: aoiDistance}
}
func (sl *zAoiList) Insert(aoi *xzAoi) {
	insertCoordinate := aoi.aoi.z
	if sl.head != nil {
		p := sl.head
		for p != nil && p.aoi.z < insertCoordinate {
			p = p.zNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p == nil { // if p == nil, insert xzaoi at the end of list
			tail := sl.tail
			tail.zNext = aoi
			aoi.zPrev = tail
			sl.tail = aoi
		} else { // otherwise, p >= xzaoi, insert xzaoi before p
			prev := p.zPrev
			aoi.zNext = p
			p.zPrev = aoi
			aoi.zPrev = prev

			if prev != nil {
				prev.zNext = aoi
			} else { // p is the head, so xzaoi should be the new head
				sl.head = aoi
			}
		}
	} else {
		sl.head = aoi
		sl.tail = aoi
	}
}
func (sl *zAoiList) Remove(aoi *xzAoi) {
	prev := aoi.zPrev
	next := aoi.zNext
	if prev != nil {
		prev.zNext = next
		aoi.zPrev = nil
	} else {
		sl.head = next
	}
	if next != nil {
		next.zPrev = prev
		aoi.zNext = nil
	} else {
		sl.tail = prev
	}
}
func (sl *zAoiList) Move(aoi *xzAoi, oldCoordinate float32) {
	coordinate := aoi.aoi.z
	if coordinate > oldCoordinate {
		// moving to next ...
		next := aoi.zNext
		if next == nil || next.aoi.z >= coordinate {
			// no need to adjust in list
			return
		}
		prev := aoi.zPrev
		//fmt.Println(1, prev, next, prev == nil || prev.zNext == xzaoi)
		if prev != nil {
			prev.zNext = next // remove xzaoi from list
		} else {
			sl.head = next // xzaoi is the head, trim it
		}
		next.zPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.zNext == next)
		prev, next = next, next.zNext
		for next != nil && next.aoi.z < coordinate {
			prev, next = next, next.zNext
			//fmt.Println(2, prev, next, prev == nil || prev.zNext == next)
		}
		//fmt.Println(3, prev, next)
		// no we have prev.X < coord && (next == nil || next.X >= coord), so insert between prev and next
		prev.zNext = aoi
		aoi.zPrev = prev
		if next != nil {
			next.zPrev = aoi
		} else {
			sl.tail = aoi
		}
		aoi.zNext = next

		//fmt.Println(4)
	} else {
		// moving to prev ...
		prev := aoi.zPrev
		if prev == nil || prev.aoi.z <= coordinate {
			// no need to adjust in list
			return
		}

		next := aoi.zNext
		if next != nil {
			next.zPrev = prev
		} else {
			sl.tail = prev // xzaoi is the head, trim it
		}
		prev.zNext = next // remove xzaoi from list

		next, prev = prev, prev.zPrev
		for prev != nil && prev.aoi.z > coordinate {
			next, prev = prev, prev.zPrev
		}
		// no we have next.X > coord && (prev == nil || prev.X <= coord), so insert between prev and next
		next.zPrev = aoi
		aoi.zNext = next
		if prev != nil {
			prev.zNext = aoi
		} else {
			sl.head = aoi
		}
		aoi.zPrev = prev
	}
}
func (sl *zAoiList) Mark(aoi *xzAoi) {
	prev := aoi.zPrev
	coordinate := aoi.aoi.z

	minCoordinate := coordinate - sl.aoiDistance
	for prev != nil && prev.aoi.z >= minCoordinate {
		prev.markVal += 1
		prev = prev.zPrev
	}

	next := aoi.zNext
	maxCoordinate := coordinate + sl.aoiDistance
	for next != nil && next.aoi.z <= maxCoordinate {
		next.markVal += 1
		next = next.zNext
	}
}
func (sl *zAoiList) ClearMark(aoi *xzAoi) {
	prev := aoi.zPrev
	coordinate := aoi.aoi.z

	minCoordinate := coordinate - sl.aoiDistance
	for prev != nil && prev.aoi.z >= minCoordinate {
		prev.markVal = 0
		prev = prev.zPrev
	}

	next := aoi.zNext
	maxCoordinate := coordinate + sl.aoiDistance
	for next != nil && next.aoi.z <= maxCoordinate {
		next.markVal = 0
		next = next.zNext
	}
}
