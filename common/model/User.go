package model

// Created: 2021-03-19T10:05:30Z

type User struct {
	BaseModel

	UserName string `json:"user_name" gorm:"index:index_user_user_name,unique;not null"`
	Password string `json:"password"`
	Freezed  bool   `json:"freezed"`

	RoleNames string `json:"role_names"`

	Description string `json:"description"`

	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

// TableName Show the name for Gorm
func (model User) TableName() string {
	return "User"
}

// ModelName Show the name of this Loggable
func (model User) ModelName() string {
	return "User"
}

// ModelID Show the ID of this Loggable
func (model User) ModelID() string {
	//return strconv.FormatInt(int64(model.ID), 10)
	return model.UUID
}
