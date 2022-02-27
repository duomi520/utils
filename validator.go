package utils

//IValidator 接口
type IValidator interface {
	Var(any, string) error
	Struct(any) error
}
