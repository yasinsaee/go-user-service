package group

import (
    "fmt"

    "github.com/yasinsaee/go-user-service/internal/context"
    "github.com/yasinsaee/go-user-service/internal/domain/group"
    "github.com/yasinsaee/go-user-service/internal/domain/permission"
    "github.com/yasinsaee/go-user-service/internal/domain/role"
    grouppb "github.com/yasinsaee/go-user-service/user-service/group"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// ====================================================
// ✅ HELPER FUNCTIONS
// ====================================================

func toPermissionPB(p *permission.Permission) *grouppb.Permission {
    if p == nil {
        return nil
    }
    return &grouppb.Permission{
        Id:          p.UniqueID,
        Name:        p.Name,
        Description: p.Description,
    }
}

func (h *Handler) getPermissionsMap(ids []string) (map[string]permission.Permission, error) {
    if len(ids) == 0 {
        return make(map[string]permission.Permission), nil
    }

    perms, err := h.perService.GetByIDs(ids)
    if err != nil {
        return nil, err
    }

    result := make(map[string]permission.Permission, len(perms))
    for _, p := range perms {
        result[p.UniqueID] = p
    }
    return result, nil
}

func (h *Handler) getRolesMap(ids []string) (map[string]role.Role, error) {
    if len(ids) == 0 {
        return make(map[string]role.Role), nil
    }

    roles, err := h.roleService.GetByIDs(ids)
    if err != nil {
        return nil, err
    }

    result := make(map[string]role.Role, len(roles))
    for _, r := range roles {
        result[r.UniqueID] = r
    }
    return result, nil
}

func (h *Handler) toRolePB(id string, rolesMap map[string]role.Role) (*grouppb.Role, error) {
    r, ok := rolesMap[id]
    if !ok {
        return nil, fmt.Errorf("role not found: %s", id)
    }

    permsMap, err := h.getPermissionsMap(r.Permissions)
    if err != nil {
        return nil, err
    }

    permPBs := make([]*grouppb.Permission, 0, len(permsMap))
    for _, p := range permsMap {
        permPBs = append(permPBs, toPermissionPB(&p))
    }

    return &grouppb.Role{
        Id:          r.UniqueID,
        Name:        r.Name,
        Description: r.Description,
        Permissions: permPBs,
    }, nil
}

func (h *Handler) toSelectedRolesPB(selectedRoles group.SelectedRoles) ([]*grouppb.SelectedRoles, error) {
    if len(selectedRoles) == 0 {
        return []*grouppb.SelectedRoles{}, nil
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

    result := make([]*grouppb.SelectedRoles, 0, len(selectedRoles))
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

        // فیلتر کردن دسترسی‌های نهایی
        finalPerms := make([]*grouppb.Permission, 0, len(rolePB.Permissions))
        for _, p := range rolePB.Permissions {
            if !deniedMap[p.Id] {
                finalPerms = append(finalPerms, p)
            }
        }

        result = append(result, &grouppb.SelectedRoles{
            Role:        rolePB,
            Permissions: finalPerms,
        })
    }

    return result, nil
}

func (h *Handler) toGroupPB(r *group.Group) (*grouppb.Group, error) {
    selectedRoles, err := h.toSelectedRolesPB(r.SelectedRoles)
    if err != nil {
        return nil, err
    }

    return &grouppb.Group{
        Id:            r.UniqueID,
        Name:          r.Name,
        Description:   r.Description,
        SelectedRoles: selectedRoles,
    }, nil
}

func (h *Handler) toMetaDataPB(r *context.MetaData) (*grouppb.MetaData, error) {
    return &grouppb.MetaData{
        Limit:       int32(r.Limit),
        TotalCounts: int32(r.TotalCounts),
        TotalPages:  int32(r.TotalPages),
        CurrentPage: int32(r.CurrentPage),
        NextPage:    int32(r.NextPage),
        Sort:        r.Sort,
    }, nil
}

func (h *Handler) checkSelectedRole(g_*group.Group, selectedRoles []*grouppb.RequestSelectedRoles) error {
    processedRoles := make(map[string]*group.SelectedRole)

    for _, selected := range selectedRoles {
        roleID := selected.GetRole()
        if roleID == "" {
            continue
        }

        r, err := h.roleService.GetByID(roleID)
        if err != nil {
            return status.Errorf(codes.InvalidArgument, "role not found: %v", err)
        }

        var deniedPermissions []string

        rolePermsMap := make(map[string]bool)
        for _, pid := range r.Permissions {
            rolePermsMap[pid] = true
        }

        requestAllowedPerms := make(map[string]bool)
        for _, pStr := range selected.Permissions {
            if !rolePermsMap[pStr] {
                return status.Errorf(codes.InvalidArgument, "permission id %s does not belong to role %s", pStr, r.Name)
            }
            requestAllowedPerms[pStr] = true
        }

        for _, pid := range r.Permissions {
            if !requestAllowedPerms[pid] {
                deniedPermissions = append(deniedPermissions, pid)
            }
        }

        processedRoles[roleID] = &group.SelectedRole{
            RoleID:           r.UniqueID,
            DeniedPermissions: deniedPermissions,
        }
    }

    var finalRoles []group.SelectedRole
    for _, role := range processedRoles {
        finalRoles = append(finalRoles, *role)
    }

    g_.SelectedRoles = finalRoles

    return nil
}