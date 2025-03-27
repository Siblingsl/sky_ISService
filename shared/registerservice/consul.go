package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"sky_ISService/config"
)

// InitConsul 初始化 Consul 客户端
func InitConsul() (*api.Client, error) {
	configConsul, err := config.InitLoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %v", err)
	}

	// 获取 Consul 配置
	consulConfig := configConsul.RegisterService.Consul

	// 创建 Consul 配置对象
	configData := &api.Config{
		Address: consulConfig.Address + ":" + consulConfig.Port, // Consul 地址
	}

	// 创建 Consul 客户端
	client, err := api.NewClient(configData)
	if err != nil {
		return nil, fmt.Errorf("无法创建 Consul 客户端: %v", err)
	}

	fmt.Println("成功连接到 Consul")

	return client, nil
}

// 注册服务到 Consul
func RegisterServiceConsul(client *api.Client, serviceName, serviceID, address string, port int) error {
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
	}

	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("服务注册失败: %v", err)
	}

	log.Printf("服务 %s 注册成功", serviceName)
	return nil
}

// 查询服务信息
//func GetServiceConsul(client *api.Client, serviceName string) ([]*api.ServiceEntry, error) {
//	services, _, err := client.Catalog().Service(serviceName, "", nil)
//	if err != nil {
//		return nil, fmt.Errorf("查询服务 %s 失败: %v", serviceName, err)
//	}
//
//	return services, nil
//}
