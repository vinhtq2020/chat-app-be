package usecase

type LogAggregator struct{}

func NewLogAggregator() *LogAggregator {
	return &LogAggregator{}
}

func (*LogAggregator) RegisterJob() {

}
