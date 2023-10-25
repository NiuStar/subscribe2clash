package req

import (
	"io"
	"net/http"
)

var Proxy string

func HttpGet(url string) (string, error) {
	/*reqIns := gorequest.New().Get(url).Timeout(time.Minute)
	if Proxy != "" {
		reqIns = reqIns.Proxy(Proxy)
	}
	_, body, errs := reqIns.End()
	if len(errs) > 0 {
		return "", errs[0]
	}*/

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	data, _ := io.ReadAll(resp.Body)
	return string(data), nil
}
