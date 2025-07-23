package validator

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors,omitempty"`
}

func (val ValidationErrors) Error() string {
	return "validation error"
}
