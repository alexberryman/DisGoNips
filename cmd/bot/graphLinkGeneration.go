package main

import (
	"DiscordGoTurnips/internal/turnips/generated-code"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type dailyPrice struct {
	DayOfWeek      int
	MorningPrice   int32
	AfternoonPrice int32
}

func linkServersCurrentPrices(m *discordgo.MessageCreate) (string, string) {
	var response string
	var reactionEmoji string
	q := turnips.New(db)
	ctx := context.Background()
	prices, err := q.GetWeeksPriceHistoryByServer(ctx, m.GuildID)
	if err != nil {
		log.Println("error fetching prices: ", err)
	}

	return buildPriceLinks(prices, response, reactionEmoji)
}

func linkUsersCurrentPrices(m *discordgo.MessageCreate) (string, string) {
	var response string
	var reactionEmoji string
	q := turnips.New(db)
	ctx := context.Background()
	prices, err := q.GetWeeksPriceHistoryByAccount(ctx, turnips.GetWeeksPriceHistoryByAccountParams{
		DiscordID: m.Author.ID,
		ServerID:  m.GuildID,
	})

	if err != nil {
		log.Println("error fetching prices: ", err)
	}

	data := make([]turnips.GetWeeksPriceHistoryByServerRow, 0)
	for _, v := range prices {
		p := turnips.GetWeeksPriceHistoryByServerRow(v)
		data = append(data, p)
	}

	return buildPriceLinks(data, response, reactionEmoji)
}

func linkAccountsPreviousPrices(m *discordgo.MessageCreate, offset int) (string, string) {
	var response string
	var reactionEmoji string
	q := turnips.New(db)
	ctx := context.Background()

	week, err := getCurrentWeek(m, q, ctx)

	prices, err := q.GetHistoricalWeekPriceHistoryByAccount(ctx, turnips.GetHistoricalWeekPriceHistoryByAccountParams{
		DiscordID: m.Author.ID,
		ServerID:  m.GuildID,
		Week:      int32(week - offset),
	})
	if err != nil {
		log.Println("error fetching prices: ", err)
	}

	data := make([]turnips.GetWeeksPriceHistoryByServerRow, 0)
	for _, v := range prices {
		p := turnips.GetWeeksPriceHistoryByServerRow(v)
		data = append(data, p)
	}

	return buildPriceLinks(data, response, reactionEmoji)
}

func linkServersPreviousPrices(m *discordgo.MessageCreate, offset int) (string, string) {
	var response string
	var reactionEmoji string
	q := turnips.New(db)
	ctx := context.Background()

	week, err := getCurrentWeek(m, q, ctx)

	prices, err := q.GetHistoricalWeekPriceHistoryByServer(ctx, turnips.GetHistoricalWeekPriceHistoryByServerParams{
		ServerID: m.GuildID,
		Week:     int32(week - offset),
	})
	if err != nil {
		log.Println("error fetching prices: ", err)
	}

	data := make([]turnips.GetWeeksPriceHistoryByServerRow, 0)
	for _, v := range prices {
		p := turnips.GetWeeksPriceHistoryByServerRow(v)
		data = append(data, p)
	}

	return buildPriceLinks(data, response, reactionEmoji)
}

func getCurrentWeek(m *discordgo.MessageCreate, q *turnips.Queries, ctx context.Context) (int, error) {
	account, _ := q.GetAccount(ctx, m.Author.ID)
	accountTimeZone, err := time.LoadLocation(account.TimeZone)
	localTime := time.Now().In(accountTimeZone)
	_, week := localTime.ISOWeek()
	return week, err
}

func buildPriceLinks(prices []turnips.GetWeeksPriceHistoryByServerRow, response string, reactionEmoji string) (string, string) {
	priceMap := make(map[string]map[string]dailyPrice)

	for _, value := range prices {
		wp := getEmptyWeeklyPrices()
		if _, ok := priceMap[value.Nickname]; ok {
			updateMorningOrAfterNoonPrice(value, priceMap)
		} else {
			priceMap[value.Nickname] = wp
			updateMorningOrAfterNoonPrice(value, priceMap)
		}
	}

	turnipLink := make(map[string]string)
	for nickname, prices := range priceMap {
		for _, d := range dayRange(Monday, Saturday) {
			if _, ok := turnipLink[nickname]; !ok {
				turnipLink[nickname] = ""
			}

			if prices[fmt.Sprint(d)].MorningPrice != 0 {
				turnipLink[nickname] += fmt.Sprintf("-%d", prices[fmt.Sprint(d)].MorningPrice)
			} else {
				turnipLink[nickname] += "-"
			}
			if prices[fmt.Sprint(d)].AfternoonPrice != 0 {
				turnipLink[nickname] += fmt.Sprintf("-%d", prices[fmt.Sprint(d)].AfternoonPrice)
			} else {
				turnipLink[nickname] += "-"
			}
		}
		response += fmt.Sprintf("%s: <https://ac-turnip.com/#%s>\n", nickname, turnipLink[nickname])
	}

	reactionEmoji = "🔗"
	return reactionEmoji, response
}

func dayRange(min, max Weekday) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = int(min) + i
	}
	return a
}

func getEmptyWeeklyPrices() map[string]dailyPrice {
	w := newWeeklyPrices()

	for _, d := range dayRange(Monday, Saturday) {
		dp := dailyPrice{
			DayOfWeek:      d,
			MorningPrice:   0,
			AfternoonPrice: 0,
		}
		w[fmt.Sprintf("%d", d)] = dp
	}
	return w
}

func newWeeklyPrices() map[string]dailyPrice {
	w := make(map[string]dailyPrice)
	return w
}

func updateMorningOrAfterNoonPrice(value turnips.GetWeeksPriceHistoryByServerRow, priceMap map[string]map[string]dailyPrice) {
	if value.AmPm == turnips.AmPmAm {
		tempPrice := priceMap[value.Nickname][fmt.Sprint(value.DayOfWeek)]
		tempPrice.MorningPrice = value.Price
		priceMap[value.Nickname][fmt.Sprint(value.DayOfWeek)] = tempPrice
	} else {
		tempPrice := priceMap[value.Nickname][fmt.Sprint(value.DayOfWeek)]
		tempPrice.AfternoonPrice = value.Price
		priceMap[value.Nickname][fmt.Sprint(value.DayOfWeek)] = tempPrice
	}
}