package dal

import "go-blog-api/services"

type User struct {
	ID       string `json:"id,omitempty"gorm:"type:varchar(200);unique_index;not null;primary"`
	Name     string `json:"name,omitempty"gorm:"type:varchar(200);not null"`
	Email    string `json:"email,omitempty"gorm:"type:varchar(200);not null"`
	Password string `json:"password,omitempty"gorm:"type:varchar(200);not null"`

	Articles []Article `gorm:"foreignkey:author_id"json:",omitempty"`
}

type DAL struct {
	DB *services.DB
}

type IDAL interface {
}

func NewDAL(db *services.DB) *DAL {
	return &DAL{
		DB: db,
	}
}

func (d *DAL) GetUserById(userID string) (*User, error) {
	user := &User{}
	err := d.DB.Where("id = ?", userID).Take(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DAL) CreateNewUser(user User) (*User, error) {
	userPass, err := services.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = userPass
	err = d.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *DAL) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := d.DB.Where("email = ?", email).Take(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
