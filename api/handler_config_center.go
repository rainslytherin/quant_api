package api

import (
	"time"

	"quant_api/models"

	"github.com/gin-gonic/gin"
)

type SyncParams struct {
	ClientID   string `json:"client_id"`
	UpdateTime int64  `json:"update_time"`
}

// ConfigCenterGetConfigs godoc
func (s *Service) ConfigCenterGetConfigs(c *gin.Context) {
	// get client_id from query params
	var param SyncParams
	if err := c.ShouldBindQuery(&param); err != nil {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}

	syncStatus, err := models.LoadSyncStatus(param.ClientID)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "获取同步状态失败: "+err.Error())
		return
	}

	if syncStatus == nil {
		SetHTTPResponse(c, -1, nil, "未找到同步状态")
		return
	}

	var updateTime int64

	if param.UpdateTime != 0 {
		updateTime = param.UpdateTime
	} else if syncStatus != nil {
		updateTime = syncStatus.UpdateTime
	}

	configs, err := models.GetConfigsAfterTime(updateTime)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "获取配置失败: "+err.Error())
		return
	}

	// 返回所有配置
	data := make(map[string]interface{})
	data["configs"] = configs

	SetHTTPResponse(c, 0, data, "查询成功")

	syncStatus.UpdateTime = time.Now().Unix()
	err = syncStatus.Save()
	if err != nil {
		s.Logger.Error("save sync status failed", "error", err)
	}
}

// ConfigCenterGetStatus godoc
func (s *Service) ConfigCenterGetStatus(c *gin.Context) {
	syncStatus, err := models.LoadAllSyncStatus()
	if err != nil {
		SetHTTPResponse(c, -1, nil, "获取同步状态失败: "+err.Error())
		return
	}

	data := make(map[string]interface{})
	data["sync_status"] = syncStatus

	SetHTTPResponse(c, 0, data, "查询成功")
}
