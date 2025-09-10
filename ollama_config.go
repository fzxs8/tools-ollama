package main

import (
	"encoding/json"
	"fmt"

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

// SaveLocalConfig 保存本地配置
func (o *OllamaConfigManager) SaveLocalConfig(config OllamaServerConfig) error {
	key := "ollama_config:local"
	config.ID = "local"
	config.Type = "local"
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化本地服务配置失败: %w", err)
	}
	return o.store.Set(key, string(data))
}

// GetLocalConfig 获取本地配置
func (o *OllamaConfigManager) GetLocalConfig() (OllamaServerConfig, error) {
	key := "ollama_config:local"
	data, err := o.store.Get(key)
	if err != nil {
		// 如果没有找到配置，返回默认配置
		return OllamaServerConfig{
			ID:         "local",
			Name:       "本地服务",
			BaseURL:    "http://localhost:11434",
			Type:       "local",
			TestStatus: "unknown",
		}, nil
	}

	var config OllamaServerConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return OllamaServerConfig{}, fmt.Errorf("反序列化本地服务配置失败: %w", err)
	}

	return config, nil
}

// GetLocalConfigPtr 获取本地配置的指针
func (o *OllamaConfigManager) GetLocalConfigPtr() (*OllamaServerConfig, error) {
	config, err := o.GetLocalConfig()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetRemoteServerByID 根据ID获取远程服务器配置
func (o *OllamaConfigManager) GetRemoteServerByID(serverID string) (*OllamaServerConfig, error) {
	servers, err := o.GetRemoteServers()
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

// SaveRemoteServers 保存远程服务器列表
func (o *OllamaConfigManager) SaveRemoteServers(servers []OllamaServerConfig) error {
	key := "ollama_config:remote_servers"
	data, err := json.Marshal(servers)
	if err != nil {
		return fmt.Errorf("序列化远程服务器列表失败: %w", err)
	}
	return o.store.Set(key, string(data))
}

// GetRemoteServers 获取远程服务器列表
func (o *OllamaConfigManager) GetRemoteServers() ([]OllamaServerConfig, error) {
	key := "ollama_config:remote_servers"
	data, err := o.store.Get(key)
	if err != nil {
		// 如果没有找到配置，返回空列表
		return []OllamaServerConfig{}, nil
	}

	var servers []OllamaServerConfig
	if err := json.Unmarshal([]byte(data), &servers); err != nil {
		return []OllamaServerConfig{}, fmt.Errorf("反序列化远程服务器列表失败: %w", err)
	}

	return servers, nil
}

// AddRemoteServer 添加远程服务器
func (o *OllamaConfigManager) AddRemoteServer(server OllamaServerConfig) error {
	servers, err := o.GetRemoteServers()
	if err != nil {
		return err
	}

	servers = append(servers, server)
	return o.SaveRemoteServers(servers)
}

// UpdateRemoteServer 更新远程服务器
func (o *OllamaConfigManager) UpdateRemoteServer(updatedServer OllamaServerConfig) error {
	servers, err := o.GetRemoteServers()
	if err != nil {
		return err
	}

	for i, server := range servers {
		if server.ID == updatedServer.ID {
			servers[i] = updatedServer
			break
		}
	}

	return o.SaveRemoteServers(servers)
}

// DeleteRemoteServer 删除远程服务器
func (o *OllamaConfigManager) DeleteRemoteServer(serverID string) error {
	servers, err := o.GetRemoteServers()
	if err != nil {
		return err
	}

	for i, server := range servers {
		if server.ID == serverID {
			servers = append(servers[:i], servers[i+1:]...)
			break
		}
	}

	return o.SaveRemoteServers(servers)
}

// SetActiveServer 设置活动服务器
func (o *OllamaConfigManager) SetActiveServer(serverID string) error {
	servers, err := o.GetRemoteServers()
	if err != nil {
		return err
	}

	// 将所有服务器设置为非活动状态，然后将指定服务器设置为活动状态
	for i, server := range servers {
		servers[i].IsActive = (server.ID == serverID)
	}

	return o.SaveRemoteServers(servers)
}

// GetActiveServer 获取活动服务器
func (o *OllamaConfigManager) GetActiveServer() (*OllamaServerConfig, error) {
	servers, err := o.GetRemoteServers()
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		if server.IsActive {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("未找到活动服务器")
}

// SaveLocalServerTestStatus 保存本地服务器的测试状态
func (o *OllamaConfigManager) SaveLocalServerTestStatus(status string) error {
	key := "ollama_config:local:test_status"
	return o.store.Set(key, status)
}

// GetLocalServerTestStatus 获取本地服务器的测试状态
func (o *OllamaConfigManager) GetLocalServerTestStatus() (string, error) {
	key := "ollama_config:local:test_status"
	status, err := o.store.Get(key)
	if err != nil {
		return "unknown", nil // 如果没有找到，默认为 "unknown"
	}
	return status, nil
}
