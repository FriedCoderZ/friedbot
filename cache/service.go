package cache

type Service struct {
	GroupEventsPool   *Pool
	PrivateEventsPool *Pool
}

func NewService() *Service {
	return &Service{}
}

const (
	smallCacheSize  = 16
	normalCacheSize = 128
	largeCacheSize  = 512
	hugeCacheSize   = 2048
)

func (s *Service) Start() {
	s.GroupEventsPool = NewPool(normalCacheSize)
	s.PrivateEventsPool = NewPool(normalCacheSize)
}
