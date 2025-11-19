package otpgrpc

import (
	"context"
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/otp"
	otpPb "github.com/yasinsaee/go-user-service/user-service/otp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	otpPb.UnimplementedOTPServiceServer
	service otp.OTPService
}

func New(service otp.OTPService) *Handler {
	return &Handler{service: service}
}

//
// CRUD
//

func (h *Handler) UpdateOTP(ctx context.Context, req *otpPb.UpdateOTPRequest) (*otpPb.UpdateOTPResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid OTP id")
	}

	o, err := h.service.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "otp not found: %v", err)
	}

	if req.GetUsed() {
		o.Used = true
	}

	if req.GetCode() != "" {
		o.Code = req.GetCode()
	}

	if req.GetTtlSeconds() > 0 {
		o.ExpiresAt = time.Now().Add(time.Duration(req.GetTtlSeconds()) * time.Second)
	}

	if err := h.service.Update(o); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update otp: %v", err)
	}

	return &otpPb.UpdateOTPResponse{
		Otp: toProto(o),
	}, nil
}

func (h *Handler) DeleteOTP(ctx context.Context, req *otpPb.DeleteOTPRequest) (*otpPb.DeleteOTPResponse, error) {
	err := h.service.Delete(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "otp not found: %v", err)
	}

	return &otpPb.DeleteOTPResponse{
		Message: "ok",
	}, nil
}

func (h *Handler) GetOTP(ctx context.Context, req *otpPb.GetOTPRequest) (*otpPb.GetOTPResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id format")
	}

	o, err := h.service.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "otp not found: %v", err)
	}

	return &otpPb.GetOTPResponse{
		Otp: toProto(o),
	}, nil
}

func (h *Handler) ListOTPs(ctx context.Context, req *otpPb.ListOTPsRequest) (*otpPb.ListOTPsResponse, error) {
	list, err := h.service.ListAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list otps: %v", err)
	}

	var pbList []*otpPb.OTP
	for _, item := range list {
		pbList = append(pbList, toProto(&item))
	}

	return &otpPb.ListOTPsResponse{
		Otps: pbList,
	}, nil
}

//
// Business Logic
//

func (h *Handler) RequestOTP(ctx context.Context, req *otpPb.RequestOTPRequest) (*otpPb.RequestOTPResponse, error) {
	receiver := req.GetReceiver()

	// نرخ‌دهی
	ok, err := h.service.CanSend(receiver)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "rate limit error: %v", err)
	}
	if !ok {
		return nil, status.Error(codes.ResourceExhausted, "too many requests, wait before retrying")
	}

	// ایجاد کد
	code := h.service.GenerateCode()

	// ذخیره در DB
	err = h.service.SaveCode(receiver, code, int(req.GetTtlSeconds()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save otp: %v", err)
	}

	// ارسال
	if err := h.service.SendCode(receiver, code); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send otp: %v", err)
	}

	return &otpPb.RequestOTPResponse{
		Message: "otp sent",
	}, nil
}

func (h *Handler) ValidateOTP(ctx context.Context, req *otpPb.ValidateOTPRequest) (*otpPb.ValidateOTPResponse, error) {
	ok, err := h.service.ValidateCode(req.GetReceiver(), req.GetCode())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid otp")
	}

	return &otpPb.ValidateOTPResponse{
		Valid: true,
	}, nil
}

//
// Helper
//

func toProto(o *otp.OTP) *otpPb.OTP {
	return &otpPb.OTP{
		Id:        o.ID.Hex(),
		UserId:    o.UserID.Hex(),
		Code:      o.Code,
		ExpiresAt: o.ExpiresAt.Unix(),
		Used:      o.Used,
	}
}
