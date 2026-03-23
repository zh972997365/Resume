package models

type BasicResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UploadResponse struct {
	BasicResponse
	File   *FileInfo   `json:"file,omitempty"`
	Files  []*FileInfo `json:"files,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

type FileListResponse struct {
	BasicResponse
	Files []*FileInfo `json:"files"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

type FileSearchResponse struct {
	BasicResponse
	Files []*FileInfo `json:"files"`
}

type InterviewListResponse struct {
	BasicResponse
	Interviews []Interview `json:"interviews"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}

type InterviewSingleResponse struct {
	BasicResponse
	Interview Interview `json:"interview"`
}

type CompanyPositionListResponse struct {
	BasicResponse
	Positions []CompanyPosition `json:"positions"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
}

type CompanyPositionSingleResponse struct {
	BasicResponse
	Position CompanyPosition `json:"position"`
}

type RecruitmentSourceListResponse struct {
	BasicResponse
	Sources []RecruitmentSource `json:"sources"`
	Total   int64               `json:"total"`
	Page    int                 `json:"page"`
	Limit   int                 `json:"limit"`
}

type RecruitmentSourceSingleResponse struct {
	BasicResponse
	Source RecruitmentSource `json:"source"`
}

type EmployeeListResponse struct {
	BasicResponse
	Employees []Employee `json:"employees"`
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}

type EmployeeSingleResponse struct {
	BasicResponse
	Employee Employee `json:"employee"`
}
