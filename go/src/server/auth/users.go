package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const hardness = 10000

func hash(username, password string) string {
	input := make([]byte, len(salt))
	copy(input, salt)
	input = append(input, (username + password)...)
	for i := 0; i < hardness; i++ {
		result := sha256.Sum256(input)
		input = result[:]
	}
	return base64.StdEncoding.EncodeToString(input)
}

func getUserId(username string) (int, error) {
	var userId int
	err := db.QueryRow("SELECT rowid FROM users WHERE username = ?", username).Scan(&userId)
	return userId, err
}

func AddUser(username, password string) error {
	_, err := db.Exec("INSERT INTO users (username, hash) VALUES (?, ?)", username, hash(username, password))
	return err
}

func AddRole(username, role string) error {
	userId, err := getUserId(username)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO user_roles (userid, role) VALUES (?, ?)", userId, role)
	return err
}

func MakeParent(parentUsername, childUsername string) error {
	parentId, err := getUserId(parentUsername)
	if err != nil {
		return err
	}
	childId, err := getUserId(childUsername)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO user_inheritance (parentid, childid) VALUES (?, ?)", parentId, childId)
	return err
}

func CheckPassword(username, password string) (bool, error) {
	var correctHash string
	err := db.QueryRow("SELECT hash FROM users WHERE username = ?", username).Scan(&correctHash)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	match := 1 == subtle.ConstantTimeCompare([]byte(correctHash), []byte(hash(username, password)))
	return match, nil
}

func getRoles(seen map[int]bool, roles map[string]bool, id int) error {
	seen[id] = true

	rows, err := db.Query("SELECT role FROM user_roles WHERE userid = ?", id)
	if err != nil {
		return err
	}
	for rows.Next() {
		var role string
		err = rows.Scan(&role)
		if err != nil {
			return err
		}
		roles[role] = true
	}

	rows, err = db.Query("SELECT childid FROM user_inheritance WHERE parentid = ?", id)
	if err != nil {
		return err
	}
	for rows.Next() {
		var childId int
		err = rows.Scan(&childId)
		if err != nil {
			return err
		}
		if !seen[childId] {
			err = getRoles(seen, roles, childId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetRoles(username string) (map[string]bool, error) {
	seen := make(map[int]bool)
	roles := make(map[string]bool)

	userId, err := getUserId(username)
	if err == sql.ErrNoRows {
		return roles, nil
	}
	if err != nil {
		return roles, err
	}

	err = getRoles(seen, roles, userId)
	return roles, err
}

func PrintUsers() error {
	var rowid int
	var username, child, role string
	users, err := db.Query("SELECT rowid, username FROM users")
	if err != nil {
		return err
	}

	for users.Next() {
		err = users.Scan(&rowid, &username)
		if err != nil {
			return err
		}
		fmt.Printf("user: %s\n", username)

		children, err := db.Query("SELECT users.username FROM users INNER JOIN user_inheritance ON users.rowid = user_inheritance.childid WHERE user_inheritance.parentid = ?", rowid)
		if err != nil {
			return err
		}
		fmt.Printf("children:")
		for children.Next() {
			err = children.Scan(&child)
			if err != nil {
				return err
			}
			fmt.Printf(" %s", child)
		}
		fmt.Printf("\n")

		roles, err := db.Query("SELECT role FROM user_roles WHERE userid = ?", rowid)
		if err != nil {
			return err
		}
		fmt.Printf("roles:")
		for roles.Next() {
			err = roles.Scan(&role)
			if err != nil {
				return err
			}
			fmt.Printf(" %s", role)
		}
		fmt.Printf("\n")
	}

	return nil
}
