package bootstrap

import (
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/utils"
	"gorm.io/gorm"
)

func InitializeAdmin(db *gorm.DB, defaultName, shareName, passwordFile string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.User{}).Count(&count).Error; err != nil {
			return err
		}
		if count != 0 {
			return nil
		}
		password, err := ensureSecret(passwordFile, 16)
		if err != nil {
			return err
		}
		hash, err := utils.EncryptPassword(password)
		if err != nil {
			return err
		}
		group := model.Group{Name: defaultName, Type: model.GroupTypeDefault}
		if err := tx.Where("type = ?", model.GroupTypeDefault).FirstOrCreate(&group).Error; err != nil {
			return err
		}
		share := model.Group{Name: shareName, Type: model.GroupTypeShare}
		if err := tx.Where("type = ?", model.GroupTypeShare).FirstOrCreate(&share).Error; err != nil {
			return err
		}
		isAdmin := true
		admin := model.User{
			Username: "admin",
			Nickname: "Admin",
			Status:   model.COMMON_STATUS_ENABLE,
			IsAdmin:  &isAdmin,
			GroupId:  group.Id,
			Password: hash,
		}
		return tx.Create(&admin).Error
	})
}
