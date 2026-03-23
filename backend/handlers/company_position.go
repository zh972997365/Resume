package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"Resume/backend/models"
	"Resume/backend/services"

	"github.com/gin-gonic/gin"
)

type CompanyPositionHandler struct {
	service *services.CompanyPositionService
}

func NewCompanyPositionHandler(service *services.CompanyPositionService) *CompanyPositionHandler {
	return &CompanyPositionHandler{service: service}
}

func (h *CompanyPositionHandler) CreatePosition(c *gin.Context) {
	var req models.CompanyPositionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("无效的请求载荷: %s", err.Error())})
		return
	}
	position, err := h.service.CreatePosition(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("创建岗位失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, models.CompanyPositionSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "岗位创建成功"},
		Position:      *position,
	})
}

func (h *CompanyPositionHandler) GetPosition(c *gin.Context) {
	id := c.Param("id") // ID is uint, but param is string
	position, err := h.service.GetPositionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": fmt.Sprintf("岗位未找到: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.CompanyPositionSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "岗位获取成功"},
		Position:      *position,
	})
}

func (h *CompanyPositionHandler) ListPositions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	order := c.DefaultQuery("order", "desc")
	keyword := c.DefaultQuery("keyword", "")
	positions, total, err := h.service.ListPositions(page, limit, sortBy, order, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("获取岗位列表失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.CompanyPositionListResponse{
		BasicResponse: models.BasicResponse{Success: true},
		Positions:     positions,
		Total:         total,
		Page:          page,
		Limit:         limit,
	})
}

func (h *CompanyPositionHandler) UpdatePosition(c *gin.Context) {
	id := c.Param("id")
	var req models.CompanyPositionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("无效的请求载荷: %s", err.Error())})
		return
	}
	position, err := h.service.UpdatePosition(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("更新岗位失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.CompanyPositionSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "岗位更新成功"},
		Position:      *position,
	})
}

func (h *CompanyPositionHandler) DeletePosition(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeletePosition(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("删除岗位失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "岗位删除成功"})
}
