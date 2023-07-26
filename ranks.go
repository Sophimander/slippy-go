package main

import (
	"math"
)

type Ranks struct {
	LowerBound float64
	UpperBound float64
	Name       string
}

var Grandmaster = Ranks{2191.75, math.Inf(1), "Grandmaster"}
var None = Ranks{0, 0, "None"}
var Pending = Ranks{0, 0, "Pending"}
var RankList = []Ranks{
	{0, 765.42, "Bronze 1"},
	{765.43, 913.71, "Bronze 2"},
	{913.72, 1054.86, "Bronze 3"},
	{1054.87, 1188.87, "Silver 1"},
	{1188.88, 1315.74, "Silver 2"},
	{1315.75, 1435.47, "Silver 3"},
	{1435.48, 1548.06, "Gold 1"},
	{1548.07, 1653.51, "Gold 2"},
	{1653.52, 1751.82, "Gold 3"},
	{1751.83, 1842.99, "Platinum 1"},
	{1843, 1927.02, "Platinum 2"},
	{1927.03, 2003.91, "Platinum 3"},
	{2003.92, 2073.66, "Diamond 1"},
	{2073.67, 2136.27, "Diamond 2"},
	{2136.28, 2191.74, "Diamond 3"},
	{2191.75, 2274.99, "Master 1"},
	{2275, 2350, "Master 2"},
	{2350, math.Inf(1), "Master 3"},
}

func get_rank(elo float64, daily_regional_placement int) Ranks {
	if elo > Grandmaster.LowerBound && daily_regional_placement > 0 {
		return Grandmaster
	}

	for _, rank := range RankList {
		if rank.LowerBound < elo && elo < rank.UpperBound {
			return rank
		}
	}

	return None
}

func (u *User) Rank() Ranks {
	if u.RankedProfile.Wins == 0 && u.RankedProfile.Losses == 0 {
		return None
	} else if (u.RankedProfile.Wins + u.RankedProfile.Losses) < 5 {
		return Pending
	} else {
		return get_rank(u.RankedProfile.Rating, u.RankedProfile.DailyGlobalPlacement)
	}
}
