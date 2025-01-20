// Only for those special documents that belong to Admin Documents DB.
// package modelsCommon
package models

type DocumentPrivilegeTable struct {
	UserRole   string `gorm:"type:varchar(3);size:3;default:'STU';not null;check:user_role IN ('STU','VOL','COR','ADM');primaryKey" json:"-"`
	DocumentID string `gorm:"not null;primaryKey"` // Part of the composite key
	CanRead    bool   `gorm:"default:false"`
	CanWrite   bool   `gorm:"default:false"`
	CanDelete  bool   `gorm:"default:false"`
}

func (DocumentPrivilegeTable) TableName() string {
	return "public.document_privilege_table"
}
