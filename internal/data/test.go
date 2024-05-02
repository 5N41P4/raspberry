package data

// Test Data struct for test purposes with the charts.
type Data struct {
	Labels   []string  `json:"labels"`   // Labels for the chart
	Datasets []Dataset `json:"datasets"` // Datasets for the chart
}

type Dataset struct {
	Label           string   `json:"label"`           // Label for the dataset
	Data            []int    `json:"data"`            // Data for the dataset
	BackgroundColor []string `json:"backgroundColor"` // Background color for the dataset
	BorderColor     string   `json:"borderColor"`     // Border color for the dataset
	BorderWidth     int      `json:"borderWidth"`     // Border width for the dataset
}
