package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"event_mgr/src/event_mgr"
	"event_mgr/src/time_line"
)

func main() {
	server_cfg := flag.String("server", "", "server config")
	base_conf_path := flag.String("base_conf_path", "", "base config path")
	flag.Parse()

	if *server_cfg == "" {
		fmt.Println("[ERROR] server_cfg is empty")
		return
	}

	svr_cfg_path := filepath.Join(*base_conf_path, "server_cfg.json")
	if _, err := os.Stat(svr_cfg_path); err != nil {
		fmt.Println("[ERROR] server_cfg is absent")
		return
	}

	daily_event_cfg_path := filepath.Join(*base_conf_path, "daily_event.json")
	if _, err := os.Stat(daily_event_cfg_path); err != nil {
		fmt.Println("[ERROR] daily_event is absent")
		return
	}

	// parse config
	content, err := ioutil.ReadFile(svr_cfg_path)
	if err != nil {
		fmt.Printf("[ERROR] cfg_path(%s) failed to load file \n", svr_cfg_path)
		return
	}
	var json_data map[string]map[string]int64
	if err := json.Unmarshal([]byte(content), &json_data); err != nil {
		fmt.Printf("[ERROR] cfg_path(%s) failed to load file \n", svr_cfg_path)
		return
	}

	conf_details, is_exist := json_data[*server_cfg]
	if !is_exist {
		fmt.Printf("[ERROR] server_cfg(%s) is empty \n", *server_cfg)
		return
	}

	time_zone := conf_details["time_zone"]

	time_line := time_line.GetTimeLine()
	daily_event_mgr_impl_proxy := event_mgr.NewDailyEventMgrImplProxy(time_zone, daily_event_cfg_path)
	time_line.InsertToTimeNodes(daily_event_mgr_impl_proxy)

	exit_time_line_channel := make(chan interface{}, 1)
	go func() {
		time_line.Start(exit_time_line_channel)
	}()

	sig_channel := make(chan os.Signal)
	signal.Notify(sig_channel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig_channel
	fmt.Println("wait 3 second to end process")
	exit_time_line_channel <- 1
	time.Sleep(3 * time.Second)
}
