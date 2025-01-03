package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int64
	DateUntil string
}

func main() {
	db, err := sql.Open("mysql", "root:rootroot@/outline_db")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	rows, err := SelectAllUsers(db)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		var ID_Str string
		var isTimePassed bool = false

		err := rows.Scan(&user.ID, &user.DateUntil)

		if err != nil {
			panic(err)
		}

		ID_Str = strconv.FormatInt(user.ID, 10)
		isTimePassed, err = CheckIfItPasses(user.DateUntil)

		if err != nil {
			panic(err)
		}

		if isTimePassed {
			_, err = DeleteUser(db, user.ID)
		}

		if err != nil {
			fmt.Println("Не удалось удалить запись в таблице users с ID = " + ID_Str)
			panic(err)
		} else {
			fmt.Println("User с ID = " + ID_Str + " был удален")
		}
	}

	err = rows.Err()

	if err != nil {
		panic(err)
	}

}

func CheckIfItPasses(datetime string) (bool, error) {
	result := false
	dateUntilTime, err := time.Parse("2006-01-02 15:04:05", datetime)

	if err != nil {
		return false, err
	}

	result = time.Now().After(dateUntilTime)

	return result, nil
}

func SelectAllUsers(db *sql.DB) (*sql.Rows, error) {
	result, err := db.Query("SELECT id, date_until FROM users")

	if err != nil {
		return nil, err
	}

	return result, nil
}

func InsertUser(db *sql.DB, name string, datetime time.Time) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO users SET name = ?, datetime = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, datetime)

	if err != nil {
		return 0, err
	}

	lastInsertId, _ := result.LastInsertId()

	return lastInsertId, nil
}

func SelectUser(db *sql.DB, id int64) (*sql.Row, error) {
	stmt, err := db.Prepare("SELECT * FROM users WHERE id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	result := stmt.QueryRow(id)

	return result, nil
}

func DeleteUser(db *sql.DB, id int64) (sql.Result, error) {
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(id)

	if err != nil {
		return nil, err
	}

	return result, nil
}
