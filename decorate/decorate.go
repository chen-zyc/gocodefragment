package main

import "fmt"

func main() {
	var b = BasicFunc(func() {
		fmt.Println("basic function")
	})
	var d1 = Decorator(func(f BasicFunc) BasicFunc {
		return BasicFunc(func() {
			fmt.Println("decorator#1 before")
			f()
			fmt.Println("decorator#1 finish")
		})
	})
	var d2 = Decorator(func(f BasicFunc) BasicFunc {
		return BasicFunc(func() {
			fmt.Println("decorator#2 before")
			f()
			fmt.Println("decorator#2 finish")
		})
	})

	f := Decorate(b, d1, d2)
	f()
	// Output:
	// decorator#2 before
	// decorator#1 before
	// basic function
	// decorator#1 finish
	// decorator#2 finish
}

// 基础函数
type BasicFunc func()

// 装饰者，装饰传入的BasicFunc，返回新的BasicFunc
type Decorator func(BasicFunc) BasicFunc

// 装饰，用ds[0]装饰f, 然后用ds[1]装饰ds[0]...
// 这是在摆多米若骨牌，返回的函数执行时，就是骨牌被推倒时。
func Decorate(f BasicFunc, ds ...Decorator) BasicFunc {
	decorated := f
	for _, d := range ds {
		decorated = d(decorated)
	}
	return decorated
}
