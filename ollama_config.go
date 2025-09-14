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

// NewOllamaConfigManager 创建新的Ollama配置管理器, 并在启动时执行数据迁移
func NewOllamaConfigManager(store *duolasdk.AppStore, logger *core.AppLog) *OllamaConfigManager {
	o := &OllamaConfigManager{
		store:  store,
		logger: logger.WithPrefix("ConfigManager"),
	}
	// 应用启动时执行一次数据迁移
	o.migrateServersToHashFormat()
	return o
}

// migrateServersToHashFormat 将旧的 string 格式服务器列表迁移到新的 Hash 格式
// 这是一个内部方法, 只在初始化时调用一次
func (o *OllamaConfigManager) migrateServersToHashFormat() {
	// 首先检查新格式是否存在, 如果存在, 则无需迁移
	dataMap, err := o.store.HGetAll(serversKey)
	if err == nil && len(dataMap) > 0 {
		o.logger.Debug("已是新版Hash存储格式，无需迁移。")
		return
	}

	// 尝试用旧方法(Get)读取，以进行数据迁移
	o.logger.Info("未找到Hash格式的配置，尝试从旧格式(string)迁移...")
	oldData, oldErr := o.store.Get(serversKey)
	if oldErr != nil || oldData == "" {
		o.logger.Info("未发现旧格式配置，无需迁移。")
		return
	}

	o.logger.Info("发现旧格式的服务器配置，开始迁移...")
	var oldServers []types.OllamaServerConfig
	if err := json.Unmarshal([]byte(oldData), &oldServers); err != nil {
		o.logger.Error("反序列化旧格式配置失败，迁移中断", "error", err)
		return
	}

	// 将数据以新的Hash格式写回
	for _, server := range oldServers {
		// 使用 AddServer 内部的逻辑来写入
		if err := o.AddServer(server); err != nil {
			o.logger.Error("迁移服务器时写入新格式失败", "serverID", server.ID, "error", err)
		}
	}

	// 成功迁移后, 删除旧的string类型的key
	if err := o.store.Delete(serversKey); err != nil {
		o.logger.Error("删除旧配置键失败", "error", err)
	}

	o.logger.Info("旧数据迁移到新Hash格式成功")
}

// GetServers 获取所有服务列表 (纯读取操作)
func (o *OllamaConfigManager) GetServers() ([]types.OllamaServerConfig, error) {
	o.logger.Debug("开始获取所有服务器配置")

	dataMap, err := o.store.HGetAll(serversKey)
	if err != nil {
		o.logger.Error("从存储中获取服务器列表失败", "error", err)
		// 即使出错也返回空列表, 避免前端出错
		return []types.OllamaServerConfig{}, nil
	}

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
		// 使用索引来修改切片中的元素
		if servers[i].ID == serverID {
			servers[i].IsActive = true
			found = true
		} else {
			servers[i].IsActive = false
		}
	}

	if !found {
		return fmt.Errorf("未找到要设置为活动的服务器ID: %s", serverID)
	}

	// 批量更新
	for _, server := range servers {
		if err := o.UpdateServer(server); err != nil {
			o.logger.Error("更新服务器活动状态时失败", "serverID", server.ID, "error", err)
			// 即使部分失败也继续
		}
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

	// 注意：这里的自动设置默认活动服务器的副作用逻辑已被移除，以防止意外的写入操作。
	// 如果没有活动的服务器，调用者应该负责处理这种情况。
	if len(servers) > 0 {
		o.logger.Warn("未找到明确的活动服务器，将默认返回列表中的第一个，但不会将其持久化为活动状态。")
		return &servers[0], nil
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
