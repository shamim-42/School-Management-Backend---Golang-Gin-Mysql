package account

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

type SignupPayload struct {
	Nickname string `gorm:"size:255;not null;unique" json:"nickname"`
	Email    string `gorm:"size:100; null;unique" json:"email"`
	Mobile   string `gorm:"size:100;not null;unique" json: "mobile"`
	Mobile2  string `gorm:"size:100;not null;unique" json: "mobile2"`
	Address  string `gorm:"size:255;" json:"address"`
	Password string `gorm:"size:100;not null;" json:"password"`
}

type SigninPayload struct {
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserPayload struct {
	Nickname string `gorm:"size:255;not null;unique" json:"nickname"`
	Email    string `gorm:"size:100; null;unique" json:"email"`
	Mobile   string `gorm:"size:100;not null;unique" json: "mobile"`
	Mobile2  string `gorm:"size:100;not null;unique" json: "mobile2"`
	Address  string `gorm:"size:255; null;" json:"address"`
}

type Role struct {
	Id        uint32 `gorm:"primary_key; auto_increment" json: "id"`
	Role      string `gorm: "size: 16" json:"role"`
	CreatedAt string `gorm: "default: CURRENT_TIMESTAMP" json: "created_at"`
	UpdatedAt string `gorm: "default: CURRENT_TIMESTAMP" json: "updated_at"`
}

type User struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Nickname string `gorm:"size:255;null" json:"nickname"`
	Fullname string `gorm:"size:255;null" json:"fullname"`
	Email    string `gorm:"size:100;unique" json:"email"`
	Mobile   string `gorm:"size:32;not null;unique" json: "mobile"`
	Password string `gorm:"size:100;not null;" json:"password"`
	Address  string `gorm:"size:1024; null" json:"address"`

	// CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	// UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}
