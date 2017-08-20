package main

import (
    "net/http"

    "github.com/labstack/echo"
	"time"
	"fmt"
)

func main() {
    e := echo.New()
    db := ConnectDB()
	fmt.Println("qweqwe")
	e.GET("/", func(c echo.Context) error {

		db.Create(&Product{Code: "200", Time: time.Now()})
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1325"))
	fmt.Println("badbad")
	defer db.Close()
}
