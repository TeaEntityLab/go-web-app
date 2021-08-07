package fileutils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	//modMine "go-web-app/model"
	"go-web-app/common/util/timeutils"
)

const (
	MEDIA_FILE_TYPE_IMAGE  = "image"
	MEDIA_FILE_TYPE_VIDEO  = "video"
	MEDIA_FILE_TYPE_ASSETS = "assets"
	MEDIA_FILE_TYPE_ELSE   = "else"

	MEDIA_FILE_EXT_PNG  = ".png"
	MEDIA_FILE_EXT_JPG  = ".jpg"
	MEDIA_FILE_EXT_JPEG = ".jpeg"
	MEDIA_FILE_EXT_GIF  = ".gif"
	MEDIA_FILE_EXT_MP4  = ".mp4"
	MEDIA_FILE_EXT_MOV  = ".mov"
	MEDIA_FILE_EXT_ZIP  = ".zip"

	MEDIA_FILE_EXT_CSV  = ".csv"
	MEDIA_FILE_EXT_XLSX = ".xlsx"
	MEDIA_FILE_EXT_XLSM = ".xlsm"
	MEDIA_FILE_EXT_XLTM = ".xltm"
	MEDIA_FILE_EXT_XLTX = ".xltx"

	Kibibyte Base2Bytes = 1024
	KiB                 = Kibibyte
	Mebibyte            = Kibibyte * 1024
	MiB                 = Mebibyte
	Gibibyte            = Mebibyte * 1024
	GiB                 = Gibibyte
	Tebibyte            = Gibibyte * 1024
	TiB                 = Tebibyte
	Pebibyte            = Tebibyte * 1024
	PiB                 = Pebibyte
	Exbibyte            = Pebibyte * 1024
	EiB                 = Exbibyte
)

type Base2Bytes int64

func GetMimeByPath(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return GetMimeByOsFile(file)
}

func GetMimeByOsFile(file *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Reset the read pointer if necessary.
	file.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetMimeByMultipartFile(file *multipart.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := (*file).Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Reset the read pointer if necessary.
	(*file).Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func CopyToTempFileByMultipartFile(extension string, input *multipart.File) (*os.File, string, error) {
	if input == nil {
		return nil, "", nil
	}

	randName := GenRandNameWithExtension(extension)
	// filename = GenTempFilePath(randName)
	// log.Debug(filename)

	newDirectoryName := timeutils.GetTimestampString()
	newDirectoryPath := filepath.Join(os.TempDir(), newDirectoryName)
	os.MkdirAll(newDirectoryPath, os.ModePerm)

	file, err := ioutil.TempFile(newDirectoryPath, randName)
	if err != nil || file == nil {
		log.Debug("Err", err)
		return file, "", err
	}
	//defer file.Close()
	filePath := file.Name()

	(*input).Seek(0, 0)
	_, err = io.Copy(file, *input)
	return file, filePath, err
}

func GenTempFilePath(fileName string) string {
	return filepath.Join(os.TempDir(), fileName)
}

func GenTempFileName(prefix, suffix string) string {
	return prefix + GenRandName() + suffix
}

func GenRandName() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)

	return hex.EncodeToString(randBytes)
}

func GenRandNameWithExtension(extension string) string {
	return fmt.Sprintf("%s.%s", GenRandName(), extension)
}

func GetExtension(url string) string {
	return strings.ToLower(filepath.Ext(url))
}

//func GenVideoFilenameByResolutionAndFormat(url, resolution, videoFormat string) string {
//	urlWithoutExtension := strings.TrimSuffix(url, filepath.Ext(url))
//	return fmt.Sprintf("%s_%s.%s", urlWithoutExtension, GenVideoSignatureByResolutionAndFormat(resolution, videoFormat), videoFormat)
//}
//
//func GenVideoSignatureByResolutionAndFormat(resolution, videoFormat string) string {
//	return fmt.Sprintf(modMine.CreativeVideoExtraFormatsFormatString, resolution, videoFormat)
//}

// GenMd5FromFile
func GenMd5FromFile(file *os.File) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, file)

	if err != nil {
		//
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), err
}

// JoinUrl ...
func JoinUrl(prefix string, relativePaths ...string) string {
	u, err := url.Parse(prefix)
	if err != nil {
		return strings.Join(relativePaths, "/")
	}

	u.Path = path.Join(append([]string{u.Path}, relativePaths...)...)
	return u.String()
}

// Copy the src file to dst. Any existing file will be overwritten and will not copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
