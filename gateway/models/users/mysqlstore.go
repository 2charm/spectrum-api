package users

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" //mysql driver
	// "github.com/info441/assignments-andrewhwang10/servers/gateway/indexes"
)

//MySQLStore represents a users.Store backed by MySQL.
type MySQLStore struct {
	Client *sql.DB
}

//NewMySQLStore constructs a new MySQLStore
func NewMySQLStore(db *sql.DB) *MySQLStore {
	//initialize and return a new MySQLStore struct
	if db != nil {
		return &MySQLStore{
			Client: db,
		}
	}
	return nil
}

//Store implementation

//GetByID returns the User with the given ID
func (mss *MySQLStore) GetByID(id int64) (*User, error) {
	user := &User{}
	row := mss.Client.QueryRow("select * from users where id=?", id)
	if err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
		&user.FirstName, &user.LastName); err != nil {
		return nil, err
	}
	return user, nil
}

//GetByEmail returns the User with the given email
func (mss *MySQLStore) GetByEmail(email string) (*User, error) {
	user := &User{}
	row := mss.Client.QueryRow("select * from users where email=?", email)
	if err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
		&user.FirstName, &user.LastName); err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

//GetByUserName returns the User with the given Username
func (mss *MySQLStore) GetByUserName(username string) (*User, error) {
	user := &User{}
	row := mss.Client.QueryRow("select * from users where user_name=?", username)
	if err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
		&user.FirstName, &user.LastName); err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (mss *MySQLStore) Insert(user *User) (*User, error) {
	insq := "insert into users(email, pass_hash, user_name, first_name, last_name) values (?, ?, ?, ?, ?)"
	res, err := mss.Client.Exec(insq, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName)
	if err != nil {
		log.Printf("Issue executing sql statement: %v", err)
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Print("Error retrieving last insert id")
		return nil, err
	}
	user.ID = id
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (mss *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	insq := "update users set first_name=?, last_name=? where id=?"
	_, err := mss.Client.Exec(insq, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, err
	}
	return mss.GetByID(id)
}

//Delete deletes the user with the given ID
func (mss *MySQLStore) Delete(id int64) error {
	insq := "delete from users where id=?"
	_, err := mss.Client.Exec(insq, id)
	if err != nil {
		return ErrUserNotFound
	}
	return nil
}
