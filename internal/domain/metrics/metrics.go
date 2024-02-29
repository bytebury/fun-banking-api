package metrics

import "time"

type ApplicationInfo struct {
	Name                 string `json:"name"`
	Version              string `json:"version"`
	Message              string `json:"message"`
	NumberOfUsers        int64  `json:"number_of_users"`
	NumberOfBanks        int64  `json:"number_of_banks"`
	NumberOfCustomers    int64  `json:"number_of_customers"`
	NumberOfTransactions int64  `json:"number_of_transactions"`
}

type WeeklyInsights struct {
	Week  int `json:"week"`
	Count int `json:"count"`
}

type VisitorByDay []struct {
	Date          time.Time `json:"date"`
	UserCount     int       `json:"user_count"`
	CustomerCount int       `json:"customer_count"`
	TotalCount    int       `json:"total_count"`
}
