package time_line

import (
	"fmt"
	"sort"
	"time"

	"event_mgr/src/event_mgr"
	"event_mgr/src/utils"
)

var time_line_instance *TimeLine

type TimeLine struct {
	eventNodes event_mgr.EventMgrProxyList
	errorNodes event_mgr.EventMgrProxyList
}

func GetTimeLine() *TimeLine {
	if time_line_instance == nil {
		time_line_instance = &TimeLine{}
	}
	return time_line_instance
}

func (self *TimeLine) InsertToTimeNodes(nodes ...*event_mgr.EventMgrProxy) {
	now := utils.UtcTs()
	for _, node := range nodes {
		if err := node.RefreshEvents(now); err != nil {
			fmt.Println("==========", err)
			self.errorNodes = append(self.eventNodes, node)
		} else {
			self.eventNodes = append(self.eventNodes, node)
		}
	}
	self.sortTimeNodes()
}

func (self *TimeLine) Start(exit_channel chan interface{}) {
	break_flag := false
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(time.Second)
		for !break_flag {
			select {
			case <-exit_channel:
				break_flag = true
				break
			case <-ticker.C:
				self.tick()
			default:
				time.Sleep(time.Millisecond)
			}
		}

		done <- true
	}()

	fmt.Println("[TimeLine] awaiting signal")
	<-done
	fmt.Println("[TimeLine] exiting")
}

func (self *TimeLine) notifyObserver(node *event_mgr.EventMgrProxy) {
	fmt.Println("[notifyObserver]", node.TriggerEvents)
}

func (self *TimeLine) sortTimeNodes() { sort.Sort(self.eventNodes) }

func (self *TimeLine) tick() {
	resort_flag := false
	now_ts := utils.UtcTs()
	error_node_idx := []int{}
	for idx := 0; idx < len(self.eventNodes); idx++ {
		event_node := self.eventNodes[idx]
		if now_ts >= event_node.TriggerTs {
			self.notifyObserver(event_node) // get event_node to notify

			// add next event_node
			if err := event_node.RefreshEvents(now_ts); err != nil {
				error_node_idx = append(error_node_idx, idx)
			} else {
				resort_flag = true
			}
		} else {
			break
		}
	}

	// check error nodes
	recover_error_nodes := []int{}
	for idx := 0; idx < len(self.errorNodes); idx++ {
		error_node := self.errorNodes[idx]
		if err := error_node.RefreshEvents(now_ts); err == nil {
			resort_flag = true
			recover_error_nodes = append(recover_error_nodes, idx)
		}
	}

	// handle error nodes
	for idx := len(error_node_idx) - 1; idx >= 0; idx-- {
		error_node := self.eventNodes[error_node_idx[idx]]
		self.errorNodes = append(self.errorNodes, error_node)
		self.eventNodes = append(self.eventNodes[:idx], self.eventNodes[idx+1:]...)
	}

	// handle recover nodes
	for idx := len(recover_error_nodes) - 1; idx >= 0; idx-- {
		recover_node := self.errorNodes[error_node_idx[idx]]
		self.eventNodes = append(self.eventNodes, recover_node)
		self.errorNodes = append(self.errorNodes[:idx], self.errorNodes[idx+1:]...)
		resort_flag = true
	}

	// resort
	if resort_flag {
		self.sortTimeNodes()
	}
}
