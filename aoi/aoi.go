package aoi

type AOI struct {
	x    float32
	z    float32
	dist float32
	Data interface{}

	callback Callback
	implData interface{}
}

func InitAOI(aoi *AOI, dist float32, data interface{}, callback Callback) {
	aoi.dist = dist
	aoi.Data = data
	aoi.callback = callback
}

type Callback interface {
	OnEnterAOI(other *AOI)
	OnLeaveAOI(other *AOI)
}

type Manager interface {
	Enter(aoi *AOI, x, z float32)
	Leave(aoi *AOI)
	Moved(aoi *AOI, x, z float32)
}
