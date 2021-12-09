package repository

import (
	"server-health/model"
	"time"

	"gorm.io/gorm"
)

type wishListRepository struct {
	db *gorm.DB
}

type IWishListRepository interface {
	AddWishList(line_id string, path string) error
	RemoveWishList(line_id string, path string) error
	GetPersonWishList(line_id string) ([]model.WishList, error)
	GetCheckableWishList() ([]model.WishList, error)
	UpdateWishListFailure(line_id string, path string) error
}

func NewWishListRepository(db *gorm.DB) wishListRepository {
	return wishListRepository{
		db: db,
	}
}
func (h wishListRepository) AddWishList(line_id string, path string) error {
	wishList := model.WishList{
		LineID:      line_id,
		Path:        path,
		LastFailure: nil,
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
	err := h.db.Table(model.Tables.WishList).Where("line_id = ?", line_id).Find(&wishLists).Error
	return wishLists, err
}
func (h wishListRepository) GetCheckableWishList() ([]model.WishList, error) {
	wishLists := make([]model.WishList, 0)
	err := h.db.Table(model.Tables.WishList).Where("last_failure IS NULL").Find(&wishLists).Error
	return wishLists, err
}

func (h wishListRepository) UpdateWishListFailure(line_id string, path string) error {
	now := time.Now()
	err := h.db.Table(model.Tables.WishList).Where("line_id = ?", line_id).Where("path = ?", path).Update("last_failure", now).Error
	return err
}
