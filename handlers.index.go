package main

import (
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/walkamongus/card-search/internal/hsapi"
	"github.com/walkamongus/card-search/internal/util"
)

// handleIndexPage contains business logic for
// handling requests to the root of the web app
func handleIndexPage(c *gin.Context) {
	// Instantiate a new API client if one does not exist
	if client == nil {
		client = hsapi.NewHSClient(viper.GetString("client-id"), viper.GetString("client-secret"), false)
	}

	// Pull druid and warlock cards with required attributes
	// from the API
	druids, err := client.SearchCards(map[string]string{
		"locale":   "en_US",
		"class":    "druid",
		"rarity":   "legendary",
		"manaCost": "7,8,9,10,11,12",
	})
	if err != nil {
		globalLog.Error(err, err.Error())
		AbortWithPage(500, err, c)
		return
	}
	warlocks, err := client.SearchCards(map[string]string{
		"locale":   "en_US",
		"class":    "warlock",
		"rarity":   "legendary",
		"manaCost": "7,8,9,10,11,12",
	})
	if err != nil {
		globalLog.Error(err, err.Error())
		AbortWithPage(500, err, c)
		return
	}

	// Retrieve card metadata from the API for
	// name, rarity, class, etc. ID lookups
	metadata, err := client.GetMetadata()
	if err != nil {
		globalLog.Error(err, err.Error())
		AbortWithPage(500, err, c)
		return
	}

	// Collect all retrieved cards
	cards := append(druids.Cards, warlocks.Cards...)

	// Parse cards and extract necessary attributes while
	// performing metadata lookups to retrieve text from IDs
	var cardsData []map[string]string
	for _, card := range cards {
		c := map[string]string{
			"id":     strconv.FormatInt(card.ID, 10),
			"name":   card.Name,
			"type":   util.GetName(card.CardTypeID, metadata.Types),
			"image":  card.Image,
			"rarity": util.GetRarityName(card.RarityID, metadata.Rarities),
			"set":    util.GetSetName(card.CardSetID, metadata.Sets),
			"class":  util.GetClassName(card.ClassID, metadata.Classes),
		}
		cardsData = append(cardsData, c)
	}

	// Shuffle the cards randomly and pick first 10
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cardsData), func(i, j int) { cardsData[i], cardsData[j] = cardsData[j], cardsData[i] })
	tableData := cardsData[0:10]

	// Sort chosen cards by their ID
	sort.Slice(tableData, func(i, j int) bool { return cardsData[i]["id"] < cardsData[j]["id"] })

	// Render the index view and pass in the filtered card data
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":   "Home Page",
			"payload": tableData,
		},
	)
}

func AbortWithPage(code int, err error, c *gin.Context) {
	c.HTML(
		code,
		"error.html",
		gin.H{
			"statusCode":    code,
			"statusMessage": err.Error(),
		},
	)
	c.Error(err)
	c.Abort()
}
