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
	Search(keyword string, page int, limit int) (*shared.List[model.Organization], error)
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

func (uc *orgUseCase) CreateOrg(orgDto CreateOrgStruct) (*model.Organization, error) {
	if err := uc.toolkit.Validator.ValidateStruct(orgDto); err != nil {
		fields := uc.toolkit.Validator.GetValidationErrors(err)
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
	if err := uc.orgRepo.Create(org); err != nil {
		return nil, errs.NewInternalError("failed to create organization", err)
	}
	return org, nil
}

func (uc *orgUseCase) UpdateOrg(id string, orgDto UpdateOrgStruct) (*model.Organization, error) {
	if err := uc.toolkit.Validator.ValidateStruct(orgDto); err != nil {
		fields := uc.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid organization data", fields)
	}
	org, err := uc.orgRepo.GetOneByID(id)
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
	if err := uc.orgRepo.Save(org); err != nil {
		return nil, errs.NewInternalError("failed to update organization", err)
	}
	return org, nil
}

func (uc *orgUseCase) DeleteOrg(id string) error {
	if err := uc.orgRepo.Delete(id); err != nil {
		return errs.NewInternalError("failed to delete organization", err)
	}
	return nil
}

func (uc *orgUseCase) GetOrgByID(id string) (*model.Organization, error) {
	org, err := uc.orgRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("Organization not found", err)
	}
	return org, nil
}

func (uc *orgUseCase) GetOrgByEmail(email string) (*model.Organization, error) {
	org, err := uc.orgRepo.GetOneByEmail(email)
	if err != nil {
		return nil, errs.NewNotFoundError("Organization not found", err)
	}
	return org, nil
}

func (uc *orgUseCase) GetOrgByUserID(userID string) (*model.Organization, error) {
	user, err := uc.userRepo.GetOneByID(userID)
	if err != nil {
		return nil, errs.NewNotFoundError("User not found", err)
	}
	org, err := uc.orgRepo.GetOneByID(user.OrgID)
	if err != nil {
		return nil, errs.NewNotFoundError("Failed to get organization", err)
	}
	return org, nil
}

func (uc *orgUseCase) GetOrgByWorkspaceID(workspaceID string) (*model.Organization, error) {
	workspace, err := uc.workspaceRepo.GetOneByID(workspaceID)
	if err != nil {
		return nil, errs.NewNotFoundError("Workspace not found", err)
	}
	org, err := uc.orgRepo.GetOneByID(workspace.OrgID)
	if err != nil {
		return nil, errs.NewNotFoundError("Failed to get organization", err)
	}
	return org, nil
}

func (uc *orgUseCase) GetAll(page int, limit int) (*shared.List[model.Organization], error) {
	orgs, total, err := uc.orgRepo.GetAll(page, limit)
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

func (uc *orgUseCase) Search(keyword string, page int, limit int) (*shared.List[model.Organization], error) {
	orgs, total, err := uc.orgRepo.Search(keyword, page, limit)
	if err != nil {
		return nil, errs.NewInternalError("Failed to search organizations", err)
	}
	return &shared.List[model.Organization]{
		Items: orgs,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}
