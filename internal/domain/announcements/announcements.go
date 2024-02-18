package announcements

import "funbanking/internal/infrastructure/persistence"

func RunMigrations() {
	persistence.DB.AutoMigrate(&Announcement{})
}
