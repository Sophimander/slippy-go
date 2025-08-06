package slippygo

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/machinebox/graphql"
)

var slippiEndpoint string = "https://internal.slippi.gg"
var slippiQuery string = `
query UserProfilePageQuery($cc: String, $uid: String) {
  getUser(connectCode: $cc, fbUid: $uid) {
    ...userProfilePage
  }
}

fragment userProfilePage on User {
  fbUid
  displayName
  connectCode {
    code
  }
  status
  activeSubscription {
    level
    hasGiftSub
  }
  rankedNetplayProfile {
    ...profileFields
  }
  rankedNetplayProfileHistory {
    ...profileFields
    season {
      id
      startedAt
      endedAt
      name
      status
    }
  }
}

fragment profileFields on NetplayProfile {
  id
  ratingOrdinal
  ratingUpdateCount
  wins
  losses
  dailyGlobalPlacement
  dailyRegionalPlacement
  continent
  characters {
    character
    gameCount
  }
}
`

type SlippiClient struct {
	graphqlClient *graphql.Client
	Log           func(s string)
}

func NewClient() *SlippiClient {
	sc := &SlippiClient{
		graphqlClient: graphql.NewClient(slippiEndpoint),
		Log:           func(s string) {},
	}
	return sc
}

func (sc *SlippiClient) logf(format string, args ...any) {
	sc.Log(fmt.Sprintf(format, args...))
}

func (sc *SlippiClient) Run(code string) (User, error) {
	sc.logf("Run(%v)", code)
	if !validConnectCode(code) {
		return User{}, fmt.Errorf("Invalid connect code")
	}
	code = strings.ToUpper(code)
	req := graphql.NewRequest(slippiQuery)
	req.Var("cc", code)
	ctx := context.Background()
	var resp slippiResponse
	err := sc.graphqlClient.Run(ctx, req, &resp)
	if err != nil {
		sc.logf("Error: %v", err)
		return User{}, err
	}
	return resp.User, nil
}

type slippiResponse struct {
	User User `json:"getUser"`
}

type User struct {
	Uid              string                        `json:"fbUid"`
	DisplayName      string                        `json:"displayName"`
	ConnectCode      ConnectCode                   `json:"connectCode"`
	Status           string                        `json:"status"`
	SubscriptionInfo SubscriptionInfo              `json:"activeSubscription"`
	RankedProfile    RankedNetplayProfile          `json:"rankedNetplayProfile"`
	RankedHistory    []RankedNetplayProfileHistory `json:"rankedNetplayProfileHistory"`
}

type RankedNetplayProfile struct {
	Id                     string      `json:"id"`
	Rating                 float64     `json:"ratingOrdinal"`
	RatingUpdateCount      int         `json:"ratingUpdateCount"`
	Wins                   int         `json:"wins"`
	Losses                 int         `json:"losses"`
	DailyGlobalPlacement   int         `json:"dailyGlobalPlacement"`
	DailyRegionalPlacement int         `json:"dailyRegionalPlacement"`
	Continent              string      `json:"continent"`
	Characters             []Character `json:"characters"`
}

type Character struct {
	Name      string `json:"character"`
	GameCount int    `json:"gameCount"`
}

type Season struct {
	Id        string `json:"id"`
	StartedAt string `json:"startedAt"`
	EndedAt   string `json:"endedAt"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

type RankedNetplayProfileHistory struct {
	Id                     string      `json:"id"`
	Rating                 float64     `json:"ratingOrdinal"`
	RatingUpdateCount      int         `json:"ratingUpdateCount"`
	Wins                   int         `json:"wins"`
	Losses                 int         `json:"losses"`
	DailyGlobalPlacement   int         `json:"dailyGlobalPlacement"`
	DailyRegionalPlacement int         `json:"dailyRegionalPlacement"`
	Continent              string      `json:"continent"`
	Characters             []Character `json:"characters"`
	Seasons                Season      `json:"season"`
}

type ConnectCode struct {
	Code string `json:"code"`
}

type SubscriptionInfo struct {
	Level  string `json:"level"`
	Gifted bool   `json:"hasGiftSub"`
}

func validConnectCode(code string) bool {
	// This regex pattern does not cover this totally, because regexp lacks some perl features
	if len(code) < 3 && len(code) > 9 {
		return false
	}
	pattern := regexp.MustCompile("^[a-zA-Z]{1,7}#[0-9]{1,7}$")
	return pattern.Match([]byte(code))
}
