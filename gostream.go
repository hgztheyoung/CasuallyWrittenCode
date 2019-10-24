package play_ground

type StreamIntf interface {
	stream()
}

type StreamList struct {
	SList []StreamIntf
}

func (s StreamList) stream() {}

type StreamSuspension struct {
	Suspension func() StreamIntf
}

func (s StreamSuspension) stream() {}

type StreamValue struct {
	Val interface{}
}

func (s StreamValue) stream() {}
