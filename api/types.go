package api

import (
	"encoding/hex"
	"math/rand"
	"strings"
)

const (
	appVersion = "23.5.0.3873"
	apiVersion = "1.0.0"
	// sdkVersion          = "1.55.2.1"
	// deviceString        = "iPhone12,8"
	userAgentSynDir     = "Companion/23.4.2.3822 CFNetwork/978.0.7 Darwin/18.7.0"
	userAgentWebBrowser = "Mozilla/5.0 (iPhone; CPU iPhone OS 12_5_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148"
	eadpClientSecret    = "ltM2J0cYMRHJyR1wABxk2lgXkSI2OwetRFO7Yd8nC0Zf9MQQB2rTmsOPaEsARBdLqCC98XTZWcynlTM1"
	fifaGamePrefixURL   = "https://utas.mob.v1.fut.ea.com/ut/game/fifa23"
)

const (
	// funCaptchaURL       = "file://"
	funCaptchaSURL      = "ea-api.arkoselabs.com"
	funCaptchaPageURL   = "https://www.easports.com/fifa/ultimate-team/web-app/"
	funCaptchaPublicKey = "A4EECF77-AC87-8C8D-5754-BF882F72063B"
)

const (
	pinEventTimeFormat = "2006-01-02T15:04:05.999Z"
)

var (
	src = rand.New(rand.NewSource(1))
)

type Console uint64

const (
	PC Console = iota
	PS4
	PS5
	XBOXONE
	XBOXSERIES
)

func (c Console) ID() string {
	switch c {
	case PC:
		return "FFA23PCC"
	case PS4:
		return "FFA23PS4"
	case PS5:
		return "FFA23PS5"
	case XBOXONE:
		return "FFA23XBO"
	case XBOXSERIES:
		return "FFA23XSX"
	}
	return "UNKNOWN"
}

func (c Console) Name() string {
	switch c {
	case PC:
		return "pc"
	case PS4:
		return "ps4"
	case PS5:
		return "ps5"
	case XBOXONE:
		return "xbox_one"
	case XBOXSERIES:
		return "xbsx"
	}
	return "unk"
}

func getStringBetween(str, start, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	newStr := str[s+len(start):]
	e := strings.Index(newStr, end)
	if e == -1 {
		return ""
	}
	return newStr[:e]
}

func randHex(n int) string {
	b := make([]byte, n/2)
	if _, err := src.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

type PinEventsResponse struct {
	Status string `json:"status,omitempty"`
}

type CampaignActiveResponse struct {
	StartTime                int    `json:"startTime,omitempty"`
	EndTime                  int    `json:"endTime,omitempty"`
	ServerCrtTime            int    `json:"serverCrtTime,omitempty"`
	RemainingRewards         int    `json:"remainingRewards,omitempty"`
	UserXp                   int    `json:"userXp,omitempty"`
	CampaignBg               string `json:"campaignBg,omitempty"`
	CampaignWidgetBg         string `json:"campaignWidgetBg,omitempty"`
	Title                    string `json:"title,omitempty"`
	NextCampaignStartTime    int64  `json:"nextCampaignStartTime,omitempty"`
	HasPreviousCampaign      bool   `json:"hasPreviousCampaign,omitempty"`
	NeedsAutoClaim           bool   `json:"needsAutoClaim,omitempty"`
	NeedsMilestonesAutoClaim bool   `json:"needsMilestonesAutoClaim,omitempty"`
}

type UserMassInfoResponse struct {
	Errors                                Errors               `json:"errors,omitempty"`
	Settings                              Settings             `json:"settings,omitempty"`
	UserInfo                              UserInfo             `json:"userInfo,omitempty"`
	PurchasedItems                        PurchasedItems       `json:"purchasedItems,omitempty"`
	LoanPlayerClientData                  LoanPlayerClientData `json:"loanPlayerClientData,omitempty"`
	Squad                                 Squad                `json:"squad,omitempty"`
	ClubUser                              ClubUser             `json:"clubUser,omitempty"`
	ActiveMessages                        ActiveMessages       `json:"activeMessages,omitempty"`
	OnboardingClientData                  OnboardingClientData `json:"onboardingClientData,omitempty"`
	IsHighTierReturningUser               bool                 `json:"isHighTierReturningUser,omitempty"`
	IsPlayerPicksTemporaryStorageNotEmpty bool                 `json:"isPlayerPicksTemporaryStorageNotEmpty,omitempty"`
	LtpEventsHubData                      LtpEventsHubData     `json:"ltpEventsHubData,omitempty"`
}

type Errors struct {
}

type Configs struct {
	Value int    `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

type Settings struct {
	Configs []Configs `json:"configs,omitempty"`
}

type BidTokens struct {
}

type Currencies struct {
	Name       string `json:"name,omitempty"`
	Funds      int    `json:"funds,omitempty"`
	FinalFunds int    `json:"finalFunds,omitempty"`
}

type Actives struct {
	ID                 int64         `json:"id,omitempty"`
	Timestamp          int           `json:"timestamp,omitempty"`
	Formation          string        `json:"formation,omitempty"`
	Untradeable        bool          `json:"untradeable,omitempty"`
	AssetID            int           `json:"assetId,omitempty"`
	Rating             int           `json:"rating,omitempty"`
	ItemType           string        `json:"itemType,omitempty"`
	ResourceID         int           `json:"resourceId,omitempty"`
	Owners             int           `json:"owners,omitempty"`
	DiscardValue       int           `json:"discardValue,omitempty"`
	ItemState          string        `json:"itemState,omitempty"`
	Cardsubtypeid      int           `json:"cardsubtypeid,omitempty"`
	LastSalePrice      int           `json:"lastSalePrice,omitempty"`
	StatsList          []interface{} `json:"statsList,omitempty"`
	LifetimeStats      []interface{} `json:"lifetimeStats,omitempty"`
	AttributeList      []interface{} `json:"attributeList,omitempty"`
	Teamid             int           `json:"teamid,omitempty"`
	Rareflag           int           `json:"rareflag,omitempty"`
	LeagueID           int           `json:"leagueId,omitempty"`
	Pile               int           `json:"pile,omitempty"`
	Cardassetid        int           `json:"cardassetid,omitempty"`
	Value              int           `json:"value,omitempty"`
	Category           int           `json:"category,omitempty"`
	Manufacturer       string        `json:"manufacturer,omitempty"`
	Name               string        `json:"name,omitempty"`
	ResourceGameYear   int           `json:"resourceGameYear,omitempty"`
	Authenticity       bool          `json:"authenticity,omitempty"`
	Year               int           `json:"year,omitempty"`
	IsPlatformSpecific bool          `json:"isPlatformSpecific,omitempty"`
	Description        string        `json:"description,omitempty"`
	Biodescription     string        `json:"biodescription,omitempty"`
	Stadiumid          int           `json:"stadiumid,omitempty"`
	Capacity           int           `json:"capacity,omitempty"`
	TifoSupportType    int           `json:"tifoSupportType,omitempty"`
	TifoRestricted     bool          `json:"tifoRestricted,omitempty"`
	BannerRestricted   bool          `json:"bannerRestricted,omitempty"`
	BallRestricted     bool          `json:"ballRestricted,omitempty"`
	PreferredTime1     int           `json:"preferredTime1,omitempty"`
	PreferredTime2     int           `json:"preferredTime2,omitempty"`
	PreferredWeather   int           `json:"preferredWeather,omitempty"`
	Undiscardable      bool          `json:"undiscardable,omitempty"`
	Tier               int           `json:"tier,omitempty"`
	MyStadium          bool          `json:"myStadium,omitempty"`
	AttributeArray     []int         `json:"attributeArray,omitempty"`
	ShowCasePriority   int           `json:"showCasePriority,omitempty"`
	Weightrare         int           `json:"weightrare,omitempty"`
	Header             string        `json:"header,omitempty"`
	ChantsCount        int           `json:"chantsCount,omitempty"`
}

type Positions struct {
	Index int `json:"index,omitempty"`
	Value int `json:"value,omitempty"`
}

type Instructions struct {
	Index int `json:"index,omitempty"`
	Value int `json:"value,omitempty"`
}

type Styles struct {
	Index int `json:"index,omitempty"`
	Value int `json:"value,omitempty"`
}

type Tactics struct {
	SquadID        int            `json:"squadId,omitempty"`
	Tactic         string         `json:"tactic,omitempty"`
	LastUpdateTime int            `json:"lastUpdateTime,omitempty"`
	Formation      string         `json:"formation,omitempty"`
	Positions      []Positions    `json:"positions,omitempty"`
	Instructions   []Instructions `json:"instructions,omitempty"`
	Styles         []Styles       `json:"styles,omitempty"`
}

type Squad struct {
	ID         int           `json:"id,omitempty"`
	Valid      bool          `json:"valid,omitempty"`
	PersonaID  interface{}   `json:"personaId,omitempty"`
	Formation  string        `json:"formation,omitempty"`
	Rating     int           `json:"rating,omitempty"`
	Chemistry  int           `json:"chemistry,omitempty"`
	Manager    []interface{} `json:"manager,omitempty"`
	Players    []interface{} `json:"players,omitempty"`
	DreamSquad bool          `json:"dreamSquad,omitempty"`
	Changed    interface{}   `json:"changed,omitempty"`
	SquadName  string        `json:"squadName,omitempty"`
	StarRating int           `json:"starRating,omitempty"`
	Captain    interface{}   `json:"captain,omitempty"`
	Kicktakers []interface{} `json:"kicktakers,omitempty"`
	Actives    []interface{} `json:"actives,omitempty"`
	NewSquad   interface{}   `json:"newSquad,omitempty"`
	SquadType  string        `json:"squadType,omitempty"`
	Custom     interface{}   `json:"custom,omitempty"`
	Tactics    []Tactics     `json:"tactics,omitempty"`
}

type SquadListResponse struct {
	Squad         []Squad `json:"squad,omitempty"`
	ActiveSquadID int     `json:"activeSquadId,omitempty"`
}

type UnopenedPacks struct {
	PreOrderPacks  int `json:"preOrderPacks,omitempty"`
	RecoveredPacks int `json:"recoveredPacks,omitempty"`
}

type Reliability struct {
	Reliability         int `json:"reliability,omitempty"`
	StartedMatches      int `json:"startedMatches,omitempty"`
	FinishedMatches     int `json:"finishedMatches,omitempty"`
	MatchUnfinishedTime int `json:"matchUnfinishedTime,omitempty"`
}

type Feature struct {
	Trade  int `json:"trade,omitempty"`
	Rivals int `json:"rivals,omitempty"`
}

type UserInfo struct {
	PersonaID                  int64             `json:"personaId,omitempty"`
	ClubName                   string            `json:"clubName,omitempty"`
	ClubAbbr                   string            `json:"clubAbbr,omitempty"`
	Draw                       int               `json:"draw,omitempty"`
	Loss                       int               `json:"loss,omitempty"`
	Credits                    int               `json:"credits,omitempty"`
	BidTokens                  BidTokens         `json:"bidTokens,omitempty"`
	Currencies                 []Currencies      `json:"currencies,omitempty"`
	Trophies                   int               `json:"trophies,omitempty"`
	Won                        int               `json:"won,omitempty"`
	Actives                    []Actives         `json:"actives,omitempty"`
	Established                string            `json:"established,omitempty"`
	DivisionOffline            int               `json:"divisionOffline,omitempty"`
	DivisionOnline             int               `json:"divisionOnline,omitempty"`
	PersonaName                string            `json:"personaName,omitempty"`
	SquadList                  SquadListResponse `json:"squadList,omitempty"`
	UnopenedPacks              UnopenedPacks     `json:"unopenedPacks,omitempty"`
	Purchased                  bool              `json:"purchased,omitempty"`
	Reliability                Reliability       `json:"reliability,omitempty"`
	SeasonTicket               bool              `json:"seasonTicket,omitempty"`
	AccountCreatedPlatformName string            `json:"accountCreatedPlatformName,omitempty"`
	UnassignedPileSize         int               `json:"unassignedPileSize,omitempty"`
	Feature                    Feature           `json:"feature,omitempty"`
	SessionCoinsBankBalance    int               `json:"sessionCoinsBankBalance,omitempty"`
}

type ItemData struct {
	ID                      int64       `json:"id,omitempty"`
	Timestamp               int         `json:"timestamp,omitempty"`
	Formation               string      `json:"formation,omitempty"`
	Untradeable             bool        `json:"untradeable,omitempty"`
	AssetID                 int         `json:"assetId,omitempty"`
	Rating                  int         `json:"rating,omitempty"`
	ItemType                string      `json:"itemType,omitempty"`
	ResourceID              int         `json:"resourceId,omitempty"`
	Owners                  int         `json:"owners,omitempty"`
	DiscardValue            int         `json:"discardValue,omitempty"`
	ItemState               string      `json:"itemState,omitempty"`
	Cardsubtypeid           int         `json:"cardsubtypeid,omitempty"`
	LastSalePrice           int         `json:"lastSalePrice,omitempty"`
	InjuryType              string      `json:"injuryType,omitempty"`
	InjuryGames             int         `json:"injuryGames,omitempty"`
	PreferredPosition       string      `json:"preferredPosition,omitempty"`
	Contract                int         `json:"contract,omitempty"`
	Teamid                  int         `json:"teamid,omitempty"`
	Rareflag                int         `json:"rareflag,omitempty"`
	PlayStyle               int         `json:"playStyle,omitempty"`
	LeagueID                int         `json:"leagueId,omitempty"`
	Assists                 int         `json:"assists,omitempty"`
	LifetimeAssists         int         `json:"lifetimeAssists,omitempty"`
	LoyaltyBonus            int         `json:"loyaltyBonus,omitempty"`
	Pile                    interface{} `json:"pile,omitempty"`
	Nation                  int         `json:"nation,omitempty"`
	ResourceGameYear        int         `json:"resourceGameYear,omitempty"`
	AttributeArray          []int       `json:"attributeArray,omitempty"`
	StatsArray              []int       `json:"statsArray,omitempty"`
	LifetimeStatsArray      []int       `json:"lifetimeStatsArray,omitempty"`
	Skillmoves              int         `json:"skillmoves,omitempty"`
	Weakfootabilitytypecode int         `json:"weakfootabilitytypecode,omitempty"`
	Attackingworkrate       int         `json:"attackingworkrate,omitempty"`
	Defensiveworkrate       int         `json:"defensiveworkrate,omitempty"`
	Preferredfoot           int         `json:"preferredfoot,omitempty"`
	PossiblePositions       []string    `json:"possiblePositions,omitempty"`
	TradeID                 string      `json:"tradeId,omitempty"`
}

type PurchasedItems struct {
	ItemData []ItemData `json:"itemData,omitempty"`
}

type Entries struct {
	Value int `json:"value,omitempty"`
	Key   int `json:"key,omitempty"`
}

type LoanPlayerClientData struct {
	Entries []Entries `json:"entries,omitempty"`
}

type User struct {
	PersonaID int64  `json:"personaId,omitempty"`
	Persona   string `json:"persona,omitempty"`
}

type ClubUser struct {
	User []User `json:"user,omitempty"`
}

type ActiveMessages struct {
	ActiveMessage []interface{} `json:"activeMessage,omitempty"`
}

type OnboardingClientData struct {
	Entries []Entries `json:"entries,omitempty"`
}

type LtpEventsHubData struct {
	NeedsLTPEventAutoclaim bool `json:"needsLTPEventAutoclaim,omitempty"`
}

type PlayStatsResponse struct {
	UserOptIn         bool                `json:"userOptIn,omitempty"`
	PlayerHealthStats []PlayerHealthStats `json:"playerHealthStats,omitempty"`
	CrtServerTime     int64               `json:"crtServerTime,omitempty"`
}

type PlayerHealthStats struct {
	ID               string `json:"id,omitempty"`
	Current          int    `json:"current,omitempty"`
	Previous         int    `json:"previous,omitempty"`
	Total            int    `json:"total,omitempty"`
	Limit            int    `json:"limit,omitempty"`
	NextRolloverTime int    `json:"nextRolloverTime,omitempty"`
}

type LiveMessageResponse struct {
	MessageList       []MessageList `json:"messageList,omitempty"`
	MessagesRead      int           `json:"messagesRead,omitempty"`
	MessagesAvailable int           `json:"messagesAvailable,omitempty"`
	PromoUpdate       []interface{} `json:"promoUpdate,omitempty"`
}

type Trackurls struct {
}

type Attributes struct {
	Style           string `json:"style,omitempty"`
	Size            string `json:"size,omitempty"`
	Alignment       string `json:"alignment,omitempty"`
	Colour          string `json:"colour,omitempty"`
	HighlightColour string `json:"highlightColour,omitempty"`
	RenderType      string `json:"renderType,omitempty"`
	CountdownTime   string `json:"countdownTime,omitempty"`
	LocalID         string `json:"localId,omitempty"`
}

type Renders struct {
	Type       string     `json:"type,omitempty"`
	Name       string     `json:"name,omitempty"`
	Value      string     `json:"value,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
}

type MessageList struct {
	TrackingTag  string        `json:"trackingTag,omitempty"`
	Screen       string        `json:"screen,omitempty"`
	MessageID    int           `json:"messageId,omitempty"`
	Priority     int           `json:"priority,omitempty"`
	TmtLink      string        `json:"tmtLink,omitempty"`
	Trackurls    Trackurls     `json:"trackurls,omitempty"`
	Subtype      string        `json:"subtype,omitempty"`
	DoNotDisplay string        `json:"doNotDisplay,omitempty"`
	Renders      []Renders     `json:"renders,omitempty"`
	Promotions   []interface{} `json:"promotions,omitempty"`
}

type StadiumResponse struct {
	StadiumID         int     `json:"stadiumId,omitempty"`
	Name              string  `json:"name,omitempty"`
	LastUsedTierLevel int     `json:"lastUsedTierLevel,omitempty"`
	Tiers             []Tiers `json:"tiers,omitempty"`
}

type Slots struct {
	SlotID          int      `json:"slotId,omitempty"`
	Locked          bool     `json:"locked,omitempty"`
	StadiumArea     int      `json:"stadiumArea,omitempty"`
	SlotName        string   `json:"slotName,omitempty"`
	Home            bool     `json:"home,omitempty"`
	VanityTypes     []int    `json:"vanityTypes,omitempty"`
	Attributes      []int    `json:"attributes,omitempty"`
	ClientSlotID    int      `json:"clientSlotId,omitempty"`
	DefaultUnlocked bool     `json:"defaultUnlocked,omitempty"`
	ActiveSlot      bool     `json:"activeSlot,omitempty"`
	Emptiable       bool     `json:"emptiable,omitempty"`
	SortOrder       int      `json:"sortOrder,omitempty"`
	ItemData        ItemData `json:"itemData,omitempty"`
}

type Tiers struct {
	Tier  int     `json:"tier,omitempty"`
	Slots []Slots `json:"slots,omitempty"`
}

type AuctionInfo struct {
	TradeID           int64    `json:"tradeId,omitempty"`
	ItemData          ItemData `json:"itemData,omitempty"`
	TradeState        string   `json:"tradeState,omitempty"`
	BuyNowPrice       int      `json:"buyNowPrice,omitempty"`
	CurrentBid        int      `json:"currentBid,omitempty"`
	Offers            int      `json:"offers,omitempty"`
	Watched           bool     `json:"watched,omitempty"`
	BidState          string   `json:"bidState,omitempty"`
	StartingBid       int      `json:"startingBid,omitempty"`
	ConfidenceValue   int      `json:"confidenceValue,omitempty"`
	Expires           int      `json:"expires,omitempty"`
	SellerName        string   `json:"sellerName,omitempty"`
	SellerEstablished int      `json:"sellerEstablished,omitempty"`
	SellerID          int      `json:"sellerId,omitempty"`
	TradeOwner        bool     `json:"tradeOwner,omitempty"`
	CoinsProcessed    int      `json:"coinsProcessed,omitempty"`
	TradeIDStr        string   `json:"tradeIdStr,omitempty"`
}

type WatchlistResponse struct {
	Total       int           `json:"total,omitempty"`
	Credits     int           `json:"credits,omitempty"`
	AuctionInfo []AuctionInfo `json:"auctionInfo,omitempty"`
}

type TradepileResponse struct {
	Credits     int           `json:"credits,omitempty"`
	AuctionInfo []AuctionInfo `json:"auctionInfo,omitempty"`
	BidTokens   BidTokens     `json:"bidTokens,omitempty"`
}

type TradeStatusLiteResponse struct {
	AuctionInfo []AuctionInfo `json:"auctionInfo,omitempty"`
}

type TransferMarketResponse struct {
	AuctionInfo []AuctionInfo `json:"auctionInfo,omitempty"`
	BidTokens   BidTokens     `json:"bidTokens,omitempty"`
}

type SettingsResponse struct {
	Configs []struct {
		Value int    `json:"value,omitempty"`
		Type  string `json:"type,omitempty"`
	} `json:"configs,omitempty"`
}

type PriceLimit struct {
	Source   string `json:"source,omitempty"`
	DefID    int    `json:"defId,omitempty"`
	ItemID   int64  `json:"itemId,omitempty"`
	MinPrice int    `json:"minPrice,omitempty"`
	MaxPrice int    `json:"maxPrice,omitempty"`
}

type BidRequest struct {
	Bid int `json:"bid,omitempty"`
}

type BidResponse struct {
	Credits     int           `json:"credits,omitempty"`
	AuctionInfo []AuctionInfo `json:"auctionInfo,omitempty"`
	BidTokens   BidTokens     `json:"bidTokens,omitempty"`
	Currencies  []Currencies  `json:"currencies,omitempty"`
}

type ItemRequest struct {
	ItemData []ItemData `json:"itemData,omitempty"`
}

type ItemResponse struct {
	ItemData []ItemData `json:"itemData,omitempty"`
}

type AuctionHouseRequest struct {
	BuyNowPrice int      `json:"buyNowPrice,omitempty"`
	Duration    int      `json:"duration,omitempty"`
	ItemData    ItemData `json:"itemData,omitempty"`
	StartingBid int      `json:"startingBid,omitempty"`
}

type ObjectiveProgressList struct {
	ObjectiveID   int `json:"objectiveId,omitempty"`
	State         int `json:"state,omitempty"`
	ProgressCount int `json:"progressCount,omitempty"`
}

type ScmpGroupProgressList struct {
	GroupID               int                     `json:"groupId,omitempty"`
	State                 int                     `json:"state,omitempty"`
	ObjectiveProgressList []ObjectiveProgressList `json:"objectiveProgressList,omitempty"`
	GroupType             int                     `json:"groupType,omitempty"`
}

type ScmpCategoryProgressList struct {
	CategoryID            int                     `json:"categoryId,omitempty"`
	ScmpGroupProgressList []ScmpGroupProgressList `json:"scmpGroupProgressList,omitempty"`
}

type LearningGroupProgressList struct {
	CategoryID            int                     `json:"categoryId,omitempty"`
	ScmpGroupProgressList []ScmpGroupProgressList `json:"scmpGroupProgressList,omitempty"`
}

type DynamicObjectivesUpdates struct {
	NeedsGroupsRefresh        bool                        `json:"needsGroupsRefresh,omitempty"`
	ScmpCategoryProgressList  []ScmpCategoryProgressList  `json:"scmpCategoryProgressList,omitempty"`
	LearningGroupProgressList []LearningGroupProgressList `json:"learningGroupProgressList,omitempty"`
	NeedsAutoClaim            bool                        `json:"needsAutoClaim,omitempty"`
	NeedsMilestonesAutoClaim  bool                        `json:"needsMilestonesAutoClaim,omitempty"`
}

type AuctionHouseResponse struct {
	ID                       int64                    `json:"id,omitempty"`
	IDStr                    string                   `json:"idStr,omitempty"`
	DynamicObjectivesUpdates DynamicObjectivesUpdates `json:"dynamicObjectivesUpdates,omitempty"`
}

type TradeStatusResponse struct {
	Credits     int           `json:"credits,omitempty"`
	AuctionInfo []AuctionInfo `json:"auctionInfo,omitempty"`
	BidTokens   BidTokens     `json:"bidTokens,omitempty"`
	Currencies  []Currencies  `json:"currencies,omitempty"`
}

type TradeIDList struct {
	ID    int64  `json:"id,omitempty"`
	IDStr string `json:"idStr,omitempty"`
}

type RelistResponse struct {
	TradeIDList              []TradeIDList            `json:"tradeIdList,omitempty"`
	DynamicObjectivesUpdates DynamicObjectivesUpdates `json:"dynamicObjectivesUpdates,omitempty"`
}

type PinEventResponse struct {
	Status string `json:"status"`
}

type CaptchaDataResponse struct {
	Blob string `json:"blob"`
	Pk   string `json:"pk"`
}

type CaptchaValidateRequest struct {
	FunCaptchaToken string `json:"funCaptchaToken"`
}

type ConnectAuthResponse struct {
	Code string `json:"code"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      any    `json:"id_token"`
}

type UtAuthRequest struct {
	ClientVersion  int    `json:"clientVersion"`
	Ds             string `json:"ds"`
	GameSku        string `json:"gameSku"`
	Identification struct {
		AuthCode    string `json:"authCode"`
		RedirectURL string `json:"redirectUrl"`
	} `json:"identification"`
	IsReadOnly       bool   `json:"isReadOnly"`
	Locale           string `json:"locale"`
	Method           string `json:"method"`
	NucleusPersonaID int64  `json:"nucleusPersonaId"`
	PriorityLevel    int    `json:"priorityLevel"`
	SKU              string `json:"sku"`
}

type UtAuthResponse struct {
	Protocol       string `json:"protocol"`
	IPPort         string `json:"ipPort"`
	ServerTime     string `json:"serverTime"`
	LastOnlineTime string `json:"lastOnlineTime"`
	Sid            string `json:"sid"`
	PhishingToken  string `json:"phishingToken"`
}
