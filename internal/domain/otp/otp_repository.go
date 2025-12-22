package otp

import "go.mongodb.org/mongo-driver/bson"

type OTPRepository interface {
	// ایجاد OTP جدید
	Create(otp *Otp) error

	// پیدا کردن OTP بر اساس ID
	FindByID(id any) (*Otp, error)

	// پیدا کردن OTP بر اساس نام (اگر نیاز باشد)
	FindByName(name string) (*Otp, error)

	// بروزرسانی OTP
	Update(otp *Otp) error

	// حذف OTP
	Delete(id any) error

	// لیست تمام OTPها
	List() (Otps, error)

	// نسخه قدیمی برای backward compatibility
	FindByReceiverAndCode(receiver, code string) (*Otp, error)

	// نسخه جدید با Type
	FindByReceiverAndCodeAndType(receiver, code string, otpType OtpType) (*Otp, error)

	// لیست OTPها بر اساس Type
	ListByType(otpType OtpType) (Otps, error)

	// حذف OTPهای منقضی شده
	DeleteExpiredOtps() error

	// شمارش OTPها با query دلخواه
	Count(q bson.M) (int, error)
}
