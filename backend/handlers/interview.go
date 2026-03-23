package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"Resume/backend/models"
	"Resume/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InterviewHandler struct {
	interviewService *services.InterviewService
	storageService   *services.StorageService
}

func NewInterviewHandler(interviewService *services.InterviewService, storageService *services.StorageService) *InterviewHandler {
	return &InterviewHandler{
		interviewService: interviewService,
		storageService:   storageService,
	}
}

func (h *InterviewHandler) CreateInterview(c *gin.Context) {
	var req models.InterviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 获取友好的错误消息
		errorMsg := getFriendlyErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errorMsg,
		})
		return
	}

	interview, err := h.interviewService.CreateInterview(&req)
	if err != nil {
		// 检查是否是简历重复使用的错误
		if strings.Contains(err.Error(), "该简历文件已被其他面试记录使用") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("创建面试记录失败: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusCreated, models.InterviewSingleResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: "面试记录创建成功",
		},
		Interview: *interview,
	})
}

// 友好的错误消息转换函数
func getFriendlyErrorMessage(err error) string {
	// 导入包：github.com/go-playground/validator/v10
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// 将技术性错误消息转换为用户友好的提示
		fieldMap := map[string]string{
			"CandidateName":       "姓名",
			"PhoneNumber":         "手机号",
			"Email":               "邮箱",
			"CompanyPositionID":   "应聘岗位",
			"ResumeFileID":        "简历文件",
			"RecruitmentSourceID": "招聘来源",
			"InterviewRound":      "面试轮次",
			"InterviewMethod":     "面试形式",
			"InterviewTime":       "面试时间",
			"InterviewerID":       "面试官",
			"Suggestion":          "面试建议",
		}

		var friendlyErrors []string
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			chineseName := fieldMap[fieldName]
			if chineseName == "" {
				chineseName = fieldName
			}

			switch fieldErr.Tag() {
			case "required":
				friendlyErrors = append(friendlyErrors, chineseName+"不能为空")
			case "email":
				friendlyErrors = append(friendlyErrors, chineseName+"格式不正确")
			default:
				friendlyErrors = append(friendlyErrors, chineseName+"填写有误")
			}
		}

		if len(friendlyErrors) > 0 {
			if len(friendlyErrors) == 1 {
				return friendlyErrors[0]
			}
			return "请检查以下必填项: " + strings.Join(friendlyErrors, "，")
		}
	}

	// 默认错误消息
	return "请求参数有误，请检查填写内容"
}

// 翻译验证错误为中文
func translateValidationErrors(errors validator.ValidationErrors) string {
	var errorMessages []string

	for _, err := range errors {
		fieldName := err.Field()
		tag := err.Tag()

		// 字段名称映射为中文
		chineseFieldName := map[string]string{
			"CandidateName":   "姓名",
			"PhoneNumber":     "手机号",
			"Email":           "邮箱",
			"AppliedPosition": "应聘岗位",
			"ResumeFileID":    "简历文件",
			"Source":          "招聘来源",
			"InterviewRound":  "面试轮次",
			"InterviewMethod": "面试形式",
			"InterviewTime":   "面试时间",
			"InterviewerName": "面试官",
			"Suggestion":      "面试建议",
		}[fieldName]

		if chineseFieldName == "" {
			chineseFieldName = fieldName
		}

		// 验证规则映射为中文提示
		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("%s不能为空", chineseFieldName)
		case "email":
			message = fmt.Sprintf("%s格式不正确", chineseFieldName)
		case "min":
			message = fmt.Sprintf("%s长度不够", chineseFieldName)
		case "max":
			message = fmt.Sprintf("%s长度超限", chineseFieldName)
		default:
			message = fmt.Sprintf("%s验证失败", chineseFieldName)
		}

		errorMessages = append(errorMessages, message)
	}

	if len(errorMessages) == 1 {
		return errorMessages[0]
	}
	return "以下字段填写有误: " + strings.Join(errorMessages, "，")
}

func (h *InterviewHandler) GetInterview(c *gin.Context) {
	id := c.Param("id")
	interview, err := h.interviewService.GetInterviewByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": fmt.Sprintf("面试记录未找到: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.InterviewSingleResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: "面试记录获取成功",
		},
		Interview: *interview,
	})
}

func (h *InterviewHandler) ListInterviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	order := c.DefaultQuery("order", "desc")
	keyword := c.DefaultQuery("keyword", "")
	suggestion := models.Suggestion(c.DefaultQuery("suggestion", ""))

	interviews, total, err := h.interviewService.ListInterviews(page, limit, sortBy, order, keyword, suggestion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("获取面试记录失败: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, models.InterviewListResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
		},
		Interviews: interviews,
		Total:      total,
		Page:       page,
		Limit:      limit,
	})
}

func (h *InterviewHandler) UpdateInterview(c *gin.Context) {
	id := c.Param("id")
	var req models.InterviewUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("无效的请求载荷: %s", err.Error())})
		return
	}

	interview, err := h.interviewService.UpdateInterview(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("更新面试记录失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, models.InterviewSingleResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: "面试记录更新成功",
		},
		Interview: *interview,
	})
}

func (h *InterviewHandler) DeleteInterview(c *gin.Context) {
	id := c.Param("id")
	err := h.interviewService.DeleteInterview(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("删除面试记录失败: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "面试记录删除成功"})
}

func (h *InterviewHandler) UploadResumeForInterview(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.UploadResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "文件上传失败: " + err.Error(),
			},
		})
		return
	}

	fileInfo, err := h.storageService.Save(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UploadResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "文件保存失败: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: "简历文件上传成功",
		},
		File: fileInfo,
	})
}
