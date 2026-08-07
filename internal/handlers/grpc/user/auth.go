package usergrpc

import (
	"context"

	"github.com/yasinsaee/go-user-service/internal/app/config"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"github.com/yasinsaee/go-user-service/pkg/jwt"
	userpb "github.com/yasinsaee/go-user-service/user-service/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// auth service
func (h *Handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	u, err := h.service.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	userAccesses := h.toUserJwtMeta(u)

	tokenConfig := jwt.TokenConfig{
		ID:           u.UniqueID,
		Username:     req.GetUsername(),
		UserAccesses: userAccesses,
	}
	accessToken, _, err := tokenConfig.GenerateAccessToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	refreshToken, _, err := tokenConfig.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	if err := h.service.StoreRefreshToken(u.UniqueID, refreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store refresh token: %v", err)
	}

	return &userpb.LoginResponse{
		User:         h.toUserPB(u, false),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *Handler) Register(ctx context.Context, req *userpb.RegisterUser) (*userpb.UserResponse, error) {
	u := &user.User{
		FirstName:    req.GetFirstName(),
		LastName:     req.GetLastName(),
		Username:     req.GetUsername(),
		ProfileImage: req.GetProfileImage(),
		Email:        req.GetEmail(),
		PhoneNumber:  req.GetPhoneNumber(),
		Password:     req.GetPassword(),
		Group:        req.GetGroup(),
		TenantID:     req.GetTenantId(),
	}
	lType := config.GetEnv("LOGIN_TYPE", "phone")
	username := req.GetUsername()
	switch lType {
	case "phone":
		username = req.GetPhoneNumber()
	case "email":
		username = req.GetEmail()
	case "both":

	}
	if err := h.service.Register(username, u); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &userpb.UserResponse{
		User: h.toUserPB(u, false),
	}, nil
}

func (h *Handler) ResetPassword(ctx context.Context, req *userpb.ResetPasswordUser) (*userpb.UserResponse, error) {
	u, err := h.service.GetByUsername(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to found user: %v", err)
	}

	if err := h.service.ResetPassword(u, req.GetCurrentPassword(), req.GetNewPassword(), req.GetRepeatNewPassword()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reset password user: %v", err)
	}

	return &userpb.UserResponse{
		User: h.toUserPB(u, false),
	}, nil
}

func (h *Handler) UpdatePassword(ctx context.Context, req *userpb.UpdatePasswordUser) (*userpb.UserResponse, error) {
	var (
		u   *user.User
		err error
	)

	if req.GetId() == "" {
		u, err = h.service.GetByUsername(req.GetUsername())
	} else {
		u, err = h.service.GetByID(req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to found user: %v", err)
	}

	if err := h.service.UpdatePassword(u, req.GetNewPassword(), req.GetRepeatNewPassword()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reset password user: %v", err)
	}

	return &userpb.UserResponse{
		User: h.toUserPB(u, false),
	}, nil
}

func (h *Handler) RefreshToken(ctx context.Context, req *userpb.RefreshTokenRequest) (*userpb.RefreshTokenResponse, error) {
	refreshToken := req.GetRefreshToken()
	if refreshToken == "" {
		return nil, status.Errorf(codes.Unauthenticated, "refresh token is required")
	}

	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid or expired refresh token")
	}

	userID := claims.ID

	exists, err := h.service.ValidateRefreshToken(userID, refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if !exists {
		return nil, status.Errorf(codes.Unauthenticated, "refresh token has been revoked")
	}

	u, err := h.service.GetByID(userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	userAccesses := h.toUserJwtMeta(u)

	tc := jwt.TokenConfig{
		ID:           u.UniqueID,
		Username:     u.Username,
		UserAccesses: userAccesses,
	}

	accessToken, accessExpTime, err := tc.GenerateAccessToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	newRefreshToken, _, err := tc.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	} else {
		if revokeErr := h.service.RevokeRefreshToken(userID, refreshToken); revokeErr != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate token")
		}

		if storeErr := h.service.StoreRefreshToken(userID, newRefreshToken); storeErr != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate token")
		}
	}

	return &userpb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    timestamppb.New(accessExpTime),
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *userpb.RefreshTokenRequest) (*userpb.LogoutResponse, error) {
	refreshToken := req.GetRefreshToken()
	if refreshToken == "" {
		return nil, status.Errorf(codes.Unauthenticated, "refresh token is required")
	}

	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid or expired refresh token")
	}

	if claims != nil && claims.ID != "" {
		userID := claims.ID

		if revokeErr := h.service.RevokeRefreshToken(userID, refreshToken); revokeErr != nil {
			return nil, status.Errorf(codes.Internal, "failed to revoke refresh token: %v", revokeErr)
		}
	}

	return &userpb.LogoutResponse{
		Success: true,
		Message: "logged out successfully",
	}, nil
}
