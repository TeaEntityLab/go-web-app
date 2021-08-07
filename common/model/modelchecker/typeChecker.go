package modelchecker

import (
	// "math"
	// "net/url"
	// "sort"
	// "unicode/utf8"
	//
	// jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	//
	// mod "go-web-app/common/model"
	// objectUtils "go-web-app/common/util/objutils"
)

var (
	ErrTypeCheckInvalidFields     = errors.New("Invalid fields")
	ErrTypeCheckNonExistentObject = errors.New("Non-existent object")
)
