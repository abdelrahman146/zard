package usecase

import (
	"github.com/abdelrahman146/zard/service/account/pkg/model"
	"github.com/abdelrahman146/zard/service/account/pkg/repo"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/errs"
)

type WorkspaceUseCase interface {
	CreateWorkSpace(wsDto *CreateWorkspaceStruct) (*model.Workspace, error)
	UpdateWorkSpace(id string, wsDto *UpdateWorkspaceStruct) (*model.Workspace, error)
	ResetApiKey(id string) (apikey string, err error)
	DeleteWorkSpace(id string) error
	GetWorkSpaceByID(id string) (*model.Workspace, error)
	GetWorkSpaceByApiKey(apiKey string) (*model.Workspace, error)
	GetAll(page int, limit int) (*shared.List[model.Workspace], error)
	GetAllByOrgID(orgID string, page int, limit int) (*shared.List[model.Workspace], error)
	Search(keyword string, page int, limit int) (*shared.List[model.Workspace], error)
}

func NewWorkspaceUseCase(toolkit shared.Toolkit, wsRepo repo.WorkspaceRepo) WorkspaceUseCase {
	return &wsUseCase{
		toolkit: toolkit,
		wsRepo:  wsRepo,
	}
}

type wsUseCase struct {
	toolkit shared.Toolkit
	wsRepo  repo.WorkspaceRepo
}

func (uc *wsUseCase) CreateWorkSpace(wsDto *CreateWorkspaceStruct) (*model.Workspace, error) {
	if err := uc.toolkit.Validator.ValidateStruct(wsDto); err != nil {
		fields := uc.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid workspace data", fields)
	}
	ws := &model.Workspace{
		Name:    wsDto.Name,
		Website: wsDto.Website,
		OrgID:   wsDto.OrgID,
	}
	if err := uc.wsRepo.Create(ws); err != nil {
		return nil, errs.NewInternalError("failed to create workspace", err)
	}
	return ws, nil
}

func (uc *wsUseCase) UpdateWorkSpace(id string, wsDto *UpdateWorkspaceStruct) (*model.Workspace, error) {
	if err := uc.toolkit.Validator.ValidateStruct(wsDto); err != nil {
		fields := uc.toolkit.Validator.GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid workspace data", fields)
	}
	ws, err := uc.wsRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("workspace not found", err)
	}
	if wsDto.Website != nil {
		ws.Website = wsDto.Website
	}
	if wsDto.Name != nil {
		ws.Name = *wsDto.Name
	}
	if err := uc.wsRepo.Save(ws); err != nil {
		return nil, errs.NewInternalError("failed to update workspace", err)
	}
	return ws, nil
}

func (uc *wsUseCase) ResetApiKey(id string) (apikey string, err error) {
	ws, err := uc.wsRepo.ResetApiKey(id)
	if err != nil {
		return "", errs.NewInternalError("failed to reset api key", err)
	}
	return ws.ApiKey, nil
}

func (uc *wsUseCase) DeleteWorkSpace(id string) error {
	_, err := uc.wsRepo.GetOneByID(id)
	if err != nil {
		return errs.NewNotFoundError("workspace not found", err)
	}
	return uc.wsRepo.Delete(id)
}

func (uc *wsUseCase) GetWorkSpaceByID(id string) (*model.Workspace, error) {
	ws, err := uc.wsRepo.GetOneByID(id)
	if err != nil {
		return nil, errs.NewNotFoundError("workspace not found", err)
	}
	return ws, nil
}

func (uc *wsUseCase) GetWorkSpaceByApiKey(apiKey string) (*model.Workspace, error) {
	ws, err := uc.wsRepo.GetOneByApiKey(apiKey)
	if err != nil {
		return nil, errs.NewNotFoundError("workspace not found", err)
	}
	return ws, nil
}

func (uc *wsUseCase) GetAll(page int, limit int) (*shared.List[model.Workspace], error) {
	workspaces, total, err := uc.wsRepo.GetAll(page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to get workspaces", err)
	}
	return &shared.List[model.Workspace]{
		Items: workspaces,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (uc *wsUseCase) GetAllByOrgID(orgID string, page int, limit int) (*shared.List[model.Workspace], error) {
	workspaces, total, err := uc.wsRepo.GetAllByOrgID(orgID, page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to get workspaces", err)
	}
	return &shared.List[model.Workspace]{
		Items: workspaces,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (uc *wsUseCase) Search(keyword string, page int, limit int) (*shared.List[model.Workspace], error) {
	workspaces, total, err := uc.wsRepo.Search(keyword, page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to search workspaces", err)
	}
	return &shared.List[model.Workspace]{
		Items: workspaces,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
