package metrics

type MetricService interface {
	GetApplicationInfo() (ApplicationInfo, error)
	GetUsersInfo() ([]WeeklyInsights, error)
	GetVisitorsInfo() (VisitorByDay, error)
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

func (s metricService) GetVisitorsInfo() (VisitorByDay, error) {
	return s.metricRepository.GetVisitorsByDay()
}
