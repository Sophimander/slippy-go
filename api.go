package slippygo

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/machinebox/graphql"
)

var slippiEndpoint string = "https://gql-gateway-dot-slippi.uc.r.appspot.com/graphql"
var slippiQuery string = `
query ($cc: String!) {
    getConnectCode(code: $cc){
        user {
            fbUid
            displayName
            status
            connectCode {
                code
            }
            activeSubscription {
                level
                hasGiftSub
                __typename
            }
            rankedNetplayProfile {
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
                    __typename
                }
                __typename
            }
            __typename
        }
    }
}`

var TimeData = []string{
	"MORS#762",
	"XATU#0",
	"RON#404",
	"ANZ#139",
	"SHOOP#0",
	"POLY#832",
	"GOOPY#1",
	"BEL#306",
	"JEO#807",
	"SO#0",
	"NAT#4713",
}

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

func (sc *SlippiClient) logf(format string, args ...interface{}) {
	sc.Log(fmt.Sprintf(format, args...))
}

func (sc *SlippiClient) Run(code string) (User, error) {
	sc.logf("Run(%v)", code)
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
	var user, err2 = convertUser(&resp.ConnectCode.User)
	if err2 != nil {
		return User{}, err2
	}
	return user, nil
}

type slippiResponse struct {
	ConnectCode getConnectCode `json:"getConnectCode"`
}

type getConnectCode struct {
	User sUser `json:"user"`
}

type User struct {
	Uid              string               `json:"fbUid"`
	DisplayName      string               `json:"displayName"`
	ConnectCode      ConnectCode          `json:"connectCode"`
	Status           string               `json:"status"`
	SubscriptionInfo SubscriptionInfo     `json:"activeSubscription"`
	RankedProfile    RankedNetplayProfile `json:"rankedNetplayProfile"`
}

type RankedNetplayProfile struct {
	Id                     int         `json:"id"`
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

type sUser struct {
	Uid              string                `json:"fbUid"`
	DisplayName      string                `json:"displayName"`
	ConnectCode      ConnectCode           `json:"connectCode"`
	Status           string                `json:"status"`
	SubscriptionInfo SubscriptionInfo      `json:"activeSubscription"`
	RankedProfile    sRankedNetplayProfile `json:"rankedNetplayProfile"`
}

type ConnectCode struct {
	Code string `json:"code"`
}

type SubscriptionInfo struct {
	Level  string `json:"level"`
	Gifted bool   `json:"hasGiftSub"`
}

type sRankedNetplayProfile struct {
	Id                     string       `json:"id"`
	Rating                 float64      `json:"ratingOrdinal"`
	RatingUpdateCount      int          `json:"ratingUpdateCount"`
	Wins                   int          `json:"wins"`
	Losses                 int          `json:"losses"`
	DailyGlobalPlacement   int          `json:"dailyGlobalPlacement"`
	DailyRegionalPlacement int          `json:"dailyRegionalPlacement"`
	Continent              string       `json:"continent"`
	Characters             []sCharacter `json:"characters"`
}

type sCharacter struct {
	Name      string `json:"character"`
	GameCount int    `json:"gameCount"`
}

func convertCharacter(sc *sCharacter) (Character, error) {
	char := Character{
		Name:      sc.Name,
		GameCount: sc.GameCount,
	}
	return char, nil
}

func convertRankedProfile(rp *sRankedNetplayProfile) (RankedNetplayProfile, error) {
	var Characters []Character
	for _, vchar := range rp.Characters {
		oldChar, err := convertCharacter(&vchar)
		if err != nil {
			return RankedNetplayProfile{}, err
		}
		Characters = append(Characters, oldChar)
	}
	newId, err := strconv.ParseInt(rp.Id, 0, 64)
	if err != nil {
		return RankedNetplayProfile{}, err
	}
	var rankedProfile = RankedNetplayProfile{
		Id:                     int(newId),
		Rating:                 rp.Rating,
		RatingUpdateCount:      rp.RatingUpdateCount,
		Wins:                   rp.Wins,
		Losses:                 rp.Losses,
		DailyGlobalPlacement:   rp.DailyGlobalPlacement,
		DailyRegionalPlacement: rp.DailyRegionalPlacement,
		Continent:              rp.Continent,
		Characters:             Characters,
	}
	return rankedProfile, nil
}

func convertUser(su *sUser) (User, error) {
	rankedProfile, err := convertRankedProfile(&su.RankedProfile)
	if err != nil {
		return User{}, err
	}
	var user = User{
		Uid:              su.Uid,
		DisplayName:      su.DisplayName,
		ConnectCode:      su.ConnectCode,
		Status:           su.Status,
		SubscriptionInfo: su.SubscriptionInfo,
		RankedProfile:    rankedProfile,
	}
	return user, nil
}

func validConnectCode(code string) bool {
	// This regex pattern does not cover this totally, because regexp lacks some perl features
	if len(code) < 3 && len(code) > 9 {
		return false
	}
	pattern := regexp.MustCompile("^[a-zA-Z]{1,7}#[0-9]{1,7}$")
	return pattern.Match([]byte(code))
}

func main() {
	fmt.Println("hi")
}
