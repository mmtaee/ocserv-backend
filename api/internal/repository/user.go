package repository

import (
	"api/pkg/utils"
	"context"
	"errors"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserRepository struct {
	Admin AdminRepositoryInterface
	Staff StaffRepositoryInterface
	db    *gorm.DB
}

type UserRepositoryInterface interface {
	Login(c context.Context, username, password string, rememberMe bool) (string, error)
	Logout(context.Context) error
	ChangePassword(c context.Context, oldPassword, newPassword string) error
	CreateToken(c context.Context, id uint, expireAt time.Time) (string, error)
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Admin: NewAdminRepository(),
		Staff: NewStaffRepository(),
		db:    database.Connection(),
	}
}

func (r *UserRepository) Login(c context.Context, username, passwd string, rememberMe bool) (string, error) {
	var (
		user     models.User
		expireAt time.Time
	)
	err := r.db.WithContext(c).Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}

	if !utils.Check(passwd, user.Password, user.Salt) {
		return "", errors.New("invalid username and password")
	}
	if rememberMe {
		expireAt = time.Now().Add(time.Hour * 24 * 30)
	} else {
		expireAt = time.Now().Add(time.Hour * 24)
	}
	token, err := r.CreateToken(c, user.ID, expireAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *UserRepository) Logout(c context.Context) error {
	userID := c.Value("userID")
	token := c.Value("token")
	return r.db.WithContext(c).
		Where("token = ? AND user_id = ? ", token, userID).
		Delete(&models.UserToken{}).Error
}

func (r *UserRepository) ChangePassword(c context.Context, oldPasswd, newPasswd string) error {
	var user models.User
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, c.Value("userID")).Error; err != nil {
			return err
		}
		if !utils.Check(oldPasswd, user.Password, user.Salt) {
			return errors.New("incorrect old password")
		}
		pass := utils.NewPassword(newPasswd)
		user.Password = pass.Hash
		user.Salt = pass.Salt
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *UserRepository) CreateToken(c context.Context, id uint, expireAt time.Time) (string, error) {
	token := models.UserToken{
		UserID:   id,
		Token:    utils.Create(id, expireAt),
		ExpireAt: &expireAt,
	}
	err := r.db.WithContext(c).Create(&token).Error
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
