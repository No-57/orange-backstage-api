package auth

type Token struct {
	ID           uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	AccountID    uint64 `gorm:"column:account_id;index;NOT NULL"`
	AccessToken  string `gorm:"column:access_token;type:text;NOT NULL"`
	RefreshToken string `gorm:"column:refresh_token;NOT NULL"`
}
