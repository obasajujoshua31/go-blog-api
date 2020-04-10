package dal

import (
	"time"
)

type Comment struct {
	ID string `json:"id,omitempty"gorm:"varchar(255);not null;unique_index;primary"`
	Content string `json:"content,omitempty"gorm:"varchar(255);not null"`
	ReviewerID string `json:"reviewer_id,omitempty"gorm:"varchar(255);not null"`
	ArticleID string `json:"article_id,omitempty"gorm:"varchar(255);not null"`
	CreatedAt time.Time `json:"created_at,omitempty"gorm:"varchar(255);not null"`
	NoOfLikes int `json:"no_of_likes"gorm:"integer;default:0"`

	Likes []Like `json:"likes"gorm:"foreignkey:comment_id"`
	User User `gorm:"foreignkey:reviewer_id"json:"user,omitempty"`
}

func (d *DAL) CreateComment(comment Comment, article Article) (*Comment, error) {
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
	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&art).Where("id = ?", article.ID).Update("no_of_comments", article.NoOfComments + 1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &comment, tx.Commit().Error

}

func (d *DAL) GetComments() (comments []Comment, err error) {
	err = d.DB.Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (d *DAL) GetCommentByID(commentID string) (comment *Comment, err error) {
	comm := Comment{}
	err = d.DB.Preload("User").Where("id = ?", commentID).Take(&comm).Error
	if err != nil {
		return nil, err
	}

	return &comm, nil
}

func (d *DAL) UpdateComment(comment Comment) (*Comment, error) {

	err := d.DB.Model(&Comment{}).Where("id = ?", comment.ID).Update("content", comment.Content).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (d *DAL) DeleteComment(commentID string, article Article) error {
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
	if err := tx.Where("id = ?", commentID).Delete(&Comment{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&art).Where("id = ?", article.ID).Update("no_of_comments", article.NoOfComments - 1).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (d *DAL) GetCommentOnArticle(articleID string) (comments []Comment, err error) {
	err = d.DB.Where("article_id = ?", articleID).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (d *DAL) GetCommentByUserID(userID string) (comments []Comment, err error) {
	err = d.DB.Where("reviewer_id = ?", userID).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil

}