package usecase

import (
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/service/account/pkg/repo"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/errs"
)

type OrgUseCase interface {
	CreateOrg(orgDto CreateOrgStruct) (*model.Organization, error)
	UpdateOrg(id string, orgDto UpdateOrgStruct) (*model.Organization, error)
	DeleteOrg(id string) error
	GetOrgByID(id string) (*model.Organization, error)
	GetOrgByEmail(email string) (*model.Organization, error)
	GetOrgByUserID(userID string) (*model.Organization, error)
	GetOrgByWorkspaceID(workspaceID string) (*model.Organization, error)
	GetAll(page int, limit int) (*shared.List[model.Organization], error)
}

func NewOrgUseCase(toolkit shared.Toolkit, orgRepo repo.OrgRepo, userRepo repo.UserRepo, workspaceRepo repo.WorkspaceRepo) OrgUseCase {
	return &orgUseCase{
		toolkit:       toolkit,
		orgRepo:       orgRepo,
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
	}
}

type orgUseCase struct {
	toolkit       shared.Toolkit
	orgRepo       repo.OrgRepo
	userRepo      repo.UserRepo
	workspaceRepo repo.WorkspaceRepo
}

func (u orgUseCase) CreateOrg(orgDto CreateOrgStruct) (*model.Organization, error) {
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

func (u orgUseCase) UpdateOrg(id string, orgDto UpdateOrgStruct) (*model.Organization, error) {
	if err := u.toolkit.Validator.ValidateStruct(orgDto); err != nil {
		fields := u.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid organization data", fields)
	}
	org, err := u.orgRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("Organization not found", err)
	}
	if orgDto.Name != nil {
		org.Name = *orgDto.Name
	}
	if orgDto.Website != nil {
		org.Website = orgDto.Website
	}
	if orgDto.Email != nil {
		org.Email = *orgDto.Email
	}
	if orgDto.Phone != nil {
		org.Phone = orgDto.Phone
	}
	if orgDto.Country != nil {
		org.Country = *orgDto.Country
	}
	if orgDto.City != nil {
		org.City = *orgDto.City
	}
	if orgDto.Address != nil {
		org.Address = *orgDto.Address
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
