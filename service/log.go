package service

import (
	"context"

	"github.com/viveknathani/kkrh/entity"
)

// CreateLog will insert a new log in the database
func (service *Service) CreateLog(ctx context.Context, l *entity.Log) error {

	if !isValidLog(l, false) {
		return ErrInvalidLog
	}

	service.Logger.Info("database: insert log start.", zapReqID(ctx))
	err := service.Repo.CreateLog(l)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrNoLogInsert
	}
	service.Logger.Info("database: insert log complete.", zapReqID(ctx))
	return nil
}

// StartLog will create a new log in the database with endTime as 0
func (service *Service) StartLog(ctx context.Context, l *entity.Log) error {

	l.EndTime = 0

	if !isValidLog(l, true) {
		return ErrInvalidLog
	}
	service.Logger.Info("database: insert log start.", zapReqID(ctx))
	err := service.Repo.CreateLog(l)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrNoLogInsert
	}
	service.Logger.Info("database: insert log complete.", zapReqID(ctx))
	return nil
}

// EndLog will update a log's endTime
func (service *Service) EndLog(ctx context.Context, userId string, id string, endTime int64) error {

	service.Logger.Info("database: update log start.", zapReqID(ctx))
	err := service.Repo.UpdateLog(userId, id, endTime)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return ErrNoLogUpdate
	}
	service.Logger.Info("database: update log complete.", zapReqID(ctx))
	return nil
}

// GetPendingLogs will fetch pending logs from DB for the given user.
func (service *Service) GetPendingLogs(ctx context.Context, userId string) (*[]entity.Log, error) {

	service.Logger.Info("database: fetch logs start.", zapReqID(ctx))
	list, err := service.Repo.GetPendingLogs(userId, 0) // 0 signifies pending
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return nil, ErrNoLogFetch
	}
	service.Logger.Info("database: fetch logs complete.", zapReqID(ctx))
	return list, nil
}

func isValidLog(l *entity.Log, ignoreEndTime bool) bool {

	if l == nil {
		return false
	}

	check := true

	if !ignoreEndTime {
		check = check && (l.EndTime > l.StartTime)
	}

	check = check && (l.Latitude != 0)
	check = check && (l.Longitude != 0)
	check = check && (l.StartTime != 0)
	check = check && (l.Activity != "")
	return check
}
