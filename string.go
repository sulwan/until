package until

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// Md5Sum 获取字符串的md5值
func Md5Sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GenRandomString 生成随机字符串
// length 生成长度
// specialChar 是否生成特殊字符
func GenRandomString(length int, specialChar string) string {

	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special := "!@#%$*.="

	if specialChar == "yes" {
		letterBytes = letterBytes + special
	}

	chars := []byte(letterBytes)

	if length == 0 {
		return ""
	}

	clen := len(chars)
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return ""
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

// CryptPassword 加密密码
func CryptPassword(password, salt string) string {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	return fmt.Sprintf("%x", dk)
}

// 简化Base64
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// 结构体拼接字符串
func QueryStr(query_balance_post interface{}, upost bool, slice []string) string {
	var splice string
	var post_string []string
	t := reflect.TypeOf(query_balance_post)
	v := reflect.ValueOf(query_balance_post)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Type().Kind() == reflect.Struct {
			structField := v.Field(i).Type()
			for j := 0; j < structField.NumField(); j++ {
				str := structField.Field(j).Name
				if upost {
					post_string = append(post_string, str)
				} else {
					if !InStringSlice(str, slice) {
						post_string = append(post_string, str)
					}
				}
			}
			continue
		}
		post_string = append(post_string, t.Field(i).Name)
	}
	sort.Strings(post_string)
	for _, vel := range post_string {
		t := v.FieldByName(vel)
		t_string := formatAtom(t)
		if len(t_string) > 2 {
			if !upost {
				splice = splice + vel + "=" + t_string + "&"
			} else {
				splice = splice + vel + "=" + url.QueryEscape(t_string) + "&"
			}
		}
	}
	splice = strings.Trim(splice, "&")
	splice = strings.Replace(splice, "%22", "", -1)
	splice = strings.Replace(splice, "\"", "", -1)
	if upost {
		// fmt.Println("--------------提交字符串-------------")
		// fmt.Println(splice)
	} else {
		// fmt.Println("--------------签名字符串-------------")
		splice = strings.Replace(splice, "\"", "", -1)
	}
	return splice
}

// 根据反射类型进行转换
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', 2, 32)
	default:
		return v.Type().String() + " value"
	}
}
