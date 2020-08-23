package client

type Peer2PeerDiscovery struct {
	server   string // 线上服务器
	metadata string // 元数据
}

// 新建一个点对点的服务发现（其实就是没有服务发现）
func NewPeer2PeerDiscovery(server, metadata string) ServiceDiscovery {
	return &Peer2PeerDiscovery{
		server:   server,
		metadata: metadata,
	}
}

// 获取服务返回一个固定的地址
func (p *Peer2PeerDiscovery) GetServices() []*KVPair {
	return []*KVPair{
		{Key: p.server, Value: p.metadata},
	}
}

func (p *Peer2PeerDiscovery) WatchService() chan []*KVPair {
	return nil
}

func (p *Peer2PeerDiscovery) RemoveService(ch chan []*KVPair) {

}

func (p *Peer2PeerDiscovery) Clone(servicePath string) ServiceDiscovery {
	return p
}

func (p *Peer2PeerDiscovery) SetFilter(filter ServiceDiscoveryFilter) {

}

func (p *Peer2PeerDiscovery) Close() {

}
