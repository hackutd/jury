package src

type Result struct {
	Success bool
	Message string
}

func NewResult(success bool, message string) Result {
	return Result{
		Success: success,
		Message: message,
	}
}

func ResultOk() Result {
	return Result{
		Success: true,
		Message: "",
	}
}
