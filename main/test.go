package main

import "github.com/hekonsek/aws-app"

func main() {
	err := (&aws_app.Application{}).Create()
	if err != nil {
		panic(err)
	}
}
