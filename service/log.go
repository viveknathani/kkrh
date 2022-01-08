package service

import (
	"github.com/viveknathani/kkrh/entity"
)

// CreateLog will insert a new log in the database
func (service *Service) CreateLog(l *entity.Log) error {

	if !isValidLog(l, false) {
		return ErrInvalidLog
	}

	service.logger.Info("Inserting log into database.")
	err := service.repo.CreateLog(l)
	if err != nil {
		service.logger.Error(err.Error())
		return ErrNoLogInsert
	}
	service.logger.Info("Log inserted.")
	return nil
}

// StartLog will create a new log in the database with endTime as 0
func (service *Service) StartLog(l *entity.Log) error {

	l.EndTime = 0

	if !isValidLog(l, true) {
		return ErrInvalidLog
	}
	service.logger.Info("Inserting log into database.")
	err := service.repo.CreateLog(l)
	if err != nil {
		service.logger.Error(err.Error())
		return ErrNoLogInsert
	}
	service.logger.Info("Log inserted.")
	return nil
}

// EndLog will update a log's endTime
func (service *Service) EndLog(id string, endTime int64) error {

	service.logger.Info("Updating log in database.")
	err := service.repo.UpdateLog(id, endTime)
	if err != nil {
		service.logger.Error(err.Error())
		return ErrNoLogUpdate
	}
	service.logger.Info("Log updated.")
	return nil
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
