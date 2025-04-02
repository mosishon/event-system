package services

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/event-system/models"
	"github.com/event-system/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// سرویس احراز هویت که منطق کسب و کار مربوط به احراز هویت رو مدیریت میکنه
type AuthService struct {
	UserRepo *repositories.UserRepository
}

// یه نمونه جدید از سرویس احراز هویت میسازه
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

// ثبت نام کاربر رو انجام میده
func (s *AuthService) Register(req models.RegisterRequest) (*models.TokenResponse, error) {
	// چک میکنه ببینه نام کاربری قبلا وجود داره یا نه
	_, err := s.UserRepo.GetByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// چک میکنه ببینه ایمیل قبلا وجود داره یا نه
	_, err = s.UserRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// رمز عبور رو هش میکنه
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("error processing registration")
	}

	// کاربر رو میسازه
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.UserRepo.Create(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, errors.New("error processing registration")
	}

	// توکن رو میسازه
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, errors.New("error processing registration")
	}

	return &models.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// ورود کاربر رو مدیریت میکنه
func (s *AuthService) Login(req models.LoginRequest) (*models.TokenResponse, error) {
	// کاربر رو با ایمیل پیدا میکنه
	user, err := s.UserRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// رمز عبور رو چک میکنه
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// توکن رو میسازه
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, errors.New("error processing login")
	}

	return &models.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// کاربر رو با آیدی پیدا میکنه
func (s *AuthService) GetUserByID(id int) (*models.UserResponse, error) {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// یه توکن JWT برای کاربر میسازه
func (s *AuthService) generateToken(user *models.User) (string, time.Time, error) {
	// کلید رمزنگاری JWT رو از متغیرهای محیطی میگیره یا از مقدار پیشفرض استفاده میکنه
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // کلید رمزنگاری پیشفرض (حتما تو محیط پروداکشن عوضش کن)
	}

	// زمان انقضای توکن رو تنظیم میکنه (24 ساعت)
	expiresAt := time.Now().Add(24 * time.Hour)

	// اطلاعات توکن رو میسازه
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      expiresAt.Unix(),
	}

	// توکن رو میسازه
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// توکن رو امضا میکنه
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// توکن JWT رو اعتبارسنجی میکنه و آیدی کاربر رو برمیگردونه
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	// کلید رمزنگاری JWT رو از متغیرهای محیطی میگیره یا از مقدار پیشفرض استفاده میکنه
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // کلید رمزنگاری پیشفرض (حتما تو محیط پروداکشن عوضش کن)
	}

	// توکن رو پردازش میکنه
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// روش امضا رو اعتبارسنجی میکنه
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// توکن رو اعتبارسنجی میکنه
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// تاریخ انقضا رو چک میکنه
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return 0, errors.New("token expired")
		}

		// آیدی کاربر رو میگیره
		userID := int(claims["id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}
