package repository

import (
	"server-health/model"

	"gorm.io/gorm"
)

type wishListRepository struct {
	db *gorm.DB
}

type IWishListRepository interface {
	AddWishList(line_id string, path string) error
	RemoveWishList(line_id string, path string) error
	GetPersonWishList(line_id string) ([]model.WishList, error)
}

func NewWishListRepository(db *gorm.DB) wishListRepository {
	return wishListRepository{
		db: db,
	}
}
func (h wishListRepository) AddWishList(line_id string, path string) error {
	wishList := model.WishList{
		LineID: line_id,
		Path:   path,
	}
	err := h.db.Table(model.Tables.WishList).Create(&wishList).Error
	return err
}
func (h wishListRepository) RemoveWishList(line_id string, path string) error {
	wishList := model.WishList{
		LineID: line_id,
		Path:   path,
	}
	err := h.db.Table(model.Tables.WishList).Where("line_id = ?", line_id).Where("path = ?", path).Delete(&wishList).Error
	return err
}
func (h wishListRepository) GetPersonWishList(line_id string) ([]model.WishList, error) {
	wishLists := make([]model.WishList, 0)
	err := h.db.Table(model.Tables.WishList).Find(&wishLists).Error
	return wishLists, err
}
