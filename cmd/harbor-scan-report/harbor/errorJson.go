package harbor

type ErrorJson struct {
	Errors []struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"errors"`
}
