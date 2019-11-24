package random

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var tokenRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
var codeRunes = []rune("0123456789")

func GeneratePassword(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = tokenRunes[rand.Intn(len(tokenRunes))]
	}
	return string(b)
}

func GenerateCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = codeRunes[rand.Intn(len(codeRunes))]
	}
	return string(b)
}

func CheckPasswd(pwd string) string {

	password := []byte(pwd)

	// 长度
	length := len(pwd)
	if length < 8 || length > 30 {
		return "Field password. The length must be between 8 and 30."
	}

	// 不能出现单词
	word := `((password)|(PASSWORD)`
	pattern_word := regexp.MustCompile(word)
	if pattern_word.Match(password) {
		return "Field password. Cannot appear the following words. password、PASSWORD"
	}

	// 包含不支持的字符
	if !regexp.MustCompile(`^[` + "`" + `()~!@#$%^&*\-_+={}\[\]:";'<>,.?/0-9a-zA-Z]+$`).Match(password) {
		return "Field password. Found unsupported characters"
	}

	// 必须数字、大小写字母、特殊字符中的三类
	include := 0
	if regexp.MustCompile(`[0-9]`).Match(password) {
		include = include + 1
	}
	if regexp.MustCompile(`[a-z]`).Match(password) {
		include = include + 1
	}
	if regexp.MustCompile(`[A-Z]`).Match(password) {
		include = include + 1
	}
	if regexp.MustCompile(`[` + "`" + `()~!@#$%^&*\-_+={}\[\]:";'<>,.?/]`).Match(password) {
		include = include + 1
	}
	if include < 3 {
		return "Field password. Must contains at least three types of numbers, uppercase letters, lowercase letters, and special characters"
	}

	// 不能含有三个连贯数字
	if regexp.MustCompile(`(012)|(123)|(234)|(345)|(456)|(567)|(678)|(789)|(987)|(876)|(765)|(654)|(543)|(432)|(321)|(210)`).Match(password) {
		return "Can't contain three consecutive numbers"
	}

	// 不能含有三个连续字母
	if regexp.MustCompile(`(abc)|(bcd)|(cde)|(def)|(efg)|(fgh)|(ghi)|(hij)|(ijk)|(jkl)|(klm)|(lmn)|(mno)|(nop)|(opq)|(pqr)|(qrs)|(rst)|(stu)|(tuv)|(uvw)|(vwx)|(wxy)|(xyz)|(zyx)|(yxw)|(xwv)|(wvu)|(vut)|(uts)|(tsr)|(srq)|(rqp)|(qpo)|(pon)|(onm)|(nml)|(mlk)|(lkj)|(kji)|(jih)|(ihg)|(hgf)|(gfe)|(fed)|(edc)|(dcb)|(cba)|(qwe)|(wer)|(ert)|(rty)|(tyu)|(yui)|(uio)|(iop)|(poi)|(oiu)|(iuy)|(uyt)|(ytr)|(tre)|(rew)|(ewq)|(asd)|(sdf)|(dfg)|(fgh)|(ghj)|(hjk)|(jkl)|(lkj)|(kjh)|(jhg)|(hgf)|(gfd)|(fds)|(dsa)|(zxc)|(xcv)|(cvb)|(vbn)|(bnm)|(mnb)|(nbv)|(bvc)|(vcx)|(cxz)|(qaz)|(wsx)|(edc)|(rfv)|(tgb)|(yhn)|(ujm)|(zaq)|(xsw)|(cde)|(vfr)|(bgt)|(nhy)|(mju)`).Match([]byte(strings.ToLower(pwd))) {
		return "Can't contain three consecutive numbers"
	}
	return ""
}
