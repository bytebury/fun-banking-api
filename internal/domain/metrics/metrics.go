package metrics

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
