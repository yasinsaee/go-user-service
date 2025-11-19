package otp

import (
	"errors"
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/otp"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"github.com/yasinsaee/go-user-service/pkg/logger"
)

// OTPServiceImpl implements otp.OTPService interface
type OTPServiceImpl struct {
	repo        otp.OTPRepository
	provider    otp.OTPProvider    // SMS / Email provider
	limiter     otp.OTPRateLimiter // Redis-based limiter
	codeTTL     time.Duration
	rateLimit   time.Duration
	userService user.UserService
}

// NewOTPService creates a new OTP service instance.
func NewOTPService(
	repo otp.OTPRepository,
	provider otp.OTPProvider,
	limiter otp.OTPRateLimiter,
	codeTTLSeconds int,
	rateLimitSeconds int,
) otp.OTPService {
	return &OTPServiceImpl{
		repo:      repo,
		provider:  provider,
		limiter:   limiter,
		codeTTL:   time.Duration(codeTTLSeconds) * time.Second,
		rateLimit: time.Duration(rateLimitSeconds) * time.Second,
	}
}

//
// CRUD (همانند PermissionService)
//

func (s *OTPServiceImpl) Create(o *otp.OTP) error {
	return s.repo.Create(o)
}

func (s *OTPServiceImpl) GetByID(id any) (*otp.OTP, error) {
	return s.repo.FindByID(id)
}

func (s *OTPServiceImpl) GetByName(name string) (*otp.OTP, error) {
	return s.repo.FindByName(name)
}

func (s *OTPServiceImpl) Update(o *otp.OTP) error {
	return s.repo.Update(o)
}

func (s *OTPServiceImpl) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *OTPServiceImpl) ListAll() (otp.OTPs, error) {
	return s.repo.List()
}

//
// Business Logic
//

// GenerateCode creates a 6-digit numeric OTP code.
func (s *OTPServiceImpl) GenerateCode() string {
	return otp.GenerateNumericCode(6) // مثل 032491
}

// SaveCode stores the OTP in the database
func (s *OTPServiceImpl) SaveCode(cellphone string, code string, ttlSeconds int) error {
	u, _ := s.userService.GetByUsername(cellphone)
	o := &otp.OTP{
		UserID:    u.ID,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
		Used:      false,
	}
	return s.repo.Create(o)
}

// ValidateCode checks whether code is valid, not expired, and not used.
func (s *OTPServiceImpl) ValidateCode(receiver string, code string) (bool, error) {
	record, err := s.repo.FindByName(receiver) // یا FindByReceiver
	if err != nil {
		return false, err
	}

	if record == nil {
		return false, errors.New("otp not found")
	}

	if record.ExpiresAt.Before(time.Now()) {
		return false, errors.New("otp expired")
	}

	if record.Used {
		return false, errors.New("otp used already")
	}

	if record.Code != code {
		return false, errors.New("invalid code")
	}

	// Mark as used
	record.Used = true
	_ = s.repo.Update(record)

	return true, nil
}

// SendCode delivers the OTP using provider (SMS, Email, ...)
func (s *OTPServiceImpl) SendCode(receiver string, code string) error {

	ok, err := s.limiter.CanSend(receiver)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("too many requests, please wait")
	}

	if err := s.provider.Send(receiver, code); err != nil {
		logger.Error("otp send failed: ", err.Error())
		return err
	}

	_ = s.limiter.MarkSend(receiver, s.rateLimit)
	return nil
}

//
// Rate Limit
//

func (s *OTPServiceImpl) CanSend(receiver string) (bool, error) {
	return s.limiter.CanSend(receiver)
}

func (s *OTPServiceImpl) MarkSend(receiver string) error {
	return s.limiter.MarkSend(receiver, s.rateLimit)
}
