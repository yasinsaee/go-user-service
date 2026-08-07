package usergrpc

import (
	"fmt"

	"github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/domain/group"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"github.com/yasinsaee/go-user-service/pkg/jwt"
	userpb "github.com/yasinsaee/go-user-service/user-service/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ====================================================
// ✅ HELPER FUNCTIONS
// ====================================================

// ====================================================
// ✅ HELPER FUNCTIONS (OPTIMIZED FOR HIGH SCALE)
// ====================================================

func toPermissionPB(p *permission.Permission) *userpb.Permission {
	if p == nil {
		return nil
	}
	return &userpb.Permission{
		Id:          p.UniqueID,
		Name:        p.Name,
		Description: p.Description,
		Key:         p.Key,
	}
}

func (h *Handler) getPermissionsMap(ids []string) (map[string]permission.Permission, error) {
	if len(ids) == 0 {
		return make(map[string]permission.Permission), nil
	}

	result := make(map[string]permission.Permission)
	missingIDs := make([]string, 0, len(ids))

	h.cacheMutex.RLock()
	for _, id := range ids {
		if p, ok := h.permCache[id]; ok {
			result[id] = *p
		} else {
			missingIDs = append(missingIDs, id)
		}
	}
	h.cacheMutex.RUnlock()

	if len(missingIDs) > 0 {
		perms, err := h.pService.GetByIDs(missingIDs)
		if err != nil {
			return nil, err
		}

		h.cacheMutex.Lock()
		for _, p := range perms {
			result[p.UniqueID] = p
			if h.permCache == nil {
				h.permCache = make(map[string]*permission.Permission)
			}
			h.permCache[p.UniqueID] = &p
		}
		h.cacheMutex.Unlock()
	}

	return result, nil
}

func (h *Handler) getRolesMap(ids []string) (map[string]role.Role, error) {
	if len(ids) == 0 {
		return make(map[string]role.Role), nil
	}

	result := make(map[string]role.Role)
	missingIDs := make([]string, 0, len(ids))

	h.cacheMutex.RLock()
	for _, id := range ids {
		if r, ok := h.roleCache[id]; ok {
			result[id] = *r
		} else {
			missingIDs = append(missingIDs, id)
		}
	}
	h.cacheMutex.RUnlock()

	if len(missingIDs) > 0 {
		roles, err := h.rService.GetByIDs(missingIDs)
		if err != nil {
			return nil, err
		}

		h.cacheMutex.Lock()
		for _, r := range roles {
			result[r.UniqueID] = r
			if h.roleCache == nil {
				h.roleCache = make(map[string]*role.Role)
			}
			h.roleCache[r.UniqueID] = &r
		}
		h.cacheMutex.Unlock()
	}

	return result, nil
}

func (h *Handler) toRolePB(id string, rolesMap map[string]role.Role) (*userpb.Role, error) {
	r, ok := rolesMap[id]
	if !ok {
		return nil, fmt.Errorf("role not found in map: %s", id)
	}

	perms, err := h.getPermissionsMap(r.Permissions)
	if err != nil {
		return nil, err
	}

	permPBs := make([]*userpb.Permission, 0, len(perms))
	for _, p := range perms {
		permPBs = append(permPBs, toPermissionPB(&p))
	}

	return &userpb.Role{
		Id:          r.UniqueID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: permPBs,
		Key:         r.Key,
	}, nil
}

func (h *Handler) toSelectedRolesPB(selectedRoles group.SelectedRoles) ([]*userpb.SelectedRoles, error) {
	if len(selectedRoles) == 0 {
		return []*userpb.SelectedRoles{}, nil
	}

	roleIDs := make([]string, 0, len(selectedRoles))
	for _, sr := range selectedRoles {
		roleIDs = append(roleIDs, sr.RoleID)
	}

	rolesMap, err := h.getRolesMap(roleIDs)
	if err != nil {
		return nil, err
	}

	deniedPermIDs := make([]string, 0)
	for _, sr := range selectedRoles {
		deniedPermIDs = append(deniedPermIDs, sr.DeniedPermissions...)
	}

	deniedPermsMap, err := h.getPermissionsMap(deniedPermIDs)
	if err != nil {
		return nil, err
	}

	result := make([]*userpb.SelectedRoles, 0, len(selectedRoles))
	for _, sr := range selectedRoles {
		rolePB, err := h.toRolePB(sr.RoleID, rolesMap)
		if err != nil {
			return nil, err
		}

		deniedMap := make(map[string]bool)
		for _, deniedID := range sr.DeniedPermissions {
			if p, ok := deniedPermsMap[deniedID]; ok {
				deniedMap[p.UniqueID] = true
			}
		}

		finalPerms := make([]*userpb.Permission, 0, len(rolePB.Permissions))
		for _, p := range rolePB.Permissions {
			if !deniedMap[p.Id] {
				finalPerms = append(finalPerms, p)
			}
		}

		result = append(result, &userpb.SelectedRoles{
			Role:        rolePB,
			Permissions: finalPerms,
		})
	}

	return result, nil
}

func (h *Handler) toGroupPB(id string) (*userpb.Group, error) {
	g, err := h.gService.GetByID(id)
	if err != nil {
		return nil, err
	}

	selectedRoles, err := h.toSelectedRolesPB(g.SelectedRoles)
	if err != nil {
		return nil, err
	}

	return &userpb.Group{
		Id:            g.UniqueID,
		Name:          g.Name,
		Description:   g.Description,
		SelectedRoles: selectedRoles,
	}, nil
}

func (h *Handler) toUserPB(u *user.User, showID bool) *userpb.User {
	var groupPB *userpb.Group
	if u.Group != "" {
		groupPB, _ = h.toGroupPB(u.Group)
	}

	if !showID {
		u.UniqueID = ""
	}

	return &userpb.User{
		Id:           u.UniqueID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		ProfileImage: u.ProfileImage,
		Group:        groupPB,
		Username:     u.Username,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		IsActive:     u.IsActive,
		IsBanned:     u.IsBanned,
		BannedAt:     timestamppb.New(u.BannedAt),
		CreatedAt:    timestamppb.New(u.CreatedAt),
		UpdatedAt:    timestamppb.New(u.UpdatedAt),
		LastLogin:    timestamppb.New(u.LastLogin),
		TenantId:     u.TenantID,
	}
}

func (h *Handler) toUserJwtMeta(u *user.User) jwt.UserAccesses {
	userAccesses := make(jwt.UserAccesses, 0)

	g, err := h.gService.GetByID(u.Group)
	if err != nil || g == nil {
		return userAccesses
	}

	if len(g.SelectedRoles) == 0 {
		return userAccesses
	}

	roleIDs := make([]string, 0, len(g.SelectedRoles))
	for _, sr := range g.SelectedRoles {
		roleIDs = append(roleIDs, sr.RoleID)
	}
	rolesMap, err := h.getRolesMap(roleIDs)
	if err != nil {
		return userAccesses
	}

	deniedPermIDs := make([]string, 0)
	for _, sr := range g.SelectedRoles {
		deniedPermIDs = append(deniedPermIDs, sr.DeniedPermissions...)
	}
	deniedPermsMap, err := h.getPermissionsMap(deniedPermIDs)
	if err != nil {
		return userAccesses
	}

	for _, sr := range g.SelectedRoles {
		role, ok := rolesMap[sr.RoleID]
		if !ok {
			continue
		}

		userAccess := jwt.UserAccess{
			RoleKey: role.Key,
		}

		deniedKeysMap := make(map[string]bool)
		for _, deniedID := range sr.DeniedPermissions {
			if p, ok := deniedPermsMap[deniedID]; ok {
				deniedKeysMap[p.Key] = true
			}
		}

		rolePermsMap, err := h.getPermissionsMap(role.Permissions)
		if err != nil {
			continue
		}

		for _, permID := range role.Permissions {
			p, ok := rolePermsMap[permID]
			if !ok {
				continue
			}

			if !deniedKeysMap[p.Key] {
				userAccess.Accesses = append(userAccess.Accesses, p.Key)
			}
		}

		userAccesses = append(userAccesses, userAccess)
	}
	return userAccesses
}

func (h *Handler) toMetaDataPB(r *context.MetaData) (*userpb.MetaData, error) {
	return &userpb.MetaData{
		Limit:       int32(r.Limit),
		TotalCounts: int32(r.TotalCounts),
		TotalPages:  int32(r.TotalPages),
		CurrentPage: int32(r.CurrentPage),
		NextPage:    int32(r.NextPage),
		Sort:        r.Sort,
	}, nil
}
