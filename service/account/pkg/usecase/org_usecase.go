package usecase

import (
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/service/account/pkg/repo"
	"github.com/abdelrahman146/zard/shared"
)

type CreateOrgRequest struct {
	Name    string  `json:"name,omitempty" validate:"required,omitempty"`
	Website *string `json:"website,omitempty"`
	Email   string  `json:"email,omitempty" validate:"required,email"`
	Phone   *string `json:"phone,omitempty" validate:"phone"`
	Country string  `json:"country,omitempty" validate:"required,omitempty,iso3166_1_alpha2"`
	City    string  `json:"city,omitempty" validate:"required,omitempty"`
	Address string  `json:"address,omitempty" validate:"required,omitempty"`
}

type UpdateOrgRequest struct {
	ID      string  `json:"id,omitempty" validate:"required"`
	Name    string  `json:"name,omitempty" validate:"required,omitempty"`
	Website *string `json:"website,omitempty"`
	Email   string  `json:"email,omitempty" validate:"required,email"`
	Phone   *string `json:"phone,omitempty" validate:"phone"`
	Country string  `json:"country,omitempty" validate:"required,omitempty,iso3166_1_alpha2"`
	City    string  `json:"city,omitempty" validate:"required,omitempty"`
	Address string  `json:"address,omitempty" validate:"required,omitempty"`
}

type OrgUseCase interface {
	CreateOrg(orgDto CreateOrgRequest) (*model.Organization, error)
	UpdateOrg(orgDto UpdateOrgRequest) (*model.Organization, error)
	DeleteOrg(id string) error
	GetOrgByID(id string) (*model.Organization, error)
	GetOrgByEmail(email string) (*model.Organization, error)
	GetOrgByUserID(userID string) (*model.Organization, error)
	GetOrgByWorkspaceID(workspaceID string) (*model.Organization, error)
	GetAll(page int, limit int) ([]model.Organization, int64, error)
}

type orgUseCase struct {
	toolkit       shared.Toolkit
	orgRepo       repo.OrgRepo
	userRepo      repo.UserRepo
	workspaceRepo repo.WorkspaceRepo
}

func (u orgUseCase) CreateOrg(orgDto CreateOrgRequest) (*model.Organization, error) {
	if err := u.toolkit.Validator.ValidateStruct(orgDto); err != nil {
		validationErrors := u.toolkit.Validator.GetValidationErrors(err)
	}
}

func (u orgUseCase) UpdateOrg(orgDto UpdateOrgRequest) (*model.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (u orgUseCase) DeleteOrg(id string) error {
	//TODO implement me
	panic("implement me")
}

func (u orgUseCase) GetOrgByID(id string) (*model.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (u orgUseCase) GetOrgByEmail(email string) (*model.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (u orgUseCase) GetOrgByUserID(userID string) (*model.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (u orgUseCase) GetOrgByWorkspaceID(workspaceID string) (*model.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (u orgUseCase) GetAll(page int, limit int) ([]model.Organization, int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewOrgUseCase(toolkit shared.Toolkit, orgRepo repo.OrgRepo, userRepo repo.UserRepo, workspaceRepo repo.WorkspaceRepo) OrgUseCase {
	return &orgUseCase{
		toolkit:       toolkit,
		orgRepo:       orgRepo,
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
	}
}
