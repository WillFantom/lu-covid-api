package rates

type ResponseContent struct {
	Title    string `json:"title"`
	Abstract string `json:"abstract`
	Content  string `json:"main"`
}

type WebResponse struct {
	Key     string            `json:"key"`
	Title   string            `json:"title"`
	Content []ResponseContent `json:"contentItems"`
}
