package obejct

type Result struct {
	IsSuccess bool

	Code int

	Message string

	Data interface{}
}

func OperateFailWithMessage(message string) Result {
	return Result{false, FalseCode, message, nil}
}

func OperateSuccessWithMessage(message string) Result {
	return Result{true, SuccessCode, message, nil}
}

func OperateSuccess() Result {
	return Result{true, SuccessCode, "成功", nil}
}

func OperateSuccess2(object interface{}) Result {
	return Result{true, SuccessCode, "成功", object}

}
