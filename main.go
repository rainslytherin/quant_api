package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quant_api/api"
	"quant_api/config"
	"quant_api/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "config.json", "config file")
}

func main() {
	flag.Parse()

	// 0.初始化全局配置
	config.InitGlobalConfig(configFile)
	defaultLogger := logger.New(config.GetGlobalConfig().Log.Level)

	// 6.创建HTTP API Service
	apiConfig := api.NewConfig(config.GetGlobalConfig())
	apiService := api.CreateService(apiConfig)
	apiService.WithLogger(defaultLogger)
	apiService.Start()

	// 创建一个通道用于接收信号
	sigChan := make(chan os.Signal, 1)

	// 将SIGINT信号发送到通道
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	// 启动一个goroutine等待信号
	go func() {
		<-sigChan

		fmt.Println("Received SIGINT signal, start to exit...")

		// 执行你想要在收到信号时执行的操作
		// 例如：关闭服务器，释放资源等
		apiService.Close()

		fmt.Println("Sleep 3 seconds to wait for all goroutines to exit...")
		time.Sleep(3 * time.Second)
		fmt.Println("Server exited.")

		// 退出程序
		os.Exit(0)
	}()

	// 在这里可以执行你的主进程的其他操作

	// 通过Sleep来模拟主进程的持续运行
	select {}
}
