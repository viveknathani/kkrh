package service

import (
	"context"

	"github.com/viveknathani/kkrh/entity"
)

// GetLogsInRange will fetch all logs for the user within the given range of startTime and endTime.
func (service *Service) GetLogsInRange(ctx context.Context, userId string, startTime int64, endTime int64) (*[]entity.Log, error) {

	service.Logger.Info("database: fetch logs start.", zapReqID(ctx))
	list, err := service.Repo.GetLogsInRange(userId, startTime, endTime)
	if err != nil {
		service.Logger.Error(err.Error(), zapReqID(ctx))
		return nil, ErrNoLogFetch
	}
	service.Logger.Info("database: fetch logs complete.", zapReqID(ctx))
	return list, nil
}
