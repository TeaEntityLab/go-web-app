package model

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	//"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type Loggable interface {
	ModelName() string
	ModelID() string
}

type CachedEntry struct {
	CreatedTime time.Time
}

func (object CachedEntry) GetCreatedTime() time.Time {
	return object.CreatedTime
}

type HTTPRequestInfo struct {
	Remote        string `json:"remote"`
	Method        string `json:"method"`
	RequestURI    string `json:"request_uri"`
	Protocol      string `json:"protocol"`
	Host          string `json:"host"`
	ContentLength int64  `json:"content_length"`
	Referer       string `json:"referer"`
	UserAgent     string `json:"user_agent"`
}

type BaseModel struct {
	UUID      string `json:"uuid" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type BaseModelForJoin struct {
	UUID      string `json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

// EnumName NOTE: Copied from proto(github.com/golang/protobuf/proto)
func EnumName(m map[int32]string, v int32) string {
	s, ok := m[v]
	if ok {
		return s
	}
	return strconv.Itoa(int(v))
}

// UnmarshalJSONEnum NOTE: Copied from proto(github.com/golang/protobuf/proto)
func UnmarshalJSONEnum(m map[string]int32, data []byte, enumName string) (int32, error) {
	if data[0] == '"' {
		// New style: enums are strings.
		var repr string
		if err := jsoniter.Unmarshal(data, &repr); err != nil {
			return -1, err
		}
		val, ok := m[repr]
		if !ok {
			return 0, fmt.Errorf("unrecognized enum %s value %q", enumName, repr)
		}
		return val, nil
	}
	// Old style: enums are ints.
	var val int32
	if err := jsoniter.Unmarshal(data, &val); err != nil {
		return 0, fmt.Errorf("cannot unmarshal %#q into enum %s", data, enumName)
	}
	return val, nil
}
