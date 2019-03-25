package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/riposa/utils/log"
	"github.com/riposa/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	earthRadius = 6378.137

	accessKey = ""
	secretKey = ""
	bucket    = "testhengha"

	getWxaCodeInterface = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
)
type TimeIt struct {
	start time.Time
	end   time.Time

	moduleName string
}

var (
	utilsLogger = log.New()
)

func rad(n float64) float64 {
	return n * math.Pi / 180.0
}

func GetWxaCode(token, scene, page string) (string, error) {

	resp, err := Requests.PostJsonWithQueryString(getWxaCodeInterface, map[string]interface{}{"scene": scene, "page": page, "is_hyaline": true, "width": 100}, nil, map[string]string{"access_token": token})
	if err != nil {
		return "", err
	}
	img := resp.Body()

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	data := img
	dataLen := int64(len(data))
	err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		utilsLogger.Exception(err)
		return "", err
	}
	return ret.Hash, nil
}

// distance unit is km
func PointSquareWithDistance(lat, lng, distance float64) (latScope [2]float64, lngScope [2]float64) {
	radLat := rad(lat)
	lngWidth := distance / (math.Cos(radLat) * 111.0)
	latWidth := distance / 111.0
	latScope[0] = lat + latWidth
	latScope[1] = lat - latWidth
	lngScope[0] = lng + lngWidth
	lngScope[1] = lng - lngWidth
	return
}

func calculateTwoPointDistance(lat1, lat2, lng1, lng2 float64) float64 {
	radLat1 := rad(lat1)
	radLat2 := rad(lat2)
	latDiff := radLat1 - radLat2
	lngDiff := rad(lng1) - rad(lng2)
	res := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(latDiff/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(lngDiff/2), 2)))
	res = res * earthRadius
	res = math.Round(res*10000) / 10000
	return res
}

// distance unit is km
func JudgeExistInRegion(lat1, lat2, lng1, lng2, distance float64) bool {
	if calculateTwoPointDistance(lat1, lat2, lng1, lng2) <= distance {
		return true
	} else {
		return false
	}
}

func MpGetUserFromHeader(c *gin.Context) (string, error) {
	author := c.GetHeader("Authorization")
	if author != "" {
		return strings.TrimPrefix(author, "UserID "), nil
	} else {
		return "", errors.New(7010)
	}
}

func GetDBFromContext(c *gin.Context, schema string) (*gorm.DB, error) {
	l, exist := c.Get(schema + "_conn")
	if !exist {
		return nil, errors.New(310)
	}
	conn, ok := l.(*gorm.DB)
	if !ok {
		return nil, errors.New(311)
	}
	return conn, nil
}

func CheckFieldLength(v string, max int, cn string) (bool, errors.Error) {
	vSlice := strings.Split(v, "")
	if len(vSlice) > max {
		return false, errors.NewFormat(30, cn, max)
	}
	return true, errors.Error{}
}

func (t *TimeIt) End() {
	t.end = time.Now()

	utilsLogger.Infof(fmt.Sprintf("Goroutine [%s] cost %.3f ms", t.moduleName, float64(t.end.UnixNano()-t.start.UnixNano())/1e6))
}

func NewTimeIt(module string) *TimeIt {
	var t TimeIt

	t.start = time.Now()
	t.moduleName = module
	return &t
}

func IsLoggedIn(c *gin.Context) bool {
	var isLoggedIn bool

	t, exist := c.Get("is_logged_in")
	if !exist {
		isLoggedIn = false
	} else {
		if tt, ok := t.(bool); ok {
			isLoggedIn = tt
		} else {
			isLoggedIn = false
		}
	}
	return isLoggedIn
}

func GetUserIDFromContext(c *gin.Context) int {
	var userID int

	t, exist := c.Get("user_id")
	if !exist {
		userID = 0
	} else {
		if tt, ok := t.(int); ok {
			userID = tt
		} else {
			userID = 0
		}
	}
	return userID
}

func Factorial(n int) int {
	var r int
	r = 1
	for ; n > 0; n-- {
		r *= n
	}
	return r
}

func Combination(n, m int) int {
	if n < m {
		return 0
	}
	return Factorial(n) / Factorial(m) * Factorial(n-m)
}

func CombinationResult(n int, m int) [][]int {
	if m < 1 || m > n {
		fmt.Println("Illegal argument. Param m must between 1 and len(nums).")
		return [][]int{}
	}

	result := make([][]int, 0, Combination(n, m))
	indexes := make([]int, n)
	for i := 0; i < n; i++ {
		if i < m {
			indexes[i] = 1
		} else {
			indexes[i] = 0
		}
	}

	result = addTo(result, indexes)
	for {
		find := false

		for i := 0; i < n-1; i++ {
			if indexes[i] == 1 && indexes[i+1] == 0 {
				find = true
				indexes[i], indexes[i+1] = 0, 1
				if i > 1 {
					moveOneToLeft(indexes[:i])
				}
				result = addTo(result, indexes)
				break
			}
		}
		if !find {
			break
		}
	}
	return result
}

func addTo(arr [][]int, ele []int) [][]int {
	newEle := make([]int, len(ele))
	copy(newEle, ele)
	arr = append(arr, newEle)
	return arr
}

func moveOneToLeft(leftNums []int) {
	sum := 0
	for i := 0; i < len(leftNums); i++ {
		if leftNums[i] == 1 {
			sum++
		}
	}
	for i := 0; i < len(leftNums); i++ {
		if i < sum {
			leftNums[i] = 1
		} else {
			leftNums[i] = 0
		}
	}
}

func FindNumsByIndexes(nums []int, indexes [][]int) [][]int {
	if len(indexes) == 0 {
		return [][]int{}
	}
	result := make([][]int, len(indexes))
	for i, v := range indexes {
		line := make([]int, 0)
		for j, v2 := range v {
			if v2 == 1 {
				line = append(line, nums[j])
			}
		}
		result[i] = line
	}
	return result
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func Bargain(target, current, remaining int, seed int64) (new int, over bool) {
	var step = 80
	if current-target <= 0 {
		return -1, true
	}
	if remaining < 0 {
		return -1, true
	} else if remaining == 0 {
		return current - target, true
	}
	rd := rand.New(rand.NewSource(seed))
	standard := float64(current-target) / float64(remaining) / 2
	param := 1 + float64(rd.Intn(step))/100.0
	result := standard * param
	return int(result), false
}