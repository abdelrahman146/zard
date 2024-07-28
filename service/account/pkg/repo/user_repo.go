package repo

import (
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/shared/cache"
	"github.com/abdelrahman146/zard/shared/config"
	"github.com/abdelrahman146/zard/shared/utils"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(user *model.User) error
	Save(user *model.User) error
	UpdatePassword(id string, password string) error
	Delete(id string) error
	GetOneByID(id string) (*model.User, error)
	GetOneByEmail(email string) (*model.User, error)
	GetOneByPhone(phone string) (*model.User, error)
	Search(keyword string, page int, limit int) ([]model.User, int64, error)
	GetAll(page int, limit int) ([]model.User, int64, error)
	GetAllByOrgID(orgID string, page int, limit int) ([]model.User, int64, error)
	Total() (int64, error)
}

type userRepo struct {
	db          *gorm.DB
	cacheClient cache.Cache
	conf        config.Config
}

func NewUserRepo(db *gorm.DB, cacheClient cache.Cache, conf config.Config) UserRepo {
	return &userRepo{
		db:          db,
		cacheClient: cacheClient,
		conf:        conf,
	}
}

func (r *userRepo) hashPassword(password *string) *string {
	if password != nil {
		hashedPassword := utils.Utils.Auth.Encrypt(*password, r.conf.GetString("app.secret"))
		password = &hashedPassword
	}
	return password
}

func (r *userRepo) Create(user *model.User) error {
	user.Password = r.hashPassword(user.Password)
	return r.db.Create(user).Error
}

func (r *userRepo) Save(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) UpdatePassword(id string, password string) error {
	hashedPassword := utils.Utils.Auth.Encrypt(password, r.conf.GetString("app.secret"))
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func (r *userRepo) Delete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *userRepo) GetOneByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetOneByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetOneByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Search(keyword string, page int, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.Model(&model.User{}).Where("name LIKE ?", "%"+keyword+"%").Or("email LIKE ?", "%"+keyword+"%").Or("phone LIKE ?", "%"+keyword+"%")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) GetAll(page int, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) GetAllByOrgID(orgID string, page int, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	if err := r.db.Model(&model.User{}).Where("orgId = ?", orgID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Where("orgId = ?", orgID).Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) Total() (int64, error) {
	var total int64
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
