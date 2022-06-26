// main.go

package main

import (
	"TodoApp/app/controllers"
	"TodoApp/app/models"
	"fmt"
)

func main() {
	fmt.Println(models.Db)

	controllers.StartMainServer()
}
