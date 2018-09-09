package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dukfaar/leveBackend/leve"
)

type LeveResponse struct {
	ID         int32  `json:"id"`
	NameEN     string `json:"name_en"`
	Level      int32  `json:"class_level"`
	Class      string `json:"classjob_category"`
	Data       int32  `json:"data"`
	ExpReward  int32  `json:"exp_reward"`
	GilReward  int32  `json:"gil_reward"`
	ItemsTotal int32  `json:"items_total"`
}

func SetLeveData(leve *leve.Model, ID string) {
	leveResponse, err := http.Get("https://api.xivdb.com/leve/" + ID)

	if err != nil {
		fmt.Errorf("Error getting leve: %v", err)
		return
	}
	defer leveResponse.Body.Close()

	var leveData LeveResponse
	err = json.NewDecoder(leveResponse.Body).Decode(&leveData)

	if err != nil {
		fmt.Errorf("Error reading leve: %v", err)
		return
	}

	leve.Level = leveData.Level
	leve.Name = leveData.NameEN
	leve.XivdbID = ID
	leve.Xp = leveData.ExpReward
	leve.Class = leveData.Class
	leve.Gil = leveData.GilReward
}

func ImportLeve(ID int32, leveService leve.Service) {
	idString := strconv.FormatInt(int64(ID), 10)

	oldLeve, err := leveService.FindByXivdbID(idString)

	if err != nil {
		fmt.Errorf("Error getting leve with id %v: %v", ID, err)
	}

	if oldLeve == nil {
		var newLeve leve.Model
		SetLeveData(&newLeve, idString)
		_, err = leveService.Create(&newLeve)
	} else {
		SetLeveData(oldLeve, idString)
		_, err = leveService.Update(oldLeve.ID.Hex(), oldLeve)
	}

	if err != nil {
		fmt.Errorf("Error creating or updating leve %v: %v", idString, err)
	}
}

func CreateLeveImporter(leveService leve.Service) chan int32 {
	leveIDChan := make(chan int32)

	go func() {
		for {
			id, ok := <-leveIDChan
			ImportLeve(id, leveService)
			if !ok {
				return
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()

	return leveIDChan
}

type LeveListResponse struct {
	ID int32 `json:"id"`
}
