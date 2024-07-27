package usecase

import (
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/service/account/pkg/repo"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/errs"
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
	Name    string  `json:"name,omitempty"`
	Website *string `json:"website,omitempty"`
	Email   string  `json:"email,omitempty" validate:"email"`
	Phone   *string `json:"phone,omitempty" validate:"phone"`
	Country string  `json:"country,omitempty" validate:"iso3166_1_alpha2"`
	City    string  `json:"city,omitempty"`
	Address string  `json:"address,omitempty"`
}

type OrgUseCase interface {
	CreateOrg(orgDto CreateOrgRequest) (*model.Organization, error)
	UpdateOrg(id string, orgDto UpdateOrgRequest) (*model.Organization, error)
	DeleteOrg(id string) error
	GetOrgByID(id string) (*model.Organization, error)
	GetOrgByEmail(email string) (*model.Organization, error)
	GetOrgByUserID(userID string) (*model.Organization, error)
	GetOrgByWorkspaceID(workspaceID string) (*model.Organization, error)
	GetAll(page int, limit int) (*shared.List[model.Organization], error)
}

type orgUseCase struct {
	toolkit       shared.Toolkit
	orgRepo       repo.OrgRepo
	userRepo      repo.UserRepo
	workspaceRepo repo.WorkspaceRepo
}

func (u orgUseCase) CreateOrg(orgDto CreateOrgRequest) (*model.Organization, error) {
	if err := u.toolkit.Validator.ValidateStruct(orgDto); err != nil {
		fields := u.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid organization data", fields)
	}
	org := &model.Organization{
		Name:    orgDto.Name,
		Website: orgDto.Website,
		Email:   orgDto.Email,
		Phone:   orgDto.Phone,
		Country: orgDto.Country,
		City:    orgDto.City,
		Address: orgDto.Address,
	}
	if err := u.orgRepo.Insert(org); err != nil {
		return nil, errs.NewInternalError("failed to create organization", err)
	}
	return org, nil
}

func (u orgUseCase) UpdateOrg(id string, orgDto UpdateOrgRequest) (*model.Organization, error) {
	if err := u.toolkit.Validator.ValidateStruct(orgDto); err != nil {
		fields := u.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid organization data", fields)
	}
	org, err := u.orgRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("Organization not found", err)
	}
	if orgDto.Name != "" {
		org.Name = orgDto.Name
	}
	if orgDto.Website != nil {
		org.Website = orgDto.Website
	}
	if orgDto.Email != "" {
		org.Email = orgDto.Email
	}
	if orgDto.Phone != nil {
		org.Phone = orgDto.Phone
	}
	if orgDto.Country != "" {
		org.Country = orgDto.Country
	}
	if orgDto.City != "" {
		org.City = orgDto.City
	}
	if orgDto.Address != "" {
		org.Address = orgDto.Address
	}
	if err := u.orgRepo.Update(org); err != nil {
		return nil, errs.NewInternalError("failed to update organization", err)
	}
	return org, nil
}

func (u orgUseCase) DeleteOrg(id string) error {
	if err := u.orgRepo.Delete(id); err != nil {
		return errs.NewInternalError("failed to delete organization", err)
	}
	return nil
}

func (u orgUseCase) GetOrgByID(id string) (*model.Organization, error) {
	org, err := u.orgRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("Organization not found", err)
	}
	return org, nil
}

func (u orgUseCase) GetOrgByEmail(email string) (*model.Organization, error) {
	org, err := u.orgRepo.GetOneByEmail(email)
	if err != nil {
		return nil, errs.NewNotFoundError("Organization not found", err)
	}
	return org, nil
}

func (u orgUseCase) GetOrgByUserID(userID string) (*model.Organization, error) {
	user, err := u.userRepo.GetOneByID(userID)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	org, err := u.orgRepo.GetOneByID(user.OrgID)
	if err != nil {
		return nil, errs.NewNotFoundError("Failed to get organization", err)
	}
	return org, nil
}

func (u orgUseCase) GetOrgByWorkspaceID(workspaceID string) (*model.Organization, error) {
	workspace, err := u.workspaceRepo.GetOneByID(workspaceID)
	if err != nil {
		return nil, errs.NewNotFoundError("Workspace not found", err)
	}
	org, err := u.orgRepo.GetOneByID(workspace.OrgID)
	if err != nil {
		return nil, errs.NewNotFoundError("Failed to get organization", err)
	}
	return org, nil
}

func (u orgUseCase) GetAll(page int, limit int) (*shared.List[model.Organization], error) {
	orgs, total, err := u.orgRepo.GetAll(page, limit)
	if err != nil {
		return nil, errs.NewInternalError("Failed to get organizations", err)
	}
	return &shared.List[model.Organization]{
		Items: orgs,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func NewOrgUseCase(toolkit shared.Toolkit, orgRepo repo.OrgRepo, userRepo repo.UserRepo, workspaceRepo repo.WorkspaceRepo) OrgUseCase {
	return &orgUseCase{
		toolkit:       toolkit,
		orgRepo:       orgRepo,
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
	}
}
