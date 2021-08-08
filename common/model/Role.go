package model

type Role struct {
	BaseModel

	Name        string `json:"name" gorm:"index:index_role_name,unique;not null"`
	Permissions string `json:"permissions"`

	Description string `json:"description"`

	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

// TableName Show the name for Gorm
func (model Role) TableName() string {
	return "Role"
}

// ModelName Show the name of this Loggable
func (model Role) ModelName() string {
	return "Role"
}

// ModelID Show the ID of this Loggable
func (model Role) ModelID() string {
	return model.UUID
}

type CachedPermissionEntry struct {
	CachedEntry

	UserID string

	Global  *CachedPermission
	Tenant map[string]*CachedPermission
}

type CachedPermission map[string]bool
