package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"Resume/backend/models"
	"Resume/backend/services"

	"github.com/gin-gonic/gin"
)

type RecruitmentSourceHandler struct {
	service *services.RecruitmentSourceService
}

func NewRecruitmentSourceHandler(service *services.RecruitmentSourceService) *RecruitmentSourceHandler {
	return &RecruitmentSourceHandler{service: service}
}

func (h *RecruitmentSourceHandler) CreateSource(c *gin.Context) {
	var req models.RecruitmentSourceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("无效的请求载荷: %s", err.Error())})
		return
	}
	source, err := h.service.CreateSource(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("创建招聘来源失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, models.RecruitmentSourceSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "招聘来源创建成功"},
		Source:        *source,
	})
}

func (h *RecruitmentSourceHandler) GetSource(c *gin.Context) {
	id := c.Param("id")
	source, err := h.service.GetSourceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": fmt.Sprintf("招聘来源未找到: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.RecruitmentSourceSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "招聘来源获取成功"},
		Source:        *source,
	})
}

func (h *RecruitmentSourceHandler) ListSources(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	order := c.DefaultQuery("order", "desc")
	keyword := c.DefaultQuery("keyword", "")
	sources, total, err := h.service.ListSources(page, limit, sortBy, order, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("获取招聘来源列表失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.RecruitmentSourceListResponse{
		BasicResponse: models.BasicResponse{Success: true},
		Sources:       sources,
		Total:         total,
		Page:          page,
		Limit:         limit,
	})
}

func (h *RecruitmentSourceHandler) UpdateSource(c *gin.Context) {
	id := c.Param("id")
	var req models.RecruitmentSourceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("无效的请求载荷: %s", err.Error())})
		return
	}
	source, err := h.service.UpdateSource(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("更新招聘来源失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.RecruitmentSourceSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "招聘来源更新成功"},
		Source:        *source,
	})
}

func (h *RecruitmentSourceHandler) DeleteSource(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteSource(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("删除招聘来源失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "招聘来源删除成功"})
}
