package ds

type Stream struct {
	Id string
}

func NewStream(id string) *Stream {
	return &Stream{Id: id}
}

type Streams struct {
	streams map[string]*Stream
}

func NewStreams() *Streams {
	return &Streams{
		streams: make(map[string]*Stream),
	}
}

func (s *Streams) Register(key string) {
	stream := NewStream(key)
	s.streams[key] = stream
}

func (s *Streams) Contains(key string) bool {
	_, ok := s.streams[key]
	return ok
}
