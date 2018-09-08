package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dukfaar/goUtils/relay"
	"github.com/dukfaar/leveBackend/leve"
)

type Resolver struct {
}

func (r *Resolver) Leves(ctx context.Context, args struct {
	First  *int32
	Last   *int32
	Before *string
	After  *string
}) (*leve.ConnectionResolver, error) {
	leveService := ctx.Value("leveService").(leve.Service)

	var totalChannel = make(chan int)
	go func() {
		var total, _ = leveService.Count()
		totalChannel <- total
	}()

	var levesChannel = make(chan []leve.Model)
	go func() {
		result, _ := leveService.List(args.First, args.Last, args.Before, args.After)
		levesChannel <- result
	}()

	var (
		start string
		end   string
	)

	var leves = <-levesChannel

	if len(leves) == 0 {
		start, end = "", ""
	} else {
		start, end = leves[0].ID.Hex(), leves[len(leves)-1].ID.Hex()
	}

	hasPreviousPageChannel, hasNextPageChannel := relay.GetHasPreviousAndNextPage(len(leves), start, end, leveService)

	return &leve.ConnectionResolver{
		Models: leves,
		ConnectionResolver: relay.ConnectionResolver{
			relay.Connection{
				Total:           int32(<-totalChannel),
				From:            start,
				To:              end,
				HasNextPage:     <-hasNextPageChannel,
				HasPreviousPage: <-hasPreviousPageChannel,
			},
		},
	}, nil
}

func (r *Resolver) Leve(ctx context.Context, args struct {
	Id string
}) (*leve.Resolver, error) {
	leveService := ctx.Value("leveService").(leve.Service)

	queryleve, err := leveService.FindByID(args.Id)

	if err == nil {
		return &leve.Resolver{
			Model: queryleve,
		}, nil
	}

	return nil, err
}

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

func (r *Resolver) ImportLeves(ctx context.Context) (*bool, error) {
	/*err := permission.Check(ctx, "mutation.ImportLeves")
	if err != nil {
		return nil, err
	}*/

	leveListResponse, err := http.Get("https://api.xivdb.com/leve?columns=id")

	if err != nil {
		fmt.Errorf("Error getting leve list: %v", err)
		return nil, err
	}
	defer leveListResponse.Body.Close()

	leveList := make([]LeveListResponse, 0)
	err = json.NewDecoder(leveListResponse.Body).Decode(&leveList)

	if err != nil {
		fmt.Errorf("Error reading leve list: %v", err)
		return nil, err
	}

	go func() {
		idChan := CreateLeveImporter(ctx.Value("leveService").(leve.Service))

		for _, leve := range leveList {
			idChan <- leve.ID
		}

		close(idChan)
	}()

	result := true
	return &result, nil
}
