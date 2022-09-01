package types

type User struct {
	ID          int    `gorm:"primaryKey;autoIncrement;type:uint"`
	Name        string `gorm:"type:varchar(32);uniqueIndex"`
	Description string
}
