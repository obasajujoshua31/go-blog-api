package dal

import "go-blog-api/services"

type User struct {
	ID       string `json:"id"gorm:"type:varchar(200);unique_index;not null;primary"`
	Name     string `json:"name"gorm:"type:varchar(200);not null"`
	Email    string `json:"email"gorm:"type:varchar(200);not null"`
	Password string `json:"password"gorm:"type:varchar(200);not null"`
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

func (d *DAL) CreateNewRecord(user User) (*User, error) {
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
