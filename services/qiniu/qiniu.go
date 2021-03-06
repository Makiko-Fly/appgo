package qiniu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/oxfeeefeee/appgo"
	qndigest "github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/conf"
	qnio "github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
	qnurl "github.com/qiniu/api.v6/url"
	"io"
	"strconv"
	"strings"
)

var (
	bucketName string
	domain     string
	useHttps   bool
	putPolicy  rs.PutPolicy
	getPolicy  rs.GetPolicy
)

func init() {
	params := appgo.Conf.Qiniu
	conf.ACCESS_KEY = params.AccessKey
	conf.SECRET_KEY = params.Secret
	bucketName = params.Bucket
	domain = params.Domain
	useHttps = params.UseHttps
	putPolicy = rs.PutPolicy{
		Expires:    uint32(params.DefaultExpires),
		InsertOnly: 1,
	}
	getPolicy = rs.GetPolicy{
		Expires: uint32(params.DefaultExpires),
	}
}

func Domain() string {
	return domain
}

func GetUrl(key string, expiryUnix int64, params string) string {
	baseUrl := makeBaseUrl(key)
	url := baseUrl + "?" + params
	if expiryUnix == 0 {
		return getPolicy.MakeRequest(url, nil)
	}
	return makeRequest(url, expiryUnix, nil)
}

func PublicGetUrl(key string) string {
	return makeBaseUrl(key)
}

func NoKeyToken() string {
	putPolicy.Scope = bucketName
	return putPolicy.Token(nil)
}

func PutToken(key string) string {
	putPolicy.Scope = bucketName + ":" + key
	return putPolicy.Token(nil)
}

func PutFile(key string, r io.Reader, size int64) (string, error) {
	putPolicy.Scope = bucketName + ":" + key
	token := putPolicy.Token(nil)
	//var ret qnio.PutRet
	if err := qnio.Put2(nil, nil /*&ret*/, token, key, r, size, nil); err != nil {
		return "", err
	} else {
		return makeBaseUrl(key), nil
	}
}

func FetchToken(url, key string) (string, string, string) {
	encodedUrl := base64.URLEncoding.EncodeToString([]byte(url))
	encodedTo := base64.URLEncoding.EncodeToString([]byte(bucketName + ":" + key))
	fullPath := "/fetch/" + encodedUrl + "/to/" + encodedTo

	h := hmac.New(sha1.New, []byte(conf.SECRET_KEY))
	io.WriteString(h, fullPath+"\n")
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	token := conf.ACCESS_KEY + ":" + sign
	return fullPath, makeBaseUrl(key), token
}

func makeBaseUrl(key string) string {
	//return rs.MakeBaseUrl(domain, key)
	p := "http"
	if useHttps {
		p = "https"
	}
	return p + "://" + domain + "/" + qnurl.Escape(key)
}

func makeRequest(baseUrl string, expiryUnix int64, mac *qndigest.Mac) (privateUrl string) {
	if strings.Contains(baseUrl, "?") {
		baseUrl += "&e="
	} else {
		baseUrl += "?e="
	}
	baseUrl += strconv.FormatInt(expiryUnix, 10)

	token := qndigest.Sign(mac, []byte(baseUrl))
	return baseUrl + "&token=" + token
}
