# gocodefragment
go代码片段

## coerce

摘抄自 <https://github.com/mreiferson/go-options/blob/master/options.go>

函数声明：`func Coerce(v interface{}, opt interface{}, arg string) (interface{}, error)`

功能：根据opt的类型转换v，例如如果opt的类型是int，那么就将v转换为int类型。当v时数值字符串(如"123")且opt是time.Duration类型时，arg被当做单位，如"1s","2m"，其他时候会被忽略。

