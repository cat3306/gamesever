package protocol

type CodeType uint16

const (
	CodeNone    = CodeType(0)
	String      = CodeType(1)
	Json        = CodeType(2)
	ProtoBuffer = CodeType(3)
)

type Coder interface {
	Unmarshal([]byte, interface{}) error   //解码
	Marshal(v interface{}) ([]byte, error) //编码
	ToString() string
}

func GameCoder(codeType CodeType) Coder {
	switch codeType {
	case Json:
		return &jsonCoder{CoderType: Json}
	case String:
		return &rawString{CodeType: String}
	case ProtoBuffer:
		return &protocBufferCoder{CoderType: ProtoBuffer}
	default:
		return &rawString{}
	}
}
