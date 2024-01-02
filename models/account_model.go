package models

import "time"

type Account struct {
	AuditModel
	Name       string   `json:"name" gorm:"not null;size:100"`
	Balance    float64  `json:"balance" gorm:"type:decimal(50,2);not null;default:0.00"`
	CustomerID uint     `json:"customer_id" gorm:"not null"`
	Customer   Customer `json:"customer" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;"`
}

type DailyTransferSummary struct {
	Date         time.Time `json:"date"`
	TotalBalance float64   `json:"total_balance"`
}
