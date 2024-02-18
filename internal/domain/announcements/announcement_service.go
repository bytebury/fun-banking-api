package announcements

type AnnouncementService interface {
	FindByID(id string) (Announcement, error)
	Create(announcement *Announcement) error
	Update(id string, announcement *Announcement) error
}

type announcementService struct {
	announcementRepository AnnouncementRepository
}

func NewAnnouncementService(announcementRepository AnnouncementRepository) AnnouncementService {
	return announcementService{announcementRepository}
}

func (s announcementService) FindByID(id string) (Announcement, error) {
	var announcement Announcement
	err := s.announcementRepository.FindByID(id, &announcement)
	return announcement, err
}

func (s announcementService) Create(announcement *Announcement) error {
	return s.announcementRepository.Create(announcement)
}

func (s announcementService) Update(id string, announcement *Announcement) error {
	return s.announcementRepository.Update(id, announcement)
}
