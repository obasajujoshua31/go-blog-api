package dal

import (
	"go-blog-api/services"
	"time"
)

const (
	ArticleReaction SourceType = "article"
	CommentReaction SourceType = "comment"
)

type Like struct {
	ID string `json:"id,omitempty"gorm:"varchar(255);not null;unique_index;primary"`
	ArticleID string `json:"article_id,omitempty"gorm:"varchar(255)"`
	ReactorID string `json:"reactor_id,omitempty"gorm:"varchar(255)"`
	CommentID string `json:"comment_id,omitemtpy"gorm:"varchar(255)"`
	CreatedAt time.Time `json:"created_at,omitempty"gorm:"varchar(255);not null"`
	SourceType  SourceType `json:"source_type"gorm:"varchar(100);not null"`

	User User `gorm:"foreignkey:reviewer_id"json:"user,omitempty"`
}

type SourceType string


func (d *DAL) GetLikeForArticle(userID, articleID string) (*Like, error) {
	like := Like{}

	err := d.DB.Where("reactor_id = ? AND article_id = ?", userID, articleID).Take(&like).Error
	if err != nil {
		return nil, err
	}

	return &like, nil
}

func (d *DAL) GetLikeForComment(userID, commentID string) (*Like, error) {
	like := Like{}

	err := d.DB.Where("reactor_id = ? AND comment_id = ?", userID, commentID).Take(&like).Error
	if err != nil {
		return nil, err
	}

	return &like, nil
}

func (d *DAL) AddLikeToArticle(userID string, article Article) (*Like, error) {
	like := Like{
		ID:         services.GenerateUUID(),
		ArticleID:  article.ID,
		ReactorID:  userID,
		CreatedAt:  time.Now(),
		SourceType: ArticleReaction,
	}

	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	art := Article{}
	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&art).Where("id = ?", article.ID).Update("no_of_likes", article.NoOfLikes + 1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &like, tx.Commit().Error
}

func (d *DAL) AddLikeToComment(userID string, comment Comment) (*Like, error) {
	like := Like{
		ID:         services.GenerateUUID(),
		CommentID:  comment.ID,
		ReactorID:  userID,
		CreatedAt:  time.Now(),
		SourceType: CommentReaction,
	}

	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	comm := Comment{}
	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&comm).Where("id = ?", comment.ID).Update("no_of_likes", comment.NoOfLikes + 1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &like, tx.Commit().Error
}

func (d *DAL) RemoveLikeFromArticle(userID string, article Article) error {

	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return  err
	}
	art := Article{}
	if err := tx.Where("reactor_id = ? AND article_id = ?", userID, article.ID).Delete(&Like{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&art).Where("id = ?", article.ID).Update("no_of_likes", article.NoOfLikes - 1).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (d *DAL) RemoveLikeFromComment(userID string, comment Comment) error {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return  err
	}

	comm := Comment{}
	if err := tx.Where("reactor_id = ? AND comment_id = ?", userID, comment.ID).Delete(&Like{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&comm).Where("id = ?", comment.ID).Update("no_of_likes", comment.NoOfLikes - 1).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}