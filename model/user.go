package model

import (
	"database/sql"
)

// 全て値があるときに利用する
func InsertUser(db *sql.DB, userId string, region Region) error {
	q := `insert into users (user_id, region) values ($1,$2)`

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userId, region)
	return err
}

func GetUserByUserId(db *sql.DB, userId string) (User, error) {
	q := `select * from users where user_id=$1`

	var user User

	stmt, err := db.Prepare(q)
	if err != nil {
		return user, err
	}

	err = stmt.QueryRow(userId).Scan(&user.ID, &user.UserID, &user.Region, &user.Created)

	return user, err
}

func GetUsersByRegion(db *sql.DB, region Region) ([]User, error) {
	q := `select * from users where region=$1`

	stmt, err := db.Prepare(q)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(region)
	if err != nil {
		return nil, err
	}

	users, err := ScanUsers(rows)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UpdateUser(db *sql.DB, userId string, region Region) error {
	q := `update users set region=$1 where user_id=$2`
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(region, userId)
	return err
}

func IsUserExists(db *sql.DB, userId string) (bool, error) {
	q := `select count(*) from users where user_id=$1`

	stmt, err := db.Prepare(q)
	if err != nil {
		return false, err
	}

	var count int64
	err = stmt.QueryRow(userId).Scan(&count)
	return count == 1, nil
}

func GetUserIdsFromUsers(users []User) []string {
	var ids []string
	for i:=0; i<len(users); i++ {
		ids = append(ids, users[i].UserID)
	}
	return ids
}
