package model

import "golang.org/x/crypto/bcrypt"

// User 用户
type User struct {
	Model
	UserName  string `json:"username" gorm:"not null;unique;comment:用户名"`
	Password  string `json:"password" gorm:"not null;comment:密码"`
	Remark    string `json:"remark" gorm:"comment:备注"`
	HeaderImg string `json:"header_img" gorm:"comment:头像url"`
	IsSuper   bool   `json:"is_super" gorm:"not null;comment:是否是超级用户 0:不是 1:是"`
	IsActive  bool   `json:"is_active" gorm:"not null;comment:是否是激活状态 0:不是 1:是"`
}

// Encrypt 加密密码.
func (s *User) Encrypt() (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare 密码比较
func (s *User) Compare(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
}
