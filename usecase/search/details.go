package search

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/maengsanha/vida-backend/dataservice/kakao"

	"github.com/PuerkitoBio/goquery"
	"github.com/maengsanha/vida-backend/dataservice/instagram"
	insta "github.com/maengsanha/vida-backend/model/instagram"
	"github.com/maengsanha/vida-backend/model/kakao/local"
)

type detail struct {
	PlaceName       string `json:"place_name"`
	RoadAddressName string `json:"road_address_name"`
	Phone           string `json:"phone"`
	AveragePrice    string `json:"average_price"`
	price           int
}

// Details returns more detailed data using Kakao search results.
func Details(tag string) (details []detail) {
	loc_ch := make(chan string)
	syncer := new(sync.WaitGroup)
	syncer.Add(1)
	go func() {
		defer syncer.Done()
		parseLocation(tag, syncer, loc_ch)
	}()

	go func() {
		syncer.Wait()
		close(loc_ch)
	}()

	loc_set := make(map[string]struct{})
	for loc := range loc_ch {
		loc_set[loc] = struct{}{}
	}

	doc_map := make(map[string]local.Document)
	for loc := range loc_set {
		time.Sleep(0) // force context switching
		if resp, err := kakao.MapParserGenerator(loc)(); err == nil && len(resp.Documents) > 0 {
			doc_map[resp.Documents[0].ID] = resp.Documents[0]
		}
	}

	detail_ch := make(chan detail)
	for _, doc := range doc_map {
		syncer.Add(1)
		go func(doc local.Document) {
			defer syncer.Done()
			if price, err := parsePrice(doc.PlaceName, doc.ID); err == nil {
				detail_ch <- detail{
					PlaceName:       doc.PlaceName,
					RoadAddressName: doc.RoadAddressName,
					Phone:           doc.Phone,
					price:           price,
				}
			}
		}(doc)
	}

	go func() {
		syncer.Wait()
		close(detail_ch)
	}()

	for d := range detail_ch {
		if d.price > 0 {
			d.AveragePrice = strconv.Itoa(d.price)
		}
		details = append(details, d)
	}

	sort.Slice(details, func(i, j int) bool { return details[i].price > details[j].price })
	return
}

func parseLocation(tag string, syncer *sync.WaitGroup, ch chan<- string) {
	worker := instagram.PageParserGenerator(tag)
	for {
		page, err := worker()
		if err != nil {
			return
		}

		workers := make([]func() (insta.Post, error), len(page.GraphQL.Hashtag.EdgeHashtagToMedia.Edges))
		for i, edge := range page.GraphQL.Hashtag.EdgeHashtagToMedia.Edges { // register workers
			workers[i] = instagram.PostParserGenerator(edge.Node.Shortcode)
		}

		syncer.Add(len(workers))
		for _, worker := range workers {
			go func(worker func() (insta.Post, error)) {
				defer syncer.Done()
				if post, err := worker(); err == nil {
					ch <- post.FilterByLocation()
				}
			}(worker)
		}
	}
}

func parsePrice(name, cid string) (average int, _ error) {
	resp, err := http.Get(fmt.Sprintf("https://m.search.daum.net/kakao?w=poi&q=%s&cid=%s", name, cid))
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
