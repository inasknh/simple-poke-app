package api

type BerriesRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type Berry struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type BerriesResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous string  `json:"previous"`
	Results  []Berry `json:"results"`
}
