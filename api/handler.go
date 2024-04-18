package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"quant_api/models"

	"github.com/gin-gonic/gin"
)

func (s *Service) InitHandlers() {
	s.GET("/hello", s.hello)

	// add close_out handler
	s.POST("/stock/close_out", s.closeOut)

	// add config update handler
	s.GET("/stock/configs", s.GetStockConfigs)
	s.POST("/stock/configs", s.AddStockConfig)
	s.PUT("/stock/configs", s.UpdateStockConfig)
	s.DELETE("/stock/configs", s.DeleteStockConfig)

	s.GET("/global/configs", s.GetGlobalConfigs)
	s.POST("/global/configs", s.AddGlobalConfig)
	s.PUT("/global/configs", s.UpdateGlobalConfig)
}

func (s *Service) hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello",
	})
	return
}

type CloseOut struct {
	StockCode string `json:"stock_code"  binding:"required"`
}

func SetHTTPResponse(c *gin.Context, code int, data interface{}, message string) {
	if data == nil {
		data = make(map[string]interface{})
	}
	c.JSON(200, gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	})
}

// 反序列化结果
type CloseOutResponse struct {
	Data    CloseOutInfo `json:"data"`
	Message string       `json:"message"`
}

type CloseOutInfo struct {
	LastPrice          float64     `json:"price"`
	CanceledEnterTasks []int       `json:"canceled_enter_tasks"`
	CanceledExitTasks  []int       `json:"canceled_exit_tasks"`
	CloseOutQty        map[int]int `json:"close_out_qty"`
}

func (s *Service) closeOut(c *gin.Context) {
	var closeOut CloseOut
	if err := c.ShouldBindJSON(&closeOut); err != nil {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}

	s.Logger.Info("closeOut", "closeOut", closeOut)

	if closeOut.StockCode == "" {
		SetHTTPResponse(c, -1, nil, "stock_code 不能为空")
		return
	}

	// do something
	quantCoreBackend := s.cfg.GetBackend("quant_core")

	// build http request to quant_core_backend
	// use http request post method
	// with json body

	// 设置请求地址和方法
	url := fmt.Sprintf("http://%s/stock/close_out", quantCoreBackend)
	method := "POST"

	// 设置请求体
	jsonData, err := json.Marshal(closeOut)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "序列化请求体失败")
		return
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonData))
	if err != nil {
		SetHTTPResponse(c, -1, nil, "创建请求失败")
		return
	}

	// 设置请求头部信息
	req.Header.Set("Content-Type", "application/json")

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("请求失败: %s", err.Error())
		SetHTTPResponse(c, -1, nil, msg)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "解析响应失败:"+err.Error())
		return
	}

	var closeOutResponse CloseOutResponse

	if err := json.Unmarshal(body, &closeOutResponse); err != nil {
		SetHTTPResponse(c, -1, nil, "反序列化失败:"+err.Error())
		return
	}

	if resp.StatusCode != 200 {
		SetHTTPResponse(c, -1, nil, "计算模块处理异常， "+closeOutResponse.Message)
		return

	}

	SetHTTPResponse(c, 0, closeOutResponse, "执行完成")
	return
}

// GetStockConfigs
func (s *Service) GetStockConfigs(c *gin.Context) {
	scope := "stock"
	configs, err := models.GetConfigs(scope)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "获取配置失败: "+err.Error())
		return
	}

	// 返回所有股票配置
	data := make(map[string]interface{})
	data["configs"] = configs
	SetHTTPResponse(c, 0, data, "查询成功")
}

type StockConfig struct {
	StockCode string `json:"stock_code"  binding:"required"`
	Config    struct {
		ProdStatus *bool    `json:"prod_status,omitempty"`
		PreStatus  *bool    `json:"pre_status,omitempty"`
		UpLimit    *float64 `json:"up_limit,omitempty"`
		LowLimit   *float64 `json:"low_limit,omitempty"`
	} `json:"config"`
	UpdateUser string `json:"update_user"  binding:"required"`
}

// AddStockConfig
func (s *Service) AddStockConfig(c *gin.Context) {
	scope := "stock"
	// 添加股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}

	s.Logger.Info("AddStockConfig", "stockConfig", stockConfig)

	if stockConfig.StockCode == "" {
		SetHTTPResponse(c, -1, nil, "stock code 不能为空")
		return
	}

	value, err := json.Marshal(stockConfig.Config)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "参数格式错误")
		return
	}

	var valueObject map[string]interface{}

	if err := json.Unmarshal(value, &valueObject); err != nil {
		SetHTTPResponse(c, -1, nil, "参数格式错误")
		return
	}

	scopeConfig := models.NewConfig(scope, stockConfig.StockCode, valueObject, stockConfig.UpdateUser)
	if err := scopeConfig.Create(); err != nil {
		SetHTTPResponse(c, -1, nil, "添加失败: "+err.Error())
		return
	}

	data := make(map[string]interface{})
	data["config"] = scopeConfig
	SetHTTPResponse(c, 0, data, "添加成功")
}

// UpdateStockConfig
func (s *Service) UpdateStockConfig(c *gin.Context) {
	scope := "stock"
	// 更新股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}

	s.Logger.Info("UpdateStockConfig", "stockConfig", stockConfig)

	value, err := json.Marshal(stockConfig.Config)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "参数格式错误")
		return
	}

	var valueObject map[string]interface{}

	if err := json.Unmarshal(value, &valueObject); err != nil {
		SetHTTPResponse(c, -1, nil, "参数格式错误")
		return
	}

	if len(valueObject) == 0 {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}

	oldConfig, err := models.GetConfig(scope, stockConfig.StockCode)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "获取原配置失败: "+err.Error())
		return
	}

	if err := oldConfig.MergeValue(valueObject); err != nil {
		SetHTTPResponse(c, -1, nil, "合并配置失败: "+err.Error())
		return
	}

	oldConfig.UpdateUser = stockConfig.UpdateUser

	if err := oldConfig.Save(); err != nil {
		SetHTTPResponse(c, -1, nil, "保存配置失败: "+err.Error())
		return
	}

	data := make(map[string]interface{})
	data["config"] = oldConfig
	SetHTTPResponse(c, 0, data, "更新成功")
}

// DeleteStockConfig
func (s *Service) DeleteStockConfig(c *gin.Context) {
	scope := "stock"
	// 删除股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}
	s.Logger.Info("DeleteStockConfig", "stockConfig", stockConfig)

	if stockConfig.StockCode == "" {
		SetHTTPResponse(c, -1, nil, "stock code 不能为空")
		return
	}

	oldConfig, err := models.GetConfig(scope, stockConfig.StockCode)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "获取原配置失败: "+err.Error())
		return
	}

	if err := oldConfig.Delete(); err != nil {
		SetHTTPResponse(c, -1, nil, "删除失败: "+err.Error())
		return
	}

	SetHTTPResponse(c, 0, nil, "删除成功")
}

// GetGlobalConfigs
func (s *Service) GetGlobalConfigs(c *gin.Context) {
	scope := "global"
	configs, err := models.GetConfigs(scope)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "获取配置失败: "+err.Error())
		return
	}

	// 返回所有全局配置
	data := make(map[string]interface{})
	data["configs"] = configs
	SetHTTPResponse(c, 0, data, "查询成功")
}

type GlobalConfig struct {
	Name   string `json:"name" binding:"required"`
	Config struct {
		Broker string `json:"broker,omitempty" binding:"required"`
	} `json:"config"`
	UpdateUser string `json:"update_user"  binding:"required"`
}

// AddGlobalConfig
func (s *Service) AddGlobalConfig(c *gin.Context) {
	scope := "global"
	// 添加全局配置
	var globalConfig GlobalConfig
	if err := c.ShouldBindJSON(&globalConfig); err != nil {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}

	s.Logger.Info("AddGlobalConfig", "globalConfig", globalConfig)

	if globalConfig.Name == "" {
		SetHTTPResponse(c, -1, nil, "name 不能为空")
		return
	}

	value, err := json.Marshal(globalConfig.Config)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "参数格式错误")
		return
	}

	var valueObject map[string]interface{}

	if err := json.Unmarshal(value, &valueObject); err != nil {
		SetHTTPResponse(c, -1, nil, "参数格式错误")
		return
	}

	scopeConfig := models.NewConfig(scope, globalConfig.Name, valueObject, globalConfig.UpdateUser)
	if err := scopeConfig.Create(); err != nil {
		SetHTTPResponse(c, -1, nil, "添加失败: "+err.Error())
		return
	}

	data := make(map[string]interface{})
	data["config"] = scopeConfig
	SetHTTPResponse(c, 0, data, "添加成功")
}

// UpdateGlobalConfig
func (s *Service) UpdateGlobalConfig(c *gin.Context) {
	scope := "global"
	// 更新全局配置
	var globalConfig GlobalConfig
	if err := c.ShouldBindJSON(&globalConfig); err != nil {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}

	s.Logger.Info("UpdateGlobalConfig", "globalConfig", globalConfig)

	value, err := json.Marshal(globalConfig.Config)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "参数格式错误")
		return
	}

	var valueObject map[string]interface{}

	if err := json.Unmarshal(value, &valueObject); err != nil {
		SetHTTPResponse(c, -1, nil, "参数格式错误")
		return
	}

	if len(valueObject) == 0 {
		SetHTTPResponse(c, -1, nil, "参数错误")
		return
	}

	oldConfig, err := models.GetConfig(scope, globalConfig.Name)
	if err != nil {
		SetHTTPResponse(c, -1, nil, "获取原配置失败: "+err.Error())
		return
	}

	if err := oldConfig.MergeValue(valueObject); err != nil {
		SetHTTPResponse(c, -1, nil, "合并配置失败: "+err.Error())
		return
	}

	oldConfig.UpdateUser = globalConfig.UpdateUser

	if err := oldConfig.Save(); err != nil {
		SetHTTPResponse(c, -1, nil, "保存配置失败: "+err.Error())
		return
	}

	data := make(map[string]interface{})
	data["config"] = oldConfig
	SetHTTPResponse(c, 0, data, "更新成功")
}
