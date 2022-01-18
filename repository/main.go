package repository

import "github.com/viveknathani/kkrh/entity"

// userRepository is the method set to operate on the User entity.
type userRepository interface {

	// CreateUser will create a new user in the database and will
	// have a newly generated UUID.
	CreateUser(u *entity.User) error

	// GetUser will fetch the first found user from the database.
	GetUser(email string) (*entity.User, error)

	// DeleteUser will delete a user specified by userId.
	DeleteUser(id string) error
}

// logRepository is the method set to operate on the Log entity.
type logRepository interface {

	// CreateLog will create a new log in the database and will have
	// a newly generated UUID.
	CreateLog(l *entity.Log) error

	// UpdateLog will update an existing log entry with the new value of
	// endTime.
	UpdateLog(userID string, logId string, endTime int64) error

	// GetPendingLogs will fetch all logs for the user where stopTime matches the
	// given endTime.
	GetPendingLogs(userID string, endTime int64) (*[]entity.Log, error)

	// GetLogsInRange will fetch all logs for the user within the given range of startTime and endTime.
	GetLogsInRange(userID string, startTime int64, endTime int64) (*[]entity.Log, error)

	// DeleteLog will delete a user specified by logId.
	DeleteLog(id string) error
}

// Repository encapsulates the method set for all entities.
type Repository interface {
	userRepository
	logRepository
}
