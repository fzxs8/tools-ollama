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
func (o *OllamaConfigManager) SaveLocalConfig(baseURL string) error {
	key := "ollama_config:local"
	return o.store.Set(key, baseURL)
}

// GetLocalConfig 获取本地配置
func (o *OllamaConfigManager) GetLocalConfig() (string, error) {
	key := "ollama_config:local"
	return o.store.Get(key)
}

// SaveRemoteServers 保存远程服务器列表
func (o *OllamaConfigManager) SaveRemoteServers(servers []OllamaServerConfig) error {
	key := "ollama_config:remote_servers"
	data, err := json.Marshal(servers)
	if err != nil {
		return fmt.Errorf("failed to marshal remote servers: %w", err)
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
		return []OllamaServerConfig{}, fmt.Errorf("failed to unmarshal remote servers: %w", err)
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

	return nil, fmt.Errorf("no active server found")
}
