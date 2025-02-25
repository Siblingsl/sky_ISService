package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/fx"
)

func main() {
	// 使用 WaitGroup 来等待所有服务的启动和关闭
	var wg sync.WaitGroup

	// 创建一个新的 fx 应用
	app := fx.New(
		// 提供 WaitGroup 到容器中，供服务启动使用
		fx.Provide(func() *sync.WaitGroup { return &wg }),

		// 注册服务启动和关闭逻辑
		fx.Invoke(func(lc fx.Lifecycle, wg *sync.WaitGroup) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					// 启动服务
					startServiceWithWaitGroup("./services/auth/cmd/main.go", wg)
					startServiceWithWaitGroup("./services/system/cmd/main.go", wg)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					// 这里可以加更多的清理操作
					log.Println("所有服务即将停止")
					return nil
				},
			})
		}),
	)

	// 启动应用
	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("启动应用失败: %v", err)
	}

	// 等待服务完成并优雅关闭
	waitForShutdown()

	// 等待所有服务完成
	wg.Wait()

	// 停止应用
	if err := app.Stop(context.Background()); err != nil {
		log.Fatalf("停止应用失败: %v", err)
	}
}

// 启动服务并等待完成
func startServiceWithWaitGroup(servicePath string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startService(servicePath); err != nil {
			log.Printf("启动服务 [%s] 失败: %v", servicePath, err)
		} else {
			log.Printf("服务 [%s] 启动成功", servicePath)
		}
	}()
}

// 启动服务的通用函数
func startService(servicePath string) error {
	cmd := exec.Command("go", "run", servicePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	// 等待服务完成执行
	return cmd.Wait()
}

// 等待系统信号并优雅关闭服务
func waitForShutdown() {
	// 创建一个 channel 来接收系统的中断信号
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 等待信号
	<-signalChannel

	// 超时机制：设置超时时间避免永远等待
	shutdownTimeout := time.After(10 * time.Second)
	select {
	case <-shutdownTimeout:
		log.Println("服务关闭超时，强制退出")
	case <-signalChannel:
		log.Println("收到关闭信号，优雅退出")
	}

	// 这里可以增加更多的清理操作，比如关闭数据库连接等
	log.Println("所有子服务已优雅退出")
}
