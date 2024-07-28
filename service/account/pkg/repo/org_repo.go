package repo

import (
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/shared/config"
	"gorm.io/gorm"
)

type OrgRepo interface {
	Create(org *model.Organization) error
	Save(org *model.Organization) error
	Delete(id string) error
	GetOneByID(id string) (*model.Organization, error)
	GetOneByName(name string) (*model.Organization, error)
	GetOneByEmail(email string) (*model.Organization, error)
	Search(keyword string, page int, limit int) ([]model.Organization, int64, error)
	GetAll(page int, limit int) ([]model.Organization, int64, error)
	Total() (int64, error)
}

type orgRepo struct {
	db   *gorm.DB
	conf config.Config
}

func NewOrgRepo(db *gorm.DB, conf config.Config) OrgRepo {
	return &orgRepo{
		db:   db,
		conf: conf,
	}
}

func (r *orgRepo) Create(org *model.Organization) error {
	return r.db.Create(org).Error
}

func (r *orgRepo) Save(org *model.Organization) error {
	return r.db.Save(org).Error
}

func (r *orgRepo) Delete(id string) error {
	return r.db.Delete(&model.Organization{}, "id = ?", id).Error
}

func (r *orgRepo) GetOneByID(id string) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.Where("id = ?", id).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *orgRepo) GetOneByName(name string) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.Where("name = ?", name).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *orgRepo) GetOneByEmail(email string) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.Where("email = ?", email).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *orgRepo) Search(keyword string, page int, limit int) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64
	query := r.db.Model(&model.Organization{}).Where("name LIKE ?", "%"+keyword+"%").Or("email LIKE ?", "%"+keyword+"%").Or("phone LIKE ?", "%"+keyword+"%").Or("id LIKE ?", "%"+keyword+"%")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&orgs).Error; err != nil {
		return nil, 0, err
	}
	return orgs, total, nil
}

func (r *orgRepo) GetAll(page int, limit int) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64
	if err := r.db.Model(&model.Organization{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Offset((page - 1) * limit).Limit(limit).Find(&orgs).Error; err != nil {
		return nil, 0, err
	}
	return orgs, total, nil
}

func (r *orgRepo) Total() (int64, error) {
	var total int64
	if err := r.db.Model(&model.Organization{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
