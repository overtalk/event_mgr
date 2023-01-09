package event_mgr

import (
	"fmt"
	"testing"

	"event_mgr/src/utils"
)

func TestDailyEventMgrImpl(t *testing.T) {
	var event_mgr_interface EventMgrInterface
	event_mgr_interface = NewDailyEventMgrImpl(8)

	if err := event_mgr_interface.LoadEvents("F:/golang/event_mgr/config/daily_event.json"); err != nil {
		t.Error(err)
	}

	if err := event_mgr_interface.LoadEvents("F:/golang/event_mgr/config/daily_event1.json"); err == nil {
		t.Error("the config file is not exist, but load event success!")
	}

	has_next_trigger_ts, next_ts, _ := event_mgr_interface.GetNextEventTs(utils.UtcTs())
	fmt.Println("event_mgr_interface: ", event_mgr_interface, has_next_trigger_ts, next_ts)
}
