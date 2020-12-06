package main

import (
	"fmt"
	"promise"
)

func main() {
	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		resolve(2)
	}).Then(func(data interface{}) interface{} {
			fmt.Println("The result is:", data)
			return data.(int) * 2
		}).Then(func(data interface{}) interface{} {
			fmt.Println("Double :", data)
			return nil
		}).Catch(func(error error) error {
			fmt.Println("Error during execution:", error.Error())
			return nil
		})
	p.Await()
}
