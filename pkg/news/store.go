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
	InsertArticle(metric NewMetric, userID int64) error

	//GetIDOfCategory returns the id of the category provided
	getCategoryID(category string) (int, error)

	//GetCategoryByID returns the name of the category provided
	getCategoryByID(categoryID int) (string, error)
}

//ArticleStore represents a SQL implemented databse for News related entries
type ArticleStore struct {
	Client *sql.DB
}

//NewArticleStore constructs a new ArticleStore
func NewArticleStore(db *sql.DB) *ArticleStore {
	//initialize and return a new MySQLStore struct
	if db != nil {
		return &ArticleStore{
			Client: db,
		}
	}
	return nil
}

func (as *ArticleStore) GetByUserID(userID int64) (*Metrics, error) {
	metrics := &Metrics{}
	metrics.UserID = userID
	rows, err := as.Client.Query("select category_name, count(*) from articles inner join categories on articles.category_id=categories.category_id where user_id=? group by category_name order by 2 desc", userID)
	if err != nil {
		log.Print("Error querying for categories count")
		return nil, err
	}
	metrics.CategoryToNumArticles = map[string]int{}
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			log.Print("Error scanning categories")
			return nil, err
		}
		metrics.CategoryToNumArticles[category] = count
	}
	// for i := 1; i <= 7; i++ {
	// 	var count int
	// 	row := as.Client.QueryRow("select count(*) from articles where user_id=? and category_id=?", userID, i)
	// 	err := row.Scan(&count)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	category, err := as.getCategoryByID(i)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	metrics.CategoryToNumArticles[category] = count
	// }

	rows, err = as.Client.Query("select source_name, count(*) from articles inner join sources on articles.source_id=sources.source_id where user_id=? group by source_name order by 2 desc", userID)
	if err != nil {
		log.Print("Error querying for sources count")
		return nil, err
	}
	metrics.SourceToNumArticles = map[string]int{}
	for rows.Next() {
		var sourceName string
		var count int
		if err := rows.Scan(&sourceName, &count); err != nil {
			log.Print("Error scanning sources ")
			return nil, err
		}
		metrics.SourceToNumArticles[sourceName] = count
	}
	return metrics, nil
}

func (as *ArticleStore) InsertArticle(metric NewMetric, userID int64) error {
	insq := "insert into articles(user_id, category_id, source_id, read_on) values (?, ?, ?, ?)"
	categoryID, err := as.getCategoryID(metric.Category)
	if err != nil {
		return err
	}
	sourceID, err := as.getSourceID(metric.Source)
	if err != nil {
		sourceID, err = as.insertSource(metric.Source)
		if err != nil {
			return err
		}
	}
	_, err = as.Client.Exec(insq, userID, categoryID, sourceID, time.Now())
	if err != nil {
		log.Printf("Issue executing sql statement: %v", err)
		return err
	}
	return nil
}

func (as *ArticleStore) insertSource(sourceName string) (int, error) {
	insq := "insert into sources(source_name) values (?)"
	res, err := as.Client.Exec(insq, sourceName)
	if err != nil {
		log.Printf("Issue executing sql statement: %v", err)
		return -1, err
	}

	id, err := res.LastInsertId()
	return int(id), err
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

func (as *ArticleStore) getSourceByID(sourceID int) (string, error) {
	var source string
	row := as.Client.QueryRow("select source_name from sources where source_id=?", sourceID)
	if err := row.Scan(&source); err != nil {
		return "", err
	}
	return source, nil
}

func (as *ArticleStore) getSourceID(source string) (int, error) {
	var id int
	row := as.Client.QueryRow("select source_id from sources where source_name=?", source)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}
