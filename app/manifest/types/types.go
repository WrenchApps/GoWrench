package types

type HttpMethod string

const (
	HttpMethodGet    HttpMethod = "get"
	HttpMethodPost   HttpMethod = "post"
	HttpMethodPut    HttpMethod = "put"
	HttpMethodPatch  HttpMethod = "patch"
	HttpMethodDelete HttpMethod = "delete"
)

type HashAlg string

const (
	HashAlgSHA512 HashAlg = "SHA-512"
	HashAlgSHA256 HashAlg = "SHA-256"
	HashAlgSHA1   HashAlg = "SHA-1"
	HashAlgMD5    HashAlg = "MD5"
)
