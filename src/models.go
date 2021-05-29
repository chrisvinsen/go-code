package models

type RawCurrency struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     struct {
		USD float64 `json:"USD"`
		CAD float64 `json:"CAD"`
		IDR float64 `json:"IDR"`
		GBP float64 `json:"GBP"`
		CHF float64 `json:"CHF"`
		SGD float64 `json:"SGD"`
		INR float64 `json:"INR"`
		MYR float64 `json:"MYR"`
		JPY float64 `json:"JPY"`
		KRW float64 `json:"KRW"`
	} `json:"rates"`
}

type ResponseCurrency struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     []Rate `json:"rates"`
}

type Rate struct {
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}
