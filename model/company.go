package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Company ...
type Company struct {
	Common
	Brand     string `json:"brand,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`

	ProjectCount  uint64 `json:"project_count,omitempty"`
	EmployeeCount uint64 `json:"employee_count,omitempty"`
	TeamCount     uint64 `json:"team_count,omitempty"`
}

// CheckUserPermission ..
func (c *Company) CheckUserPermission(db *gorm.DB, userID, permission uint64) (uint64, error) {
	// 验证雇员是否属于企业
	var uc UserCompany
	if err := db.Where("user_id = ? AND company_id = ?", userID, c.ID).First(&uc).Error; err != nil {
		return 0, fmt.Errorf("您不是该企业的雇员(%s)", err)
	}
	// 验证权限
	if permission > 0 && uc.Permission < UCPManager {
		return 0, fmt.Errorf("您不是该企业的管理人员(%d)", uc.Permission)
	}
	return uc.Permission, nil
}