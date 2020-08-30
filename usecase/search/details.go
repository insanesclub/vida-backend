package search

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/maengsanha/vida-backend/usecase/kakao"

	"github.com/PuerkitoBio/goquery"
	insta "github.com/maengsanha/vida-backend/model/instagram"
	"github.com/maengsanha/vida-backend/usecase/instagram"
)

type detail struct {
	PlaceName    string `json:"place_name"`
	Phone        string `json:"phone"`
	Address      string `json:"road_address_name"`
	AveragePrice int    `json:"average_price"`
}

// Details collects more detailed information about the results of the tag search on Instagram.
func Details(tag string) []detail {
	locations := make(chan string)
	syncer := new(sync.WaitGroup)
	syncer.Add(1)
	go parsePage(tag, locations, syncer)
	go func() {
		syncer.Wait()
		close(locations)
	}()
	loc_set := make(map[string]struct{})
	for location := range locations {
		loc_set[location] = struct{}{}
	}
	details := make(chan detail)
	for location := range loc_set {
		syncer.Add(1)
		fmt.Printf("location: %s\n", location)
		go parseDetails(location, details, syncer)
	}
	go func() {
		syncer.Wait()
		close(details)
	}()
	var ds []detail
	for d := range details {
		ds = append(ds, d)
	}
	sort.Slice(ds, func(i, j int) bool {
		return ds[i].AveragePrice > ds[j].AveragePrice
	})
	return ds
}

func parsePage(tag string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()

	var (
		err    error
		parser = instagram.PageParserGenerator(tag)
	)

	for err == nil {
		var page insta.Tagpage
		page, err = parser()
		if err != nil {
			return
		}
		for _, edge := range page.GraphQL.Hashtag.EdgeHashtagToMedia.Edges {
			syncer.Add(1)
			go parseLocation(edge.Node.Shortcode, ch, syncer)
		}
	}
}

func parseLocation(shortcode string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()

	post, err := instagram.PostParserGenerator(shortcode)()
	if err != nil {
		return
	}
	loc := post.FilterByLocation()
	if len(loc) > 0 {
		ch <- loc
	}
}

func parseDetails(location string, ch chan<- detail, syncer *sync.WaitGroup) {
	defer syncer.Done()

	parser := kakao.LocalAPIParserGenerator(location)
	info, err := parser()
	if err != nil {
		return
	}
	if len(info.Documents) > 0 {
		doc := info.Documents[0]
		price, err := averagePrice(doc.PlaceName, doc.ID)
		if err != nil {
			return
		}
		ch <- detail{
			PlaceName:    doc.PlaceName,
			Phone:        doc.Phone,
			Address:      doc.RoadAddressName,
			AveragePrice: price,
		}
	}
}

func averagePrice(name, cid string) (average int, _ error) {
	resp, err := http.Get(fmt.Sprintf("https://m.search.daum.net/kakao?DA=SH2&w=poi&q=%s&cid=%s#&linked=true", name, cid))
	if err != nil {
		return average, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return average, err
	}
	prices := doc.Find("span.txt_price").
		Map(func(_ int, s *goquery.Selection) string { return s.Text() })
	if len(prices) < 1 {
		return average, nil
	}
	for _, s := range prices {
		price, err := strconv.Atoi(strings.ReplaceAll(s, ",", ""))
		if err != nil {
			return
		}
		average += price
	}
	average /= len(prices)
	return average, nil
}
