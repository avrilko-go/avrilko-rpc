package client

type Weighted struct {
	Server          string // 服务地址
	Weight          int    // todo 权重(暂时无用)
	CurrentWeight   int    // 当前权重
	EffectiveWeight int    // 有效权重
}

// 参考nginx的平滑加权轮询算法
func nextWeighted(servers []*Weighted) *Weighted {
	total := 0
	var best *Weighted
	for _, w := range servers {
		if w == nil {
			continue
		}

		w.CurrentWeight += w.EffectiveWeight
		total += w.CurrentWeight

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}

	if best == nil {
		return nil
	}

	best.CurrentWeight -= total
	return best
}
