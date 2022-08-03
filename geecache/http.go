package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"log"
	"net/http"
	"strings"
	"sync"
	pb "geecache/geecachepb"

	"github.com/golang/protobuf/proto"
)

const (
    defaultBasePath = "/_mmincache/"
    defaultReplicas = 50
)


// self记录自己的地址 主机名IP和端口 basepath为节点间的通信前缀
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string
	mu sync.Mutex
	peers *consistenthash.Map // 根据具体的key选择节点
	httpGetters map[string]*httpGetter // 客户端map
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 校验请求的格式是否正确
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 把获取的值写入response body里
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

// 实现PeerPicker
// 1 注册节点
func (p *HTTPPool)Set(peers ...string){
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas,nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _,peer := range peers{
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// 2 根据key返回peer
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool){
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key);peer != "" && peer !=p.self{
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

// 确保此类实现了该接口
var _ PeerPicker = (*HTTPPool)(nil)

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
	"%v%v/%v",
	h.baseURL,
	url.QueryEscape(in.GetGroup()),
	url.QueryEscape(in.GetKey()),
)
	res, err := http.Get(u)
	if err != nil {
	return err
}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
	return fmt.Errorf("server returned: %v", res.Status)
}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
	return fmt.Errorf("reading response body: %v", err)
}

	if err = proto.Unmarshal(bytes, out); err != nil {
	return fmt.Errorf("decoding response body: %v", err)
}

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)