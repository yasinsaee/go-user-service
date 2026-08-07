package group

import (
	"context"

	cContext "github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/context/filter"
	"github.com/yasinsaee/go-user-service/internal/domain/group"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	grouppb "github.com/yasinsaee/go-user-service/user-service/group"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	grouppb.UnimplementedGroupServiceServer
	service     group.GroupService
	roleService role.RoleService
	perService  permission.PermissionService
}

func New(service group.GroupService, roleService role.RoleService, perService permission.PermissionService) *Handler {
	return &Handler{service: service, roleService: roleService, perService: perService}
}

func (h *Handler) CreateGroup(ctx context.Context, req *grouppb.CreateGroupRequest) (*grouppb.GroupResponse, error) {
	g := &group.Group{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		TenantID:    req.GetTenantId(),
	}

	err := h.checkSelectedRole(g, req.GetSelectedRoles())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to selected role error: %v", err)
	}

	err = h.service.Create(g)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create group: %v", err)
	}

	grppb, err := h.toGroupPB(g)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map group: %v", err)
	}

	return &grouppb.GroupResponse{
		Group: grppb,
	}, nil
}

func (h *Handler) UpdateGroup(ctx context.Context, req *grouppb.UpdateGroupRequest) (*grouppb.GroupResponse, error) {
	gp, err := h.service.GetByID(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "group not found: %v", err)
	}

	if name := req.GetName(); name != "" {
		gp.Name = name
	}

	if desc := req.GetDescription(); desc != "" {
		gp.Description = desc
	}

	if teant := req.GetTenantId(); teant != "" {
		gp.TenantID = teant
	}

	if len(req.GetSelectedRoles()) > 0 {
		err := h.checkSelectedRole(gp, req.GetSelectedRoles())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to selected role error: %v", err)
		}
	}

	if err := h.service.Update(gp); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update group: %v", err)
	}

	gppb, err := h.toGroupPB(gp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map group: %v", err)
	}

	return &grouppb.GroupResponse{Group: gppb}, nil
}

func (h *Handler) GetGroup(ctx context.Context, req *grouppb.GetGroupRequest) (*grouppb.GroupResponse, error) {
	r, err := h.service.GetByID(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "group not found: %v", err)
	}

	gppb, err := h.toGroupPB(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map group: %v", err)
	}

	return &grouppb.GroupResponse{Group: gppb}, nil
}

func (h *Handler) ListGroups(ctx context.Context, req *grouppb.ListGroupsRequest) (*grouppb.ListGroupsResponse, error) {
	groups, err := h.service.ListAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list groups: %v", err)
	}

	var pbgroups []*grouppb.Group
	for _, r := range groups {
		grouppb, err := h.toGroupPB(&r)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to map group: %v", err)
		}
		pbgroups = append(pbgroups, grouppb)
	}

	return &grouppb.ListGroupsResponse{Groups: pbgroups}, nil
}

func (h *Handler) ListPaginationGroups(ctx context.Context, req *grouppb.ListPaginationGroupsRequest) (*grouppb.ListPaginationGroupsResponse, error) {
	var metaData = cContext.MetaData{
		Limit:       int(req.GetLimit()),
		Sort:        req.GetSort(),
		CurrentPage: int(req.GetPage()),
	}

	metaDataRes, groups, err := h.service.PaginationList(metaData, group.GroupFilter{SearchFilter: filter.SearchFilter{Search: req.GetSearch()}, TenantID: req.GetTenantId(), IsDelete: "false"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list groups: %v", err)
	}

	var pbgroups []*grouppb.Group
	for _, r := range groups {
		grouppb, err := h.toGroupPB(&r)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to map group: %v", err)
		}
		pbgroups = append(pbgroups, grouppb)
	}

	metaDataPB, err := h.toMetaDataPB(&metaDataRes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map metadata: %v", err)
	}

	return &grouppb.ListPaginationGroupsResponse{Groups: pbgroups, MetaData: metaDataPB}, nil
}

func (h *Handler) DeleteGroup(ctx context.Context, req *grouppb.DeleteGroupRequest) (*grouppb.DeleteGroupResponse, error) {
	err := h.service.Delete(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "group not found: %v", err)
	}

	return &grouppb.DeleteGroupResponse{Message: "ok"}, nil
}

func (h *Handler) GetByFilterGroup(ctx context.Context, req *grouppb.GetByFilterRequest) (*grouppb.GroupResponse, error) {
	g, err := h.service.GetByFilter(group.GroupFilter{IsDelete: "false", Key: req.GetKey(), TenantID: req.GetTeantId()})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "group not found: %v", err)
	}

	gppb, err := h.toGroupPB(g)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map group: %v", err)
	}

	return &grouppb.GroupResponse{Group: gppb}, nil
}
