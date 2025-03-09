package costco

import "errors"

const (
	BLOB_UNKNOWN          = "BLOB_UNKNOWN"
	BLOB_UPLOAD_INVALID   = "BLOB_UPLOAD_INVALID"
	BLOB_UPLOAD_UNKNOWN   = "BLOB_UPLOAD_UNKNOWN"
	DIGEST_INVALID        = "DIGEST_INVALID"
	MANIFEST_BLOB_UNKNOWN = "MANIFEST_BLOB_UNKNOWN"
	MANIFEST_INVALID      = "MANIFEST_INVALID"
	MANIFEST_UNKNOWN      = "MANIFEST_UNKNOWN"
	MANIFEST_UNVERIFIED   = "MANIFEST_UNVERIFIED"
	NAME_INVALID          = "NAME_INVALID"
	NAME_UNKNOWN          = "NAME_UNKNOWN"
	SIZE_INVALID          = "SIZE_INVALID"
	TAG_INVALID           = "TAG_INVALID"
	UNAUTHORIZED          = "UNAUTHORIZED"
	DENIED                = "DENIED"
	UNSUPPORTED           = "UNSUPPORTED"

	FILESYSTEM = "file-system"
	S3         = "s3"

	BASIC = "auth-basic"
)

type FsLayers struct {
	BlobSum string `json:"blobSum"`
}

type ImageManifest struct {
	Name      string     `json:"name"`
	Tag       string     `json:"tag"`
	FsLayers  []FsLayers `json:"fsLayers"`
	History   string     `json:"history"`
	Signature string     `json:"signature"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

var (
	errNoStorageSpecified = errors.New("error, no storage type specified")
	errInvalidStorageType = errors.New("error, storage type specified invalid")
	errNoAuthSpecified    = errors.New("error, no auth type specified")
	errInvalidAuthType    = errors.New("error, auth type specified invalid")
	errRepositoryNotFound = errors.New("error, repository not found")
)

type BasicAuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type FileSystemConfig struct {
	filePath string `yaml:"filePath"`
}

type MinioConfig struct {
}

type StoreConfig struct {
	kind             string            `yaml:"kind"`
	FileSystemConfig *FileSystemConfig `yaml:"fsConfig"`
	MinioConfig      *MinioConfig      `yaml:"minioConfig"`
}

type AuthConfig struct {
	kind            string           `yaml:"kind"`
	BasicAuthConfig *BasicAuthConfig `yaml:"basicAuthConfig"`
}

type Configs struct {
	Storage StoreConfig `yaml:"storage"`
	Auth    AuthConfig  `yaml:"auth"`
}
