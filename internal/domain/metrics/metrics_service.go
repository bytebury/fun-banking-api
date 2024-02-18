package metrics

type MetricService interface {
	GetApplicationInfo() (ApplicationInfo, error)
}

type metricService struct {
	metricRepository MetricRepository
}

func NewMetricService() MetricService {
	return metricService{
		metricRepository: NewMetricRepository(),
	}
}

func (s metricService) GetApplicationInfo() (ApplicationInfo, error) {
	var appInfo ApplicationInfo
	err := s.metricRepository.GetApplicationInfo(&appInfo)
	return appInfo, err
}
