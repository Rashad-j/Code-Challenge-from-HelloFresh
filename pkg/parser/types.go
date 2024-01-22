package parser

// Entry is a struct that is used to stream data over the channel
type Entry struct {
	Recipe Recipe `json:"recipe"`
	Error  error  `json:"error"`
}

type Recipe struct {
	Recipe   string `json:"recipe"`
	Postcode string `json:"postcode"`
	Delivery string `json:"delivery"`
}
