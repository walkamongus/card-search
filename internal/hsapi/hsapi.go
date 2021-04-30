// Package hsapi provides API utilities for instantiating a
// Hearthstone API client and functions for interacting
// with the API.
package hsapi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// OAuthResponse represents the JSON API response
// when a new access token is requested.
type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Sub         string `json:"sub"`
	TokenType   string `json:"token_type"`
}

// HSClient represents a Hearthstone API client.
type HSClient struct {
	ClientID     string
	ClientSecret string
	TokenExpires time.Time
	*resty.Client
	*OAuthResponse
}

// CardSearchResult represents the top-level structure of
// the API response when a search is performed on cards.
type CardSearchResult struct {
	CardCount int64  `json:"cardCount"`
	Cards     []Card `json:"cards"`
	Page      int64  `json:"page"`
	PageCount int64  `json:"pageCount"`
}

// Card represents a single card retrieved from the API.
type Card struct {
	ArtistName    interface{}   `json:"artistName"`
	CardSetID     int64         `json:"cardSetId"`
	CardTypeID    int64         `json:"cardTypeId"`
	ClassID       int64         `json:"classId"`
	Collectible   int64         `json:"collectible"`
	CropImage     string        `json:"cropImage"`
	FlavorText    string        `json:"flavorText"`
	Health        int64         `json:"health"`
	ID            int64         `json:"id"`
	Image         string        `json:"image"`
	ImageGold     string        `json:"imageGold"`
	ManaCost      int64         `json:"manaCost"`
	MultiClassIds []interface{} `json:"multiClassIds"`
	Name          string        `json:"name"`
	ParentID      int64         `json:"parentId"`
	RarityID      int64         `json:"rarityId"`
	Slug          string        `json:"slug"`
	Text          string        `json:"text"`
}

// MetadataResult represents the top-level structure of the
// API response when all metadata is retrieved from the API.
type MetadataResult struct {
	ArenaIds           []int64           `json:"arenaIds"`
	CardBackCategories []GenericMetadata `json:"cardBackCategories"`
	Classes            []Class           `json:"classes"`
	FilterableFields   []string          `json:"filterableFields"`
	GameModes          []GenericMetadata `json:"gameModes"`
	Keywords           []Keyword         `json:"keywords"`
	MinionTypes        []GenericMetadata `json:"minionTypes"`
	NumericFields      []string          `json:"numericFields"`
	Rarities           []Rarity          `json:"rarities"`
	SetGroups          []SetGroup        `json:"setGroups"`
	Sets               []Set             `json:"sets"`
	SpellSchools       []GenericMetadata `json:"spellSchools"`
	Types              []GenericMetadata `json:"types"`
}

// Set represents the metadata for a particular card set.
type Set struct {
	AliasSetIds                 []int64 `json:"aliasSetIds"`
	CollectibleCount            int64   `json:"collectibleCount"`
	CollectibleRevealedCount    int64   `json:"collectibleRevealedCount"`
	ID                          int64   `json:"id"`
	Name                        string  `json:"name"`
	NonCollectibleCount         int64   `json:"nonCollectibleCount"`
	NonCollectibleRevealedCount int64   `json:"nonCollectibleRevealedCount"`
	Slug                        string  `json:"slug"`
	Type                        string  `json:"type"`
}

// Class represents the metadata for a particular card class.
type Class struct {
	AlternateHeroCardIds []int64 `json:"alternateHeroCardIds"`
	CardID               int64   `json:"cardId"`
	HeroPowerCardID      int64   `json:"heroPowerCardId"`
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	Slug                 string  `json:"slug"`
}

// SetGroup represents the metadata for a particular card set group.
type SetGroup struct {
	CardSets  []string `json:"cardSets"`
	Icon      string   `json:"icon"`
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Standard  bool     `json:"standard"`
	Svg       string   `json:"svg"`
	Year      int64    `json:"year"`
	YearRange string   `json:"yearRange"`
}

// Rarity represents the metadata for a particular card rarity.
type Rarity struct {
	CraftingCost []int64 `json:"craftingCost"`
	DustValue    []int64 `json:"dustValue"`
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
}

// Keyword represents the metadata for a card keyword.
type Keyword struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	RefText string `json:"refText"`
	Slug    string `json:"slug"`
	Text    string `json:"text"`
}

// GenericMetadata is a generic structure that can represent any
// metadata grouping returned that only includes these field.
// Examples of these metadata groups include: spell schools,
// types, minion types, and card back categories.
type GenericMetadata struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// NewHSClient instantiates a new API client with global client
// settings such as a base URL and retry settings. It does not
// instantiate any API access tokens.
func NewHSClient(id string, secret string, debug bool) *HSClient {
	client := resty.New()
	if debug {
		client.SetDebug(true)
	}
	client.SetHostURL("https://us.api.blizzard.com")
	client.SetRetryCount(3)
	client.AddRetryCondition(func(r *resty.Response, err error) bool {
		return r.StatusCode() == http.StatusTooManyRequests
	})
	return &HSClient{
		ClientID:     id,
		ClientSecret: secret,
		Client:       client,
	}
}

// GetToken is a receiver function that will generate and store an access
// token on the client. A new token is only generated if one does not exist
// or the current token is expired. This function may be invoked before any
// API interaction to ensure a valide token exists prior to API calls.
func (c *HSClient) GetToken() error {
	if c.TokenExpires.IsZero() || c.TokenExpires.Before(time.Now().Local()) {
		resp, err := c.R().
			SetBasicAuth(c.ClientID, c.ClientSecret).
			SetFormData(map[string]string{
				"grant_type": "client_credentials",
			}).
			SetResult(&OAuthResponse{}).
			Post("https://us.battle.net/oauth/token")

		if err != nil {
			return err
		}

		if resp.StatusCode() != 200 {
			return errors.New(fmt.Sprintf("Upstream API provider returned %d: %s", resp.StatusCode(), resp.String()))
		}

		c.OAuthResponse = resp.Result().(*OAuthResponse)
		c.TokenExpires = time.Now().Local().Add(time.Duration(c.ExpiresIn) * time.Second)
		c.SetAuthToken(c.AccessToken)
	}
	return nil
}

// SearchCards is a receiver function that allows clients
// to query the card search API. A map of arbitrary search
// parameters may be passed to customize the card search.
func (c *HSClient) SearchCards(opts map[string]string) (*CardSearchResult, error) {
	err := c.GetToken()
	if err != nil {
		return nil, err
	}
	resp, err := c.R().SetQueryParams(opts).SetResult(&CardSearchResult{}).Get("/hearthstone/cards")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprintf("Upstream API provider returned %d: %s", resp.StatusCode(), resp.String()))
	}
	return resp.Result().(*CardSearchResult), nil
}

// GetMetadata is a receiver function that pulls all
// card metadata from the metadata API endpoint.
func (c *HSClient) GetMetadata() (*MetadataResult, error) {
	err := c.GetToken()
	if err != nil {
		return nil, err
	}
	resp, err := c.R().
		SetQueryParams(map[string]string{"locale": "en_US"}).
		SetResult(&MetadataResult{}).
		Get("/hearthstone/metadata")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprintf("Upstream API provider returned %d: %s", resp.StatusCode(), resp.String()))
	}
	return resp.Result().(*MetadataResult), nil
}
