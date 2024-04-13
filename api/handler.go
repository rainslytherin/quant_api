package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (s *Service) closeOut(c *gin.Context) {
	var closeOut CloseOut
	if err := c.ShouldBindJSON(&closeOut); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	s.Logger.Info("closeOut", "closeOut", closeOut)

	if closeOut.StockCode == "" {
		c.JSON(400, gin.H{
			"message": "stock_code 不能为空",
		})
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

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "创建请求失败",
		})
		return
	}

	// 设置请求头部信息
	req.Header.Set("Content-Type", "application/json")

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "请求失败: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "解析响应失败",
		})
		return
	}

	// 反序列化结果
	type CloseOutResponse struct {
		StockPrice float64 `json:"stockPrice"`
		Qty        int     `json:"qty"`
		Message    string  `json:"message"`
	}

	var closeOutResponse CloseOutResponse

	if err := json.Unmarshal(body, &closeOutResponse); err != nil {
		c.JSON(400, gin.H{
			"message": "反序列化失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"stockPrice": closeOutResponse.StockPrice,
		"qty":        closeOutResponse.Qty,
		"message":    "执行完成",
	})
}

// GetStockConfigs
func (s *Service) GetStockConfigs(c *gin.Context) {
	scope := "stock"
	configs, err := models.GetConfigs(scope)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取配置失败: " + err.Error(),
		})
		return
	}

	// 返回所有股票配置
	c.JSON(200, gin.H{
		"stockConfigs": configs,
	})
}

type StockConfig struct {
	StockCode string `json:"stock_code"  binding:"required"`
	Config    struct {
		ProdStatus bool    `json:"prod_status,omitempty"`
		PreStatus  bool    `json:"pre_status,omitempty"`
		UpLimit    float64 `json:"up_limit,omitempty"`
		LowLimit   float64 `json:"low_limit,omitempty"`
	} `json:"config"`
	UpdateUser string `json:"update_user"  binding:"required"`
}

// AddStockConfig
func (s *Service) AddStockConfig(c *gin.Context) {
	scope := "stock"
	// 添加股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	s.Logger.Info("AddStockConfig", "stockConfig", stockConfig)

	if stockConfig.StockCode == "" {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	value, err := json.Marshal(stockConfig.Config)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "参数格式错误",
		})
		return
	}

	var valueObject map[string]interface{}

	if err := json.Unmarshal(value, &valueObject); err != nil {
		c.JSON(400, gin.H{
			"message": "参数格式错误",
		})
		return
	}

	scopeConfig := models.NewConfig(scope, stockConfig.StockCode, valueObject, stockConfig.UpdateUser)
	if err := scopeConfig.Create(); err != nil {
		c.JSON(400, gin.H{
			"message": "添加失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "添加成功",
	})
}

// UpdateStockConfig
func (s *Service) UpdateStockConfig(c *gin.Context) {
	scope := "stock"
	// 更新股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	s.Logger.Info("UpdateStockConfig", "stockConfig", stockConfig)

	value, err := json.Marshal(stockConfig.Config)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "参数格式错误",
		})
		return
	}

	var valueObject map[string]interface{}

	if err := json.Unmarshal(value, &valueObject); err != nil {
		c.JSON(400, gin.H{
			"message": "参数格式错误",
		})
		return
	}

	if len(valueObject) == 0 {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	oldConfig, err := models.GetConfig(scope, stockConfig.StockCode)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取原配置失败: " + err.Error(),
		})
		return
	}

	if err := oldConfig.MergeValue(valueObject); err != nil {
		c.JSON(400, gin.H{
			"message": "合并配置失败: " + err.Error(),
		})
		return
	}

	oldConfig.UpdateUser = stockConfig.UpdateUser

	if err := oldConfig.Save(); err != nil {
		c.JSON(400, gin.H{
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "更新成功",
	})
}

// DeleteStockConfig
func (s *Service) DeleteStockConfig(c *gin.Context) {
	scope := "stock"
	// 删除股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}
	s.Logger.Info("DeleteStockConfig", "stockConfig", stockConfig)

	if stockConfig.StockCode == "" {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	oldConfig, err := models.GetConfig(scope, stockConfig.StockCode)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取原配置失败: " + err.Error(),
		})
		return
	}

	if err := oldConfig.Delete(); err != nil {
		c.JSON(400, gin.H{
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "删除成功",
	})
}

// GetGlobalConfigs
func (s *Service) GetGlobalConfigs(c *gin.Context) {
	scope := "global"
	configs, err := models.GetConfigs(scope)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取配置失败: " + err.Error(),
		})
		return
	}

	// 返回所有全局配置
	c.JSON(200, gin.H{
		"globalConfigs": configs,
	})
}

type GlobalConfig struct {
	Name       string `json:"name"  binding:"required"`
	Config     struct {
		Broker  string `json:"broker,omitempty"`
	}
	UpdateUser string `json:"update_user"  binding:"required"`
}

// AddGlobalConfig
func (s *Service) AddGlobalConfig(c *gin.Context) {
	scope := "global"
	// 添加全局配置
	var globalConfig GlobalConfig
	if err := c.ShouldBindJSON(&globalConfig); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	s.Logger.Info("AddGlobalConfig", "globalConfig", globalConfig)

	if globalConfig.Name == "" {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	value, err := json.Marshal(globalConfig.Config)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "参数格式错误",
		})
		return
	}

	var valueObject map[string]interface{}

	if err := json.Unmarshal(value, &valueObject); err != nil {
		c.JSON(400, gin.H{
			"message": "参数格式错误",
		})
		return
	}

	scopeConfig := models.NewConfig(scope, globalConfig.Name, valueObject, globalConfig.UpdateUser)
	if err := scopeConfig.Create(); err != nil {
		c.JSON(400, gin.H{
			"message": "添加失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "添加成功",
	})
}

// UpdateGlobalConfig
func (s *Service) UpdateGlobalConfig(c *gin.Context) {
	scope := "global"
	// 更新全局配置
	var globalConfig GlobalConfig
	if err := c.ShouldBindJSON(&globalConfig); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	s.Logger.Info("UpdateGlobalConfig", "globalConfig", globalConfig)

	value, err := json.Marshal(globalConfig.Config)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "参数格式错误",
		})
		return
	}

	var valueObject map[string]interface{}

	if err := json.Unmarshal(value, &valueObject); err != nil {
		c.JSON(400, gin.H{
			"message": "参数格式错误",
		})
		return
	}

	if len(valueObject) == 0 {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	oldConfig, err := models.GetConfig(scope, globalConfig.Name)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取原配置失败: " + err.Error(),
		})
		return
	}

	if err := oldConfig.MergeValue(valueObject); err != nil {
		c.JSON(400, gin.H{
			"message": "合并配置失败: " + err.Error(),
		})
		return
	}

	oldConfig.UpdateUser = globalConfig.UpdateUser

	if err := oldConfig.Save(); err != nil {
		c.JSON(400, gin.H{
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "更新成功",
	})
}