package event_mgr

import (
	"errors"
	"fmt"
)

type EventMgrInterface interface {
	GetNextEventTs(int64) (bool, int64, []string)
	LoadEvents(string) error
}

type EventMgrProxyList []*EventMgrProxy

func (self EventMgrProxyList) Len() int { return len(self) }
func (self EventMgrProxyList) Less(i, j int) bool {
	return self[i].TriggerTs < self[j].TriggerTs
}
func (self EventMgrProxyList) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

type EventMgrProxy struct {
	owner         EventMgrInterface
	TriggerTs     int64
	TriggerEvents []string
}

func NewEventMgrProxy(owner EventMgrInterface) *EventMgrProxy {
	return &EventMgrProxy{
		owner:         owner,
		TriggerTs:     0,
		TriggerEvents: []string{},
	}
}

func (self *EventMgrProxy) RefreshEvents(now_ts int64) error {
	is_success, next_ts, trigger_events := self.owner.GetNextEventTs(now_ts)
	if !is_success {
		return errors.New("failed to get next events")
	}

	fmt.Println(is_success, next_ts, trigger_events)

	self.TriggerTs = next_ts
	self.TriggerEvents = trigger_events
	return nil
}
