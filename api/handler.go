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
	s.POST("/close_out", s.closeOut)

	// add config update handler
	s.GET("/stock/configs", s.GetStockConfigs)
	s.POST("/stock/configs", s.AddStockConfig)
	s.PUT("/stock/configs", s.UpdateStockConfig)
	s.DELETE("/stock/configs", s.DeleteStockConfig)

}

func (s *Service) hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello",
	})
	return
}

type CloseOut struct {
	StockCode string
}

func (s *Service) closeOut(c *gin.Context) {
	var closeOut CloseOut
	if err := c.ShouldBindJSON(closeOut); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

	// do something
	quantCoreBackend := s.cfg.GetBackend("quant_core")

	// build http request to quant_core_backend
	// use http request post method
	// with json body

	// 设置请求地址和方法
	url := fmt.Sprintf("http://%s/close_out", quantCoreBackend)
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
			"message": "请求失败",
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
	StockCode  string                 `json:"stock_code"`
	Config     map[string]interface{} `json:"config"`
	UpdateUser string                 `json:"update_user"`
}

// AddStockConfig
func (s *Service) AddStockConfig(c *gin.Context) {
	s.Logger.Info("AddStockConfig", "request", c.Request.Body)
	scope := "stock"
	// 添加股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

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

	scopeConfig := models.NewConfig(scope, stockConfig.StockCode, value, "admin")
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
	s.Logger.Info("UpdateStockConfig", "request", c.Request.Body)
	scope := "stock"
	// 更新股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
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

	oldConfig, err := models.GetConfig(scope, stockConfig.StockCode)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取原配置失败: " + err.Error(),
		})
		return
	}

	if err := oldConfig.MergeValue(value); err != nil {
		c.JSON(400, gin.H{
			"message": "合并配置失败: " + err.Error(),
		})
		return
	}

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
	s.Logger.Info("DeleteStockConfig", "request", c.Request.Body)
	scope := "stock"
	// 删除股票配置
	var stockConfig StockConfig
	if err := c.ShouldBindJSON(&stockConfig); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}

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
