package model

// AdminTOTP 管理员 TOTP 二次校验绑定记录。
// 一个 admin 至多一条记录；secret 加密存储（AES-256-GCM）。
type AdminTOTP struct {
	Model
	AdminID        int64  `gorm:"uniqueIndex;not null" json:"admin_id"`
	SecretCipher   string `gorm:"type:text;not null" json:"-"` // base64 编码的 AES-GCM 密文，绝不输出
	Enabled        bool   `gorm:"not null;default:false" json:"enabled"`
	RecoveryCodes  string `gorm:"type:text" json:"-"` // JSON 数组，元素 {"hash":"bcrypt-hash","used_at":"..."}
	LastUsedAt     *int64 `json:"last_used_at,omitempty"`
}

// TableName 显式表名，避免 gorm 自动复数化产生歧义。
func (AdminTOTP) TableName() string {
	return "admin_totp"
}
