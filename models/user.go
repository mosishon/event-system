package models

import "time"

// کاربر رو تو سیستم نشون میده
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // رمز عبور تو پاسخ‌های JSON نشون داده نمیشه
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ساختار پاسخ برای اطلاعات کاربر (بدون اطلاعات حساس)
// ساختار پاسخ خطای API
// swagger:model
type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ساختار درخواست ورود به سیستم
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ساختار درخواست ثبت نام
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// ساختار پاسخ توکن JWT
type TokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}
