package repository

import "github.com/viveknathani/kkrh/entity"

// userRepository is the method set to operate on the User entity.
type userRepository interface {

	// CreateUser will create a new user in the database and will
	// have a newly generated UUID.
	CreateUser(u *entity.User) error

	// GetUser will fetch the first found user from the database.
	GetUser(email string) (*entity.User, error)
}

// logRepository is the method set to operate on the Log entity.
type logRepository interface {

	// CreateLog will create a new log in the database and will have
	// a newly generated UUID.
	CreateLog(l *entity.Log) error

	// UpdateLog will update an existing log entry with the new value of
	// endTime.
	UpdateLog(logId string, endTime int64) error

	// GetLogs will fetch all logs for the user where stopTime matches the
	// given endTime.
	GetLogs(userID string, endTime int64) (*[]entity.Log, error)
}

// Repository encapsulates the method set for all entities.
type Repository interface {
	userRepository
	logRepository
}
