package news

import (
	"database/sql"
	"log"
	"time"
)

//Store represents a store for News related entries
type Store interface {
	//GetByUserID returns the metrics for a given UserID
	GetByUserID(userID int64) (*Metrics, error)

	//InsertArticle inserts a new article based on category provided
	InsertArticle(category string) error

	//GetIDOfCategory returns the id of the category provided
	getIDOfCategory(category string) (int, error)

	//GetCategoryByID returns the name of the category provided
	getCategoryByID(categoryID int) (string, error)
}

//ArticleStore represents a SQL implemented databse for News related entries
type ArticleStore struct {
	Client *sql.DB
}

func (as *ArticleStore) GetByUserID(userID int64) (*Metrics, error) {
	metrics := &Metrics{}
	metrics.UserID = userID
	for i := 1; i <= 7; i++ {
		var count int
		row := as.Client.QueryRow("select count(*) from articles where user_id=? and where category_id=?", userID, i)
		err := row.Scan(&count)
		if err != nil {
			return nil, err
		}
		category, err := as.getCategoryByID(i)
		if err != nil {
			return nil, err
		}
		metrics.CategoryToNumArticles[category] = count
	}
	return metrics, nil
}

func (as *ArticleStore) InsertArticle(category string, userID int64) error {
	insq := "insert into articles(user_id, category_id, read_on) values (?, ?, ?)"
	categoryID, err := as.getCategoryID(category)
	if err != nil {
		return err
	}
	_, err = as.Client.Exec(insq, userID, categoryID, time.Now())
	if err != nil {
		log.Printf("Issue executing sql statement: %v", err)
		return err
	}
	return nil
}

func (as *ArticleStore) getCategoryByID(categoryID int) (string, error) {
	var category string
	row := as.Client.QueryRow("select category_name from categories where category_id=?", categoryID)
	if err := row.Scan(&category); err != nil {
		return "", err
	}
	return category, nil
}

func (as *ArticleStore) getCategoryID(category string) (int, error) {
	var id int
	row := as.Client.QueryRow("select category_id from categories where category_name=?", category)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}
