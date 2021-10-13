package load_balancer

import "errors"

// WeightedRound-Robin algorithm implementation

type wrrNode struct {
	weight        int64
	currentWeight int64
	value         interface{}
}

type WrrSelector struct {
	nodes []*wrrNode
}

func NewWrrSelector() *WrrSelector {
	return &WrrSelector{}
}

func (selector *WrrSelector) AddNode(weight int64, value interface{}) {
	if selector.nodes == nil {
		selector.nodes = []*wrrNode{}
	}
	selector.nodes = append(selector.nodes, &wrrNode{
		weight:        weight,
		currentWeight: weight,
		value:         value,
	})
}
func (selector *WrrSelector) Clear() {
	selector.nodes = []*wrrNode{}
}

func (selector *WrrSelector) Next() (result interface{}, err error) {
	var totalWeight int64 = 0
	var maxWeight int64 = -1
	var maxWeightItem *wrrNode = nil
	for _, item := range selector.nodes {
		totalWeight += item.weight
		item.currentWeight += item.weight
		if maxWeightItem == nil {
			maxWeight = item.currentWeight
			maxWeightItem = item
		} else if item.currentWeight > maxWeight {
			maxWeight = item.currentWeight
			maxWeightItem = item
		}
	}
	if maxWeightItem == nil {
		err = errors.New("can't find max weight item")
		return
	}
	maxWeightItem.currentWeight -= totalWeight
	result = maxWeightItem.value
	return
}
