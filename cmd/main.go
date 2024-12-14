package main

import (
	"FastGo/internal/bootstrap"
	"flag"
	"log"
)

func main() {
	// 定义命令行标志
	migrate := flag.Bool("migrate", false, "执行数据库迁移")
	flag.Parse()

	// 初始化应用
	app, err := bootstrap.Setup()
	if err != nil {
		log.Fatalf("应用初始化失败: %v", err)
	}

	// 执行数据库迁移
	if *migrate {
		bootstrap.Migrate()
		return
	}

	// 运行应用
	if err := app.Run(); err != nil {
		log.Fatalf("应用运行失败: %v", err)
	}
}
