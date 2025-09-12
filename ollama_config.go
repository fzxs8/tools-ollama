package main

import (
	"encoding/json"
	"fmt"
	"tools-ollama/types"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
)

// OllamaConfigManager Ollama配置管理器
type OllamaConfigManager struct {
	store  *duolasdk.AppStore
	logger *core.AppLog
}

const serversKey = "ollama_config:servers"

// NewOllamaConfigManager 创建新的Ollama配置管理器
func NewOllamaConfigManager(store *duolasdk.AppStore, logger *core.AppLog) *OllamaConfigManager {
	return &OllamaConfigManager{
		store:  store,
		logger: logger.WithPrefix("ConfigManager"),
	}
}

// GetServers 获取所有服务列表 (包含数据迁移逻辑)
func (o *OllamaConfigManager) GetServers() ([]types.OllamaServerConfig, error) {
	o.logger.Debug("开始获取所有服务器配置")

	// 优先尝试用新方法(HGetAll)读取
	dataMap, err := o.store.HGetAll(serversKey)
	if err == nil && len(dataMap) > 0 {
		o.logger.Debug("检测到Hash存储格式的服务器配置，直接处理")
		servers := make([]types.OllamaServerConfig, 0, len(dataMap))
		for id, data := range dataMap {
			var server types.OllamaServerConfig
			if err := json.Unmarshal([]byte(data), &server); err != nil {
				o.logger.Warn("反序列化单个服务器配置失败", "id", id, "error", err)
				continue
			}
			servers = append(servers, server)
		}
		return servers, nil
	}

	// 尝试用旧方法(Get)读取，以进行数据迁移
	o.logger.Info("未找到Hash格式的配置，尝试从旧格式(string)迁移...")
	oldData, oldErr := o.store.Get(serversKey)
	if oldErr == nil && oldData != "" {
		o.logger.Info("发现旧格式的服务器配置，开始迁移...")
		var oldServers []types.OllamaServerConfig
		if err := json.Unmarshal([]byte(oldData), &oldServers); err != nil {
			o.logger.Error("反序列化旧格式配置失败，迁移中断", "error", err)
			return []types.OllamaServerConfig{}, nil
		}

		// 删除旧的string类型的key
		if err := o.store.Delete(serversKey); err != nil {
			o.logger.Error("删除旧配置键失败", "error", err)
		}

		// 将数据以新的Hash格式写回
		for _, server := range oldServers {
			if err := o.AddServer(server); err != nil {
				o.logger.Error("迁移服务器时写入新格式失败", "serverID", server.ID, "error", err)
			}
		}
		o.logger.Info("旧数据迁移到新Hash格式成功")
		return oldServers, nil
	}

	// 如果新旧两种方式都没有数据，说明就是没有配置，返回空列表
	o.logger.Info("未找到任何服务器配置，返回空列表")
	return []types.OllamaServerConfig{}, nil
}

// GetServerByID 根据ID获取服务器配置
func (o *OllamaConfigManager) GetServerByID(serverID string) (*types.OllamaServerConfig, error) {
	o.logger.Debug("根据ID获取服务器配置", "serverID", serverID)
	data, err := o.store.HGet(serversKey, serverID)
	if err != nil {
		o.logger.Error("从存储中获取服务器失败", "serverID", serverID, "error", err)
		return nil, fmt.Errorf("找不到ID为 %s 的服务器: %w", serverID, err)
	}

	var server types.OllamaServerConfig
	if err := json.Unmarshal([]byte(data), &server); err != nil {
		o.logger.Error("反序列化服务器配置失败", "serverID", serverID, "error", err)
		return nil, fmt.Errorf("解析服务器 %s 的配置失败: %w", serverID, err)
	}

	return &server, nil
}

// AddServer 添加或更新服务器
func (o *OllamaConfigManager) AddServer(server types.OllamaServerConfig) error {
	o.logger.Info("添加或更新服务器", "serverName", server.Name, "serverID", server.ID)
	data, err := json.Marshal(server)
	if err != nil {
		o.logger.Error("序列化服务器配置失败", "serverName", server.Name, "error", err)
		return fmt.Errorf("序列化服务配置失败: %w", err)
	}

	if err := o.store.HSet(serversKey, server.ID, string(data)); err != nil {
		o.logger.Error("向存储中写入服务器配置失败", "serverName", server.Name, "error", err)
		return fmt.Errorf("保存服务配置失败: %w", err)
	}
	return nil
}

// UpdateServer 更新服务器
func (o *OllamaConfigManager) UpdateServer(updatedServer types.OllamaServerConfig) error {
	return o.AddServer(updatedServer)
}

// DeleteServer 删除服务器
func (o *OllamaConfigManager) DeleteServer(serverID string) error {
	o.logger.Info("删除服务器", "serverID", serverID)
	if err := o.store.HDel(serversKey, serverID); err != nil {
		o.logger.Error("从存储中删除服务器失败", "serverID", serverID, "error", err)
		return fmt.Errorf("删除服务器 %s 失败: %w", serverID, err)
	}
	return nil
}

// SetActiveServer 设置活动服务器
func (o *OllamaConfigManager) SetActiveServer(serverID string) error {
	o.logger.Info("设置活动服务器", "serverID", serverID)
	servers, err := o.GetServers()
	if err != nil {
		return err
	}

	found := false
	for i := range servers {
		server := servers[i]
		if server.ID == serverID {
			server.IsActive = true
			found = true
		} else {
			server.IsActive = false
		}
		if err := o.UpdateServer(server); err != nil {
			o.logger.Error("更新服务器活动状态时失败", "serverID", server.ID, "error", err)
		}
	}

	if !found {
		return fmt.Errorf("未找到要设置为活动的服务器ID: %s", serverID)
	}

	return nil
}

// GetActiveServer 获取活动服务器
func (o *OllamaConfigManager) GetActiveServer() (*types.OllamaServerConfig, error) {
	o.logger.Debug("获取活动服务器")
	servers, err := o.GetServers()
	if err != nil {
		return nil, err
	}

	for i, server := range servers {
		if server.IsActive {
			o.logger.Debug("找到活动服务器", "serverName", server.Name)
			return &servers[i], nil
		}
	}

	if len(servers) > 0 {
		o.logger.Warn("未找到明确的活动服务器，将默认使用列表中的第一个")
		firstServer := servers[0]
		firstServer.IsActive = true
		if err := o.UpdateServer(firstServer); err != nil {
			o.logger.Error("设置默认活动服务器失败", "error", err)
		}
		return &firstServer, nil
	}

	return nil, fmt.Errorf("没有可用的服务配置")
}

// UpdateServerTestStatus 更新服务的测试状态
func (o *OllamaConfigManager) UpdateServerTestStatus(serverID string, status string) error {
	o.logger.Debug("更新服务器测试状态", "serverID", serverID, "status", status)
	server, err := o.GetServerByID(serverID)
	if err != nil {
		return err
	}

	server.TestStatus = status
	return o.UpdateServer(*server)
}
