package model

var Tables = struct {
	WishList string
}{
	"WishList",
}

type WishList struct {
	LineID  string `gorm:"column:line_id" gorm:"primaryKey"`
	Path    string `gorm:"column:path" gorm:"primaryKey"`
	Routine int    `gorm:"column:routine"`
}
