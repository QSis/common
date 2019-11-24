package client

import (
	"github.com/cihub/seelog"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func RequestStructWithHeader(uri string, method string, data interface{}, header map[string]string, isJson bool, param ...interface{}) error {
	method = strings.ToUpper(method)
	client := gorequest.New()
	switch method {
	case "GET":
		client = client.Get(uri)
	case "POST":
		client = client.Post(uri)
	case "DELETE":
		client = client.Delete(uri)
	case "PUT":
		client = client.Put(uri)
	default:
		return errors.New("Unsupported method")
	}
	if isJson {
		for _, v := range param {
			client.Send(v)
		}
	} else {
		for _, v := range param {
			client = client.Query(v)
		}
	}
	for k, v := range header {
		client = client.Set(k, v)
	}
	res, body, errs := client.EndStruct(data)
	var err error = nil
	if errs != nil && len(errs) > 0 {
		seelog.Warn(method+" "+uri+" failed. Body: "+string(body), ", code: ", res, ", param: ", param)
		err = errs[0]
	} else {
		seelog.Debug(method+" "+uri+" success. Body: ", string(body))
	}
	return err
}
func GetStruct(uri string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "GET", data, nil, false, param...)
}
func PostStruct(uri string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "POST", data, nil, false, param...)
}
func PutStruct(uri string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "PUT", data, nil, false, param...)
}
func DeleteStruct(uri string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "DELETE", data, nil, false, param...)
}
func GetStructWithHeader(uri string, header map[string]string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "GET", data, header, false, param...)
}
func PostStructWithHeader(uri string, header map[string]string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "POST", data, header, false, param...)
}
func PutStructWithHeader(uri string, header map[string]string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "PUT", data, header, false, param...)
}
func DeleteStructWithHeader(uri string, header map[string]string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "DELETE", data, header, false, param...)
}
func JsonPostStruct(uri string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "POST", data, nil, true, param...)
}
func JsonPutStruct(uri string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "PUT", data, nil, true, param...)
}
func JsonDeleteStruct(uri string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "DELETE", data, nil, true, param...)
}
func JsonPostStructWithHeader(uri string, header map[string]string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "POST", data, header, true, param...)
}
func JsonPutStructWithHeader(uri string, header map[string]string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "PUT", data, header, true, param...)
}
func JsonDeleteStructWithHeader(uri string, header map[string]string, data interface{}, param ...interface{}) error {
	return RequestStructWithHeader(uri, "DELETE", data, header, true, param...)
}
func GetStringWithTimeout(uri string, timeout time.Duration) (string, error) {
	client := gorequest.New().Timeout(timeout).Get(uri)
	_, body, errs := client.End()
	var err error = nil
	if errs != nil {
		err = errs[0]
		return body, err
	}
	return body, err
}
func GetStringWithTimeoutAndProxy(uri string, timeout time.Duration, proxy string) (string, error) {
	client := gorequest.New().Timeout(timeout).Proxy(proxy).Get(uri)
	_, body, errs := client.End()
	var err error = nil
	if errs != nil {
		err = errs[0]
		return body, err
	}
	return body, err
}
func JsonPostStructWithTimeout(uri string, data interface{}, timeout time.Duration, param ...interface{}) error {
	client := gorequest.New().Timeout(timeout).Post(uri)
	for _, v := range param {
		client.Send(v)
	}

	_, _, errs := client.EndStruct(data)
	var err error = nil
	if errs != nil {
		err = errs[0]
		return err
	}
	return err
}
func GetStringWithHeaderAndTimeout(uri string, header map[string]string, timeout time.Duration) (string, error) {
	client := gorequest.New().Timeout(timeout).Get(uri)
	for k, v := range header {
		client = client.Set(k, v)
	}

	_, body, errs := client.End()
	var err error = nil
	if errs != nil {
		err = errs[0]
		return body, err
	}

	return body, err
}

func PostStructWithFormUrlencoded(uri string, data interface{}, timeout time.Duration, param ...interface{}) error {
	client := gorequest.New().Timeout(timeout).Post(uri)
	for _, v := range param {
		client.Send(v)
	}
	client.Set("Content-Type", "application/x-www-form-urlencoded")
	client.TargetType = "form"
	_, _, errs := client.EndStruct(data)
	var err error = nil
	if errs != nil {
		err = errs[0]
		return err
	}
	return err
}
