package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"Resume/backend/models"
	"Resume/backend/services"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req models.EmployeeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("无效的请求载荷: %s", err.Error())})
		return
	}
	employee, err := h.service.CreateEmployee(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("创建面试官失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, models.EmployeeSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "面试官创建成功"},
		Employee:      *employee,
	})
}

func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id := c.Param("id")
	employee, err := h.service.GetEmployeeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": fmt.Sprintf("面试官未找到: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.EmployeeSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "面试官获取成功"},
		Employee:      *employee,
	})
}

func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	order := c.DefaultQuery("order", "desc")
	keyword := c.DefaultQuery("keyword", "")
	employees, total, err := h.service.ListEmployees(page, limit, sortBy, order, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("获取面试官列表失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.EmployeeListResponse{
		BasicResponse: models.BasicResponse{Success: true},
		Employees:     employees,
		Total:         total,
		Page:          page,
		Limit:         limit,
	})
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var req models.EmployeeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("无效的请求载荷: %s", err.Error())})
		return
	}
	employee, err := h.service.UpdateEmployee(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("更新面试官失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.EmployeeSingleResponse{
		BasicResponse: models.BasicResponse{Success: true, Message: "面试官更新成功"},
		Employee:      *employee,
	})
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteEmployee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("删除面试官失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "面试官删除成功"})
}
