package repo

import (
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/config"
	"gorm.io/gorm"
)

type WorkspaceRepo interface {
	Create(workspace *model.Workspace) error
	Save(workspace *model.Workspace) error
	ResetApiKey(id string) error
	Delete(id string) error
	GetOneByID(id string) (*model.Workspace, error)
	GetOneByApiKey(apiKey string) (*model.Workspace, error)
	Search(keyword string, page int, limit int) ([]model.Workspace, int64, error)
	GetAll(page int, limit int) ([]model.Workspace, int64, error)
	GetAllByOrgID(orgID string, page int, limit int) ([]model.Workspace, int64, error)
	Total() (int64, error)
}

type workspaceRepo struct {
	db   *gorm.DB
	conf config.Config
}

func NewWorkspaceRepo(db *gorm.DB, conf config.Config) WorkspaceRepo {
	return &workspaceRepo{
		db:   db,
		conf: conf,
	}
}

func (r *workspaceRepo) generateApiKey(workspace *model.Workspace) {
	apikey := shared.Utils.Auth.CreateToken("zky", workspace.ID, r.conf.GetString("app.secret"))
	workspace.ApiKey = shared.Utils.Auth.Encrypt(apikey, r.conf.GetString("app.secret"))
}

func (r *workspaceRepo) decryptApiKey(apiKey string) (string, error) {
	return shared.Utils.Auth.Decrypt(apiKey, r.conf.GetString("app.secret"))
}

func (r *workspaceRepo) encryptApiKey(apiKey string) string {
	return shared.Utils.Auth.Encrypt(apiKey, r.conf.GetString("app.secret"))
}

func (r *workspaceRepo) Create(workspace *model.Workspace) error {
	r.generateApiKey(workspace)
	return r.db.Create(workspace).Error
}

func (r *workspaceRepo) Save(workspace *model.Workspace) error {
	return r.db.Save(workspace).Error
}

func (r *workspaceRepo) ResetApiKey(id string) error {
	workspace := &model.Workspace{ID: id}
	r.generateApiKey(workspace)
	return r.db.Model(workspace).Updates(workspace).Error
}

func (r *workspaceRepo) Delete(id string) error {
	return r.db.Delete(&model.Workspace{}, "id = ?", id).Error
}

func (r *workspaceRepo) GetOneByID(id string) (*model.Workspace, error) {
	var workspace model.Workspace
	if err := r.db.Where("id = ?", id).First(&workspace).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepo) GetOneByApiKey(apiKey string) (*model.Workspace, error) {
	var workspace model.Workspace
	encryptedApiKey := r.encryptApiKey(apiKey)
	if err := r.db.Where("apiKey = ?", encryptedApiKey).First(&workspace).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepo) Search(keyword string, page int, limit int) ([]model.Workspace, int64, error) {
	var workspaces []model.Workspace
	var total int64
	query := r.db.Where("name LIKE ?", "%"+keyword+"%").Or("id LIKE ?", "%"+keyword+"%")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&workspaces).Error; err != nil {
		return nil, 0, err
	}
	return workspaces, total, nil
}

func (r *workspaceRepo) GetAll(page int, limit int) ([]model.Workspace, int64, error) {
	var workspaces []model.Workspace
	var total int64
	if err := r.db.Model(&model.Workspace{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Offset((page - 1) * limit).Limit(limit).Find(&workspaces).Error; err != nil {
		return nil, 0, err
	}
	return workspaces, total, nil
}

func (r *workspaceRepo) GetAllByOrgID(orgID string, page int, limit int) ([]model.Workspace, int64, error) {
	var workspaces []model.Workspace
	var total int64
	if err := r.db.Where("orgId = ?", orgID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Where("orgId = ?", orgID).Offset((page - 1) * limit).Limit(limit).Find(&workspaces).Error; err != nil {
		return nil, 0, err
	}
	return workspaces, total, nil

}

func (r *workspaceRepo) Total() (int64, error) {
	var total int64
	if err := r.db.Model(&model.Workspace{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
