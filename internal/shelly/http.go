package shelly

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"io"
	"net/http"
)

type ShellyStatus struct {
	WifiSta struct {
		Connected bool   `json:"connected"`
		Ssid      string `json:"ssid"`
		IP        string `json:"ip"`
		Rssi      int    `json:"rssi"`
	} `json:"wifi_sta"`
	Cloud struct {
		Enabled   bool `json:"enabled"`
		Connected bool `json:"connected"`
	} `json:"cloud"`
	Mqtt struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`
	Time          string `json:"time"`
	Unixtime      int    `json:"unixtime"`
	Serial        int    `json:"serial"`
	HasUpdate     bool   `json:"has_update"`
	Mac           string `json:"mac"`
	CfgChangedCnt int    `json:"cfg_changed_cnt"`
	ActionsStats  struct {
		Skipped int `json:"skipped"`
	} `json:"actions_stats"`
	Relays []struct {
		Ison           bool   `json:"ison"`
		HasTimer       bool   `json:"has_timer"`
		TimerStarted   int    `json:"timer_started"`
		TimerDuration  int    `json:"timer_duration"`
		TimerRemaining int    `json:"timer_remaining"`
		Overpower      bool   `json:"overpower"`
		IsValid        bool   `json:"is_valid"`
		Source         string `json:"source"`
	} `json:"relays"`
	Emeters []struct {
		Power         float64 `json:"power"`
		Reactive      float64 `json:"reactive"`
		Pf            float64 `json:"pf"`
		Voltage       float64 `json:"voltage"`
		IsValid       bool    `json:"is_valid"`
		Total         float64 `json:"total"`
		TotalReturned float64 `json:"total_returned"`
	} `json:"emeters"`
	Update struct {
		Status      string `json:"status"`
		HasUpdate   bool   `json:"has_update"`
		NewVersion  string `json:"new_version"`
		OldVersion  string `json:"old_version"`
		BetaVersion string `json:"beta_version"`
	} `json:"update"`
	RAMTotal int `json:"ram_total"`
	RAMFree  int `json:"ram_free"`
	FsSize   int `json:"fs_size"`
	FsFree   int `json:"fs_free"`
	Uptime   int `json:"uptime"`
}

func (c *Collector) getShellyStatus(host string) (*ShellyStatus, error) {
	url := fmt.Sprintf("http://%s/status", host)
	resp, err := http.Get(url)
	if err != nil {
		level.Error(c.Logger).Log("msg", "Error in http request: "+err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		level.Error(c.Logger).Log("msg", "Error reading body: "+err.Error())
		return nil, err
	}
	shellyStatus := ShellyStatus{}
	err = json.Unmarshal(body, &shellyStatus)
	if err != nil {
		level.Error(c.Logger).Log("msg", "Error Unmarshalling body: "+err.Error())
		return nil, err
	}
	return &shellyStatus, nil
}
