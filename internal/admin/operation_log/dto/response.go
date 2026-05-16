package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

// OperationLogResponse 操作日志响应结构体
type OperationLogResponse struct {
	ID int64   `json:"id"`
	Username   string `json:"username"`
	Ip         string `json:"ip"`
	IpLocation string `json:"ip_location"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Desc       string `json:"desc"`
	Status     int    `json:"status"`
	StartTime  string `json:"start_time"`
	TimeCost   int64  `json:"time_cost"`
	UserAgent  string `json:"user_agent"`
}

func toOperationLogResponse(log model.OperationLog) OperationLogResponse {
	return OperationLogResponse{
		ID:         log.ID,
		Username:   log.Username,
		Ip:         log.Ip,
		IpLocation: log.IpLocation,
		Method:     log.Method,
		Path:       log.Path,
		Desc:       log.Desc,
		Status:     log.Status,
		StartTime:  log.StartTime.Format(constant.TIME_FORMAT),
		TimeCost:   log.TimeCost,
		UserAgent:  log.UserAgent,
	}
}

// ToOperationLogListResponse 转换操作日志列表为Response
func ToOperationLogListResponse(logList []model.OperationLog) []OperationLogResponse {
	if logList == nil {
		return []OperationLogResponse{}
	}

	logs := make([]OperationLogResponse, 0, len(logList))
	for _, log := range logList {
		logs = append(logs, toOperationLogResponse(log))
	}

	return logs
}
