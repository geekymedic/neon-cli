package tool

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
)

func MD5UUID4() string {
	hasher := md5.New()
	txt := uuid.New()
	hasher.Write(txt[:])
	return hex.EncodeToString(hasher.Sum(nil))
}