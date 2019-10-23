type StreamIntf interface {
	stream()
}

type StreamList struct {
	SList []StreamIntf
}

type StreamSuspension struct {
	Suspension func() StreamIntf
}

func (s StreamSuspension) stream() {}

type StreamValue struct {
	Val interface{}
}

func (s StreamValue) stream() {
}
