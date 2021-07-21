package router

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"reflect"
	"strings"
	"unsafe"
)

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func Md5UploadFile(f multipart.File) string {
	h := md5.New()
	_, errSeek := f.Seek(0, 0) // 重置文件指针
	if errSeek != nil {        // ignore
		return "error_seek"
	}
	_, err := io.Copy(h, f)
	if err != nil { // ignore
		return "error_copy"
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func Md5(bytes []byte) string {
	h := md5.New()
	h.Write(bytes)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func Md5String(word string) string {
	return Md5([]byte(word))
}
