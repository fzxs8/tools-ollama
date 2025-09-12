package main

import (
	"encoding/json"
	"fmt"
	"log" // Import log package

	"github.com/fzxs8/duolasdk"
)

// OllamaConfigManager Ollama配置管理器
type OllamaConfigManager struct {
	store *duolasdk.AppStore
}

// NewOllamaConfigManager 创建新的Ollama配置管理器
func NewOllamaConfigManager(store *duolasdk.AppStore) *OllamaConfigManager {
	return &OllamaConfigManager{
		store: store,
	}
}

// SaveServers 保存所有服务列表
func (o *OllamaConfigManager) SaveServers(servers []OllamaServerConfig) error {
	key := "ollama_config:servers"
	data, err := json.Marshal(servers)
	if err != nil {
		log.Printf("ERROR: 序列化服务列表失败: %v", err) // Added logging
		return fmt.Errorf("序列化服务列表失败: %w", err)
	}
	return o.store.Set(key, string(data))
}

// GetServers 获取所有服务列表
func (o *OllamaConfigManager) GetServers() ([]OllamaServerConfig, error) {
	key := "ollama_config:servers"
	data, err := o.store.Get(key)
	if err != nil || data == "" {
		log.Printf("INFO: 未找到Ollama服务配置或配置为空，返回默认本地服务. Error: %v, Data empty: %t", err, data == "") // Added logging
		// 如果没有找到配置，返回一个包含默认本地服务的列表
		return nil, nil
	}

	log.Printf("DEBUG: 从存储中获取到Ollama服务配置原始数据: %s", data) // Added logging

	var servers []OllamaServerConfig
	if err := json.Unmarshal([]byte(data), &servers); err != nil {
		log.Printf("ERROR: 反序列化Ollama服务列表失败: %v, 原始数据: %s", err, data) // Added logging
		return nil, fmt.Errorf("反序列化服务列表失败: %w", err)
	}

	return servers, nil
}

// GetServerByID 根据ID获取远程服务器配置
func (o *OllamaConfigManager) GetServerByID(serverID string) (*OllamaServerConfig, error) {
	servers, err := o.GetServers()
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		if server.ID == serverID {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("找不到ID为 %s 的远程服务器", serverID)
}

// AddServer 添加服务器
func (o *OllamaConfigManager) AddServer(server OllamaServerConfig) error {
	servers, err := o.GetServers()
	if err != nil {
		return err
	}

	servers = append(servers, server)
	return o.SaveServers(servers)
}

// UpdateServer 更新服务器
func (o *OllamaConfigManager) UpdateServer(updatedServer OllamaServerConfig) error {
	servers, err := o.GetServers()
	if err != nil {
		return err
	}

	for i, server := range servers {
		if server.ID == updatedServer.ID {
			servers[i] = updatedServer
			return o.SaveServers(servers)
		}
	}

	return fmt.Errorf("未找到要更新的服务器: %s", updatedServer.Name)
}

// DeleteServer 删除服务器
func (o *OllamaConfigManager) DeleteServer(serverID string) error {
	servers, err := o.GetServers()
	if err != nil {
		return err
	}

	newServers := make([]OllamaServerConfig, 0)
	found := false
	for _, server := range servers {
		if server.ID != serverID {
			newServers = append(newServers, server)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("未找到要删除的服务器ID: %s", serverID)
	}

	return o.SaveServers(newServers)
}

// SetActiveServer 设置活动服务器
func (o *OllamaConfigManager) SetActiveServer(serverID string) error {
	servers, err := o.GetServers()
	if err != nil {
		return err
	}

	for i := range servers {
		servers[i].IsActive = (servers[i].ID == serverID)
	}

	return o.SaveServers(servers)
}

// GetActiveServer 获取活动服务器
func (o *OllamaConfigManager) GetActiveServer() (*OllamaServerConfig, error) {
	servers, err := o.GetServers()
	if err != nil {
		return nil, err
	}

	for i, server := range servers {
		if server.IsActive {
			return &servers[i], nil
		}
	}

	// 如果没有活动的，默认第一个为活动
	if len(servers) > 0 {
		return &servers[0], nil
	}

	return nil, fmt.Errorf("没有可用的服务配置")
}

// UpdateServerTestStatus 更新服务的测试状态
func (o *OllamaConfigManager) UpdateServerTestStatus(serverID string, status string) error {
	servers, err := o.GetServers()
	if err != nil {
		return err
	}
	found := false
	for i := range servers {
		if servers[i].ID == serverID {
			servers[i].TestStatus = status
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("未找到ID为 %s 的服务器以更新状态", serverID)
	}
	return o.SaveServers(servers)
}
