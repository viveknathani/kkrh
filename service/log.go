package service

import (
	"github.com/viveknathani/kkrh/entity"
)

// CreateLog will insert a new log in the database
func (service *Service) CreateLog(l *entity.Log) error {

	if !isValidLog(l, false) {
		return ErrInvalidLog
	}

	service.Logger.Info("Inserting log into database.")
	err := service.Repo.CreateLog(l)
	if err != nil {
		service.Logger.Error(err.Error())
		return ErrNoLogInsert
	}
	service.Logger.Info("Log inserted.")
	return nil
}

// StartLog will create a new log in the database with endTime as 0
func (service *Service) StartLog(l *entity.Log) error {

	l.EndTime = 0

	if !isValidLog(l, true) {
		return ErrInvalidLog
	}
	service.Logger.Info("Inserting log into database.")
	err := service.Repo.CreateLog(l)
	if err != nil {
		service.Logger.Error(err.Error())
		return ErrNoLogInsert
	}
	service.Logger.Info("Log inserted.")
	return nil
}

// EndLog will update a log's endTime
func (service *Service) EndLog(id string, endTime int64) error {

	service.Logger.Info("Updating log in database.")
	err := service.Repo.UpdateLog(id, endTime)
	if err != nil {
		service.Logger.Error(err.Error())
		return ErrNoLogUpdate
	}
	service.Logger.Info("Log updated.")
	return nil
}

// GetPendingLogs will fetch pending logs from DB for the given user.
func (service *Service) GetPendingLogs(userId string) (*[]entity.Log, error) {

	service.Logger.Info("Fetching logs from database.")
	list, err := service.Repo.GetLogs(userId, 0) // 0 signifies pending
	if err != nil {
		service.Logger.Error(err.Error())
		return nil, ErrNoLogFetch
	}
	service.Logger.Info("Logs fetched.")
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
	check = check && (l.EndTime != 0)
	check = check && (l.Activity != "")
	return check
}
