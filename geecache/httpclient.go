package geecache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpGetter struct {
	// 要访问的远程节点地址
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	// 获取返回值并转成 []bytes类型
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}

// 确保此类实现了该接口
var _ PeerGetter = (*httpGetter)(nil)
