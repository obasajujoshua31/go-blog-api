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
	NoOfLikes int `json:"no_of_likes"gorm:"integer;default:0"`
	NoOfComments int `json:"no_of_comments"gorm:"integer;default:0"`

	User User `gorm:"foreignkey:author_id"json:"user;omitempty"`
	Comments []Comment `gorm:"foreignkey:article_id"json:"comments"`
	Likes []Like `json:"likes"gorm:"foreignkey:article_id"`
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
	err := d.DB.Preload("User").Preload("Comments").Preload("Likes").Where("id = ?", articleID).Take(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (d *DAL) GetArticlesByUserID(userID string) (articles []Article, err error) {

	err = d.DB.Where("author_id = ?", userID).Find(&articles).Error

	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (d *DAL) GetAllArticles() (articles []Article, err error) {

	err = d.DB.Find(&articles).Error

	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (d *DAL) UpdateArticle(article Article) (*Article, error) {

	err := d.DB.Model(&Article{}).Where("id = ?", article.ID).Updates(article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (d *DAL) DeleteArticle(articleID string) error {
	err := d.DB.Where("id = ?", articleID).Delete(&Article{}).Error

	if err != nil {
		return err
	}

	return nil
}
