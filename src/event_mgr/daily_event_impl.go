package event_mgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"

	"event_mgr/src/utils"
)

type DailyEventDailySignCfgItem struct {
	sortSecond int64
	EventNames []string

	Hour   int64 `json:"hour"`
	Minute int64 `json:"minute"`
	Second int64 `json:"second"`
}

func (self *DailyEventDailySignCfgItem) init(event_name string) {
	self.EventNames = []string{event_name}
	self.sortSecond = self.toSecond()
}

func (self *DailyEventDailySignCfgItem) mergeEventNames(event_names []string) {
	self.EventNames = append(self.EventNames, event_names...)
}

func (self *DailyEventDailySignCfgItem) toString() string {
	return fmt.Sprintf("[DailyEventDailySignCfgItem] EventNames(%s) sort_second(%d) trigger_time(%d:%d:%d)",
		self.EventNames, self.sortSecond, self.Hour, self.Minute, self.Second)
}

func (self *DailyEventDailySignCfgItem) isValid() bool {
	if self.Hour < 0 || self.Minute < 0 || self.Second < 0 || self.Hour >= 24 || self.Minute >= 60 || self.Second >= 60 {
		return false
	}
	return true
}

func (self DailyEventDailySignCfgItem) toSecond() int64 {
	return self.Hour*3600 + self.Minute*60 + self.Second
}

type DailyEventDailySignCfgItemSlice []*DailyEventDailySignCfgItem

func (self DailyEventDailySignCfgItemSlice) Len() int { return len(self) }
func (self DailyEventDailySignCfgItemSlice) Less(i, j int) bool {
	return self[i].toSecond() < self[j].toSecond()
}
func (self DailyEventDailySignCfgItemSlice) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

// daily_event_manager
type DailyEventMgrImpl struct {
	timeOffset    int64
	sortEventList DailyEventDailySignCfgItemSlice
}

func NewDailyEventMgrImplProxy(time_offset int64, cfg_path string) *EventMgrProxy {
	impl := NewDailyEventMgrImpl(time_offset)
	impl.LoadEvents(cfg_path)
	return NewEventMgrProxy(impl)
}

func NewDailyEventMgrImpl(time_offset int64) *DailyEventMgrImpl {
	daily_mgr := DailyEventMgrImpl{timeOffset: time_offset}
	return &daily_mgr
}

func (self *DailyEventMgrImpl) LoadEvents(cfg_path string) error {
	// parse config
	content, err := ioutil.ReadFile(cfg_path)
	if err != nil {
		return err
	}
	var json_data map[string]*DailyEventDailySignCfgItem
	if err := json.Unmarshal([]byte(content), &json_data); err != nil {
		return err
	}

	// sort all configs
	var sort_list DailyEventDailySignCfgItemSlice
	for event_name, event_cfg := range json_data {
		if !event_cfg.isValid() {
			return errors.New(fmt.Sprintf("daily_event_cfg_error, event_name: %s\n", event_name))
		}
		event_cfg.init(event_name)
		sort_list = append(sort_list, event_cfg)
	}

	if len(sort_list) == 0 {
		return nil
	}

	// sort all events
	sort.Sort(sort_list)

	// merge event trigger in same time
	to_remove_idx_list := []int{}
	for i := len(sort_list) - 1; i >= 1; i-- {
		cur := sort_list[i]
		pre := sort_list[i-1]
		if cur.sortSecond == pre.sortSecond {
			pre.mergeEventNames(cur.EventNames)
			to_remove_idx_list = append(to_remove_idx_list, i)
		}
	}

	for _, to_remove_idx := range to_remove_idx_list {
		sort_list = append(sort_list[:to_remove_idx], sort_list[to_remove_idx+1:]...)
	}
	self.sortEventList = sort_list
	return nil
}

func (self *DailyEventMgrImpl) GetNextEventTs(cur_ts int64) (bool, int64, []string) {
	if len(self.sortEventList) == 0 {
		return false, -1, nil
	}

	time_with_tz := utils.TsToUtcTime(cur_ts, self.timeOffset)
	day_sec := int64(time_with_tz.Hour()*3600 + time_with_tz.Minute()*60 + time_with_tz.Second())

	next_trigger_idx := -1
	next_trigger_ts := int64(-1)
	ts := time_with_tz.Unix()
	for idx, event := range self.sortEventList {
		if event.sortSecond > day_sec {
			next_trigger_idx = idx
			next_trigger_ts = ts + (event.sortSecond - day_sec) - self.timeOffset*3600
			break
		}
	}

	if next_trigger_idx < 0 {
		sort_second := self.sortEventList[0].sortSecond
		diff := 24*60*60 - day_sec
		next_trigger_ts = cur_ts + sort_second + diff
		next_trigger_idx = 0
	}

	return true, next_trigger_ts, self.sortEventList[next_trigger_idx].EventNames
}
