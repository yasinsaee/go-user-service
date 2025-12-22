package otp

import (
	"errors"
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/otp"
	"github.com/yasinsaee/go-user-service/internal/domain/otp/config"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
)

// OTPServiceImpl implements otp.OTPService interface
type OTPServiceImpl struct {
	repo              otp.OTPRepository
	provider          otp.OTPProvider    // SMS / Email provider
	limiter           otp.OTPRateLimiter // Redis-based limiter
	codeTTL           time.Duration
	config            config.OTPConfig
	rateLimiter       int
	maxOTPPerReceiver int
}

// NewOTPService creates a new OTP service instance.
func NewOTPService(
	repo otp.OTPRepository,
	provider otp.OTPProvider,
	limiter otp.OTPRateLimiter,
	codeTTL time.Duration,
	rateLimiter int,
	config config.OTPConfig,
	maxOTPPerReceiver int,
) otp.OTPService {
	return &OTPServiceImpl{
		repo:              repo,
		provider:          provider,
		limiter:           limiter,
		codeTTL:           codeTTL,
		config:            config,
		rateLimiter:       rateLimiter,
		maxOTPPerReceiver: maxOTPPerReceiver,
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

func (s *OTPServiceImpl) ListByType(otpType otp.OtpType) (otp.Otps, error) {
	return s.repo.ListByType(otpType)
}

func (s *OTPServiceImpl) Count(q bson.M) (int, error) {
	return s.repo.Count(q)
}

//
// Business Logic
//

// GenerateCode creates a numeric OTP code.
func (s *OTPServiceImpl) GenerateCode() string {
	return otp.GenerateCode(s.config.Length, s.config.Charset)
}

// SaveCode stores the OTP in the database with type.
func (s *OTPServiceImpl) SaveCode(receiver string, otpType otp.OtpType, code string) error {
	o := &otp.Otp{
		Receiver:  receiver,
		Code:      code,
		Type:      otpType,
		Used:      false,
		SendAt:    time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(s.codeTTL),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return s.repo.Create(o)
}

// ValidateCode checks whether code is valid for a given type, not expired, and not used.
func (s *OTPServiceImpl) ValidateCode(receiver string, otpType otp.OtpType, code string) (bool, error) {
	record, err := s.repo.FindByReceiverAndCodeAndType(receiver, code, otpType)
	if err != nil {
		return false, err
	}
	if record == nil {
		return false, errors.New("otp not found")
	}
	if record.ExpiresAt.Before(time.Now().UTC()) {
		return false, errors.New("otp expired")
	}
	if record.Used {
		return false, errors.New("otp used already")
	}

	// Mark OTP as used
	record.Used = true
	if err := s.repo.Update(record); err != nil {
		return false, errors.New("cannot mark otp as used")
	}

	_ = s.repo.DeleteExpiredOtps()
	return true, nil
}

// SendCode delivers the OTP using provider (SMS, Email, ...) for a given type
func (s *OTPServiceImpl) SendCode(receiver string, otpType otp.OtpType, code string) error {
	ok, err := s.CanSend(receiver)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("too many requests, please wait")
	}

	if s.provider != nil {
		if err := s.provider.Send(receiver, code); err != nil {
			logger.Error("otp send failed: ", err.Error())
			return err
		}
	}

	if err := s.MarkSend(receiver); err != nil {
		logger.Error("otp marked failed: ", err.Error())
		return err
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

// Hard limit Check by receiver and type
func (s *OTPServiceImpl) CheckHardLimit(receiver string, otpType otp.OtpType) (bool, error) {
	if s.maxOTPPerReceiver > 0 {
		count, _ := s.Count(bson.M{"receiver": receiver, "type": otpType})
		if count >= s.maxOTPPerReceiver {
			return false, errors.New("too many OTP requests, contact support")
		}
	}
	return true, nil
}
