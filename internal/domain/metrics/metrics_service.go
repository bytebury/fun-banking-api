package metrics

type MetricService interface {
	GetApplicationInfo() (ApplicationInfo, error)
	GetUsersInfo() ([]WeeklyInsights, error)
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

func (s metricService) GetUsersInfo() ([]WeeklyInsights, error) {
	return s.metricRepository.GetUsersInfo()
}
