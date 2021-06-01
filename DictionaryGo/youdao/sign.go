package youdao

import (
	"crypto/sha256"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

func GetSign(q string, config *Config) (signToStr, salt, curTime string) {
	uuidRandNum := uuid.NewV4()
	salt = uuidRandNum.String()
	input := TruncateString(q)
	stamp := time.Now().Unix()
	curTime = strconv.FormatInt(stamp, 10)
	instr := config.AppKey + input + salt + strconv.FormatInt(stamp, 10) + config.AppSecret
	sign := sha256.Sum256([]byte(instr))
	signToStr = HexNumToString(sign[:])
	return
}

//TruncateString 截断查询字符串(有道官方要求字符串长度大于20则取前10+字符串长度+后10，否则返回q)
func TruncateString(q string) string {
	res := make([]rune, 10)
	temp := []rune(q)
	qLen := len(temp)
	if qLen <= 20 {
		return q
	} else {
		res = temp[:10]
		strQLen := strconv.Itoa(qLen)
		res = append(res, []rune(strQLen)...)
		res = append(res, temp[qLen-10:qLen]...)
		return string(res)
	}
}

//HexNumToString 将十六进制Sign转换为字符串
func HexNumToString(hexnum []byte) (res string) {
	for _, v := range hexnum {
		str := strconv.FormatUint(uint64(v), 16)
		if len(str) == 1 {
			res = res + "0" + str
		} else {
			res += str
		}
	}
	return res
}
