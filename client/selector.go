package client

import (
	"context"
	"github.com/edwingeng/doublejump"
	"github.com/valyala/fastrand"
	"net/url"
	"sort"
	"strconv"
)

type SelectFunc func(ctx context.Context, servicePath, serviceMethod string, args interface{}) string // 负载均衡函数

// 负载均衡算法实现接口
type Selector interface {
	Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string // 负载均衡选择函数
	UpdateServer(servers map[string]string)                                                 // 更新服务列表
}

func newSelector(selectMode SelectMode, servers map[string]string) Selector {
	switch selectMode {
	case RandomSelect: // 随机算法
		return newRandomSelector(servers)
	case RoundRobin: // 轮询
		return newRoundRobinSelector(servers)
	case WeightedRoundRobin:
		return newWeightRoundRobinSelector(servers)
	case ConsistentHash:
		return newConsistentHashSelector(servers)
	case SelectByUser:
		return nil
	default: // 默认也是使用随机
		return newRandomSelector(servers)
	}
}

// 随机负载均衡算法
type randomSelector struct {
	servers []string
}

func newRandomSelector(servers map[string]string) Selector {
	ss := make([]string, 0, len(servers))
	for k, _ := range servers {
		ss = append(ss, k)
	}
	return &randomSelector{servers: ss}
}

func (r *randomSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	ss := r.servers
	if len(ss) == 0 {
		return ""
	}

	i := fastrand.Uint32n(uint32(len(ss)))

	return ss[i]
}

func (r *randomSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k, _ := range servers {
		ss = append(ss, k)
	}
	r.servers = ss
}

// 轮询
type roundRobinSelector struct {
	servers []string
	i       int
}

func newRoundRobinSelector(servers map[string]string) Selector {
	ss := make([]string, 0, len(servers))
	for k, _ := range servers {
		ss = append(ss, k)
	}
	return &roundRobinSelector{servers: ss}
}

func (r *roundRobinSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	ss := r.servers
	if len(ss) == 0 {
		return ""
	}

	i := r.i
	i = i % len(ss)
	r.i = i + 1
	return ss[i]
}

func (r *roundRobinSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k, _ := range servers {
		ss = append(ss, k)
	}
	r.servers = ss
}

// 加权轮询
type weightRoundRobinSelector struct {
	servers []*Weighted
}

func newWeightRoundRobinSelector(servers map[string]string) Selector {
	return &weightRoundRobinSelector{
		servers: createdWeighted(servers),
	}
}

// 创建权重
func createdWeighted(servers map[string]string) []*Weighted {
	ss := make([]*Weighted, 0, len(servers))
	for k, metadata := range servers {
		w := &Weighted{
			Server:          k,
			Weight:          1,
			EffectiveWeight: 1,
		}
		if v, err := url.ParseQuery(metadata); err == nil {
			ww := v.Get("weight")
			if ww != "" {
				if weight, err := strconv.Atoi(ww); err == nil {
					w.Weight = weight
					w.EffectiveWeight = weight
				}
			}
		}
		ss = append(ss, w)
	}
	return ss
}

func (w *weightRoundRobinSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	ss := w.servers

	if len(ss) == 0 {
		return ""
	}

	s := nextWeighted(ss)
	if s == nil {
		return ""
	}
	return s.Server
}

func (w *weightRoundRobinSelector) UpdateServer(servers map[string]string) {
	s := createdWeighted(servers)
	w.servers = s
}

// 一致性hash
type consistentHashSelector struct {
	servers []string
	h       *doublejump.Hash
}

func newConsistentHashSelector(servers map[string]string) Selector {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}
	sort.Slice(ss, func(i, j int) bool { // 进行从小到大排序
		return ss[i] < ss[j]
	})

	return &consistentHashSelector{
		servers: ss,
		h:       doublejump.NewHash(),
	}
}

func (c *consistentHashSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	ss := c.servers
	if len(ss) == 0 {
		return ""
	}

	key := genKey(servicePath, serviceMethod, args)
	selected, _ := c.h.Get(key).(string)
	return selected
}

func (c *consistentHashSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		c.h.Add(k)
		ss = append(ss, k)
	}

	sort.Slice(ss, func(i, j int) bool { // 进行从小到大排序
		return ss[i] < ss[j]
	})
	for _, k := range c.servers {
		if servers[k] == "" {
			c.h.Remove(k)
		}
	}
	c.servers = ss
}
