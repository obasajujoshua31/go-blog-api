package dal

import (
	"time"
)

type Article struct {
	ID        string    `json:"id"gorm:"varchar(255);not null;unique_index;primary"`
	AuthorID  string    `json:"author_id"gorm:"varchar(255);not null;"`
	Content   string    `json:"content"gorm:"varchar(255)"`
	Title     string    `json:"title"gorm:"varchar(255);not null;"`
	CreatedAt time.Time `json:"created_at"gorm:"varchar(255)"`

	User User `gorm:"foreignkey:author_id"json:"user;omitempty"`
}

func (d *DAL) CreateNewArticle(article Article) (*Article, error) {
	err := d.DB.Create(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (d *DAL) GetArticleById(articleID string) (*Article, error) {
	article := Article{}
	err := d.DB.Preload("User").Where("id = ?", articleID).Take(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (d *DAL) GetArticlesByUserID(userID string) (articles []Article, err error) {

	err = d.DB.Where("author_id = ?", userID).Take(&articles).Error

	if err != nil {
		return nil, err
	}
	return articles, nil
}
