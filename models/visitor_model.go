package models

type Visitor struct {
	AuditModel
	UserID     *uint  `gorm:"user_id"`
	CustomerID *uint  `gorm:"customer_id"`
	IPAddress  string `gotm:"ip_address"`
}
