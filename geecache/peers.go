package geecache

import pb "geecache/geecachepb"

// 节点选择接口
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 节点获取缓存值接口
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
