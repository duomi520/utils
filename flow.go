package utils

//Three 三元表达式
func Three[T any](boolExpression bool, trueReturnValue, falseReturnValue T) T {
	if boolExpression {
		return trueReturnValue
	} else {
		return falseReturnValue
	}
}

// https://github.com/reugn/go-streams
// https://zhuanlan.zhihu.com/p/452984498
