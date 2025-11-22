package otp

import (
	"errors"
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/otp"
	"github.com/yasinsaee/go-user-service/internal/domain/otp/config"
	"github.com/yasinsaee/go-user-service/pkg/logger"
)

// OTPServiceImpl implements otp.OTPService interface
type OTPServiceImpl struct {
	repo        otp.OTPRepository
	provider    otp.OTPProvider    // SMS / Email provider
	limiter     otp.OTPRateLimiter // Redis-based limiter
	codeTTL     time.Duration
	config      config.OTPConfig
	rateLimiter int
}

// NewOTPService creates a new OTP service instance.
func NewOTPService(
	repo otp.OTPRepository,
	provider otp.OTPProvider,
	limiter otp.OTPRateLimiter,
	codeTTLSeconds time.Duration,
	rateLimiter int,
	config config.OTPConfig,

) otp.OTPService {
	return &OTPServiceImpl{
		repo:        repo,
		provider:    provider,
		limiter:     limiter,
		codeTTL:     codeTTLSeconds,
		config:      config,
		rateLimiter: rateLimiter,
	}
}

//
// CRUD
//

func (s *OTPServiceImpl) Create(o *otp.Otp) error {
	return s.repo.Create(o)
}

func (s *OTPServiceImpl) GetByID(id any) (*otp.Otp, error) {
	return s.repo.FindByID(id)
}

func (s *OTPServiceImpl) GetByName(receiver string) (*otp.Otp, error) {
	return s.repo.FindByName(receiver)
}

func (s *OTPServiceImpl) Update(o *otp.Otp) error {
	o.UpdatedAt = time.Now().UTC()
	return s.repo.Update(o)
}

func (s *OTPServiceImpl) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *OTPServiceImpl) ListAll() (otp.Otps, error) {
	return s.repo.List()
}

//
// Business Logic
//

// GenerateCode creates a 6-digit numeric OTP code.
func (s *OTPServiceImpl) GenerateCode() string {
	return otp.GenerateCode(s.config.Length, s.config.Charset)
}

// SaveCode stores the OTP in the database
func (s *OTPServiceImpl) SaveCode(receiver string, code string, ttlSeconds int) error {
	o := &otp.Otp{
		Receiver:  receiver,
		Code:      code,
		Used:      false,
		SendAt:    time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(ttlSeconds) * time.Second),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return s.repo.Create(o)
}

// ValidateCode checks whether code is valid, not expired, and not used.
func (s *OTPServiceImpl) ValidateCode(receiver string, code string) (bool, error) {
	record, err := s.repo.FindByReceiverAndCode(receiver, code)
	if err != nil {
		return false, err
	}
	if (record == &otp.Otp{} || record == nil) {
		return false, errors.New("otp not found")
	}
	if record.ExpiresAt.Before(time.Now().UTC()) {
		return false, errors.New("otp expired")
	}
	if record.Used {
		return false, errors.New("otp used already")
	}
	if record.Code != code {
		return false, errors.New("invalid code")
	}

	// Mark only this OTP as used
	record.Used = true
	if err := s.repo.Update(record); err != nil {
		return false, errors.New("cannot update otp as used")
	}

	_ = s.repo.DeleteExpiredOtps()
	return true, nil
}

// SendCode delivers the OTP using provider (SMS, Email, ...)
func (s *OTPServiceImpl) SendCode(receiver string, code string) error {
	if s.limiter != nil {
		ok, err := s.limiter.CanSend(receiver)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("too many requests, please wait")
		}
	}

	if s.provider != nil {
		if err := s.provider.Send(receiver, code); err != nil {
			logger.Error("otp send failed: ", err.Error())
			return err
		}
	}

	if s.limiter != nil {
		_ = s.limiter.MarkSend(receiver, s.codeTTL)
	}
	return nil
}

//
// Rate Limit
//

func (s *OTPServiceImpl) CanSend(receiver string) (bool, error) {
	if s.limiter == nil {
		return true, nil
	}
	return s.limiter.CanSend(receiver)
}

func (s *OTPServiceImpl) MarkSend(receiver string) error {
	if s.limiter == nil {
		return nil
	}
	return s.limiter.MarkSend(receiver, s.codeTTL)
}
