package models

type Health struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	Message           string `json:"message"`
	NumberOfUsers     int64  `json:"number_of_users"`
	NumberOfBanks     int64  `json:"number_of_banks"`
	NumberOfCustomers int64  `json:"number_of_customers"`
	NumberOfTransfers int64  `json:"number_of_transfers"`
}