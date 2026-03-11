package v1

// JSONResult ...
type JSONResult struct {
	Code    int    `json:"code" description:"Code answer" example:"200"`
	Message string `json:"message,omitempty" description:"Message text with error" example:"test error"`
	Data    any    `json:"data,omitempty" description:"Block for user data"`
}
