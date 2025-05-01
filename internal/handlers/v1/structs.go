package v1

// Structure of base answer
type JSONResult struct {
	Code    int         `json:"code" description:"Code answer" example:"200"`
	Message string      `json:"message,omitempty" description:"Message text with error" example:"test error"`
	Data    interface{} `json:"data,omitempty" description:"Block for user data"`
}
