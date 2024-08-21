package main

import (
	"flag"
	"fmt"
	"service/server"
	"service/utils"
)

func main() {
	// load .env
	success := utils.InitEnv()
	if !success {
		return
	}

	// init redis
	success = utils.InitRedis()
	if !success {
		fmt.Println("InitRedis: fail")
		return
	} else {
		fmt.Println("InitRedis: success")
	}

	// init DB
	success = utils.InitDB()
	if !success {
		fmt.Println("InitDB: fail")
		return
	} else {
		fmt.Println("InitDB: success")
	}

	// start server
	flag.Parse()
	if success {
		switch flag.Arg(0) {
		case "gin":
			server.InitGinServer()
		case "grpc":
			server.InitGinServer()
		default:
			panic("choose a service")
		}
	}
}
