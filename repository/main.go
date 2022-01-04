package repository

// userRepository is the method set to operate on the User model
type userRepository interface{}

// logRepository is the method set to operate on the Log model
type logRepository interface{}

// Repository encapsulates the method set for all models
type Repository interface {
	userRepository
	logRepository
}
