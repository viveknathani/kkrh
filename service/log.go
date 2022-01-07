package service

import (
	"log"

	"github.com/viveknathani/kkrh/entity"
)

// CreateLog will insert a new log in the database
func (service *Service) CreateLog(l *entity.Log) error {

	if !isValidLog(l, false) {
		return ErrInvalidLog
	}
	err := service.repo.CreateLog(l)
	if err != nil {
		log.Print(err)
		return ErrNoLogInsert
	}
	return nil
}

// StartLog will create a new log in the database with endTime as 0
func (service *Service) StartLog(l *entity.Log) error {

	l.EndTime = 0

	if !isValidLog(l, true) {
		return ErrInvalidLog
	}
	err := service.repo.CreateLog(l)
	if err != nil {
		log.Print(err)
		return ErrNoLogInsert
	}
	return nil
}

// EndLog will update a log's endTime
func (service *Service) EndLog(id string, endTime int64) error {

	err := service.repo.UpdateLog(id, endTime)
	if err != nil {
		log.Print(err)
		return ErrNoLogUpdate
	}
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
