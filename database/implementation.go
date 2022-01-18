package database

// This file contains the implementation for the Repository interface.

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/viveknathani/kkrh/entity"
)

const (
	statementInsertUser                   = "insert into users (id, name, email, password) values ($1, $2, $3, $4);"
	statementSelectUserFromEmail          = "select * from users where email = $1;"
	statementInsertLog                    = "insert into logs (id, userId, latitude, longitude, activity, startTime, endTime, notes) values ($1, $2, $3, $4, $5, $6, $7, $8);"
	statementSelectLogsFromUserAndEndTime = "select * from logs where userId = $1 and endTime = $2;"
	statementUpdateLogsWithIdAndEndTime   = "update logs set endTime = $1 where id = $2 and userId = $3;"
	statementDeleteUser                   = "delete from users where id = $1;"
	statementDeleteLog                    = "delete from logs where id = $1;"
)

// CreateUser will create a new user in the database and will
// have a newly generated UUID.
func (db *Database) CreateUser(u *entity.User) error {

	u.Id = uuid.New().String()
	err := db.execWithTransaction(statementInsertUser, u.Id, u.Name, u.Email, u.Password)
	return err
}

// GetUser will fetch the first found user from the database.
func (db *Database) GetUser(email string) (*entity.User, error) {

	var u entity.User
	exists := false
	err := db.queryWithTransaction(statementSelectUserFromEmail, func(rows *sql.Rows) error {

		//We iterate only once since we are interested in the first occurence
		if rows.Next() {
			err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password)
			if err != nil {
				return err
			}
			exists = true
		}
		return nil
	}, email)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}
	return &u, nil
}

// CreateLog will create a new log in the database and will have
// a newly generated UUID.
func (db *Database) CreateLog(l *entity.Log) error {

	l.Id = uuid.New().String()
	err := db.execWithTransaction(statementInsertLog, l.Id, l.UserId, l.Latitude, l.Longitude, l.Activity, l.StartTime, l.EndTime, l.Notes)
	return err
}

// UpdateLog will update an existing log entry with the new value of
// endTime.
func (db *Database) UpdateLog(userID string, logId string, endTime int64) error {

	err := db.execWithTransaction(statementUpdateLogsWithIdAndEndTime, endTime, logId, userID)
	return err
}

// GetPendingLogs will fetch all logs for the user where stopTime matches the
// given endTime.
func (db *Database) GetPendingLogs(userID string, endTime int64) (*[]entity.Log, error) {

	result := make([]entity.Log, 0)
	err := db.queryWithTransaction(statementSelectLogsFromUserAndEndTime, func(rows *sql.Rows) error {
		for rows.Next() {

			var l entity.Log
			err := rows.Scan(&l.Id, &l.UserId, &l.Latitude, &l.Longitude, &l.Activity, &l.StartTime, &l.EndTime, &l.Notes)
			if err != nil {
				return err
			}
			result = append(result, l)
		}
		return nil
	}, userID, endTime)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteUser will delete a user specified by userId.
func (db *Database) DeleteUser(id string) error {
	return db.execWithTransaction(statementDeleteUser, id)
}

// DeleteLog will delete a user specified by logId.
func (db *Database) DeleteLog(id string) error {

	return db.execWithTransaction(statementDeleteLog, id)
}
