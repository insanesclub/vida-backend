package search

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	local "github.com/maengsanha/vida-backend/model/kakao"
	"github.com/maengsanha/vida-backend/usecase/kakao"

	"github.com/PuerkitoBio/goquery"
	"github.com/maengsanha/vida-backend/usecase/instagram"
)

var total int

type detail struct {
	PlaceName    string `json:"place_name"`
	Phone        string `json:"phone"`
	Address      string `json:"road_address_name"`
	AveragePrice string `json:"average_price"`
}

// Details collects more detailed information about the results of the tag search on Instagram.
func Details(tag string) (detail_infos []detail) {
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
	fmt.Printf("collect %d locations from Instagram\n", len(loc_set))

	local_docs := make(chan local.Document)
	for location := range loc_set {
		syncer.Add(1)
		go parseMapInfo(location, local_docs, syncer)
	}
	go func() {
		syncer.Wait()
		close(local_docs)
	}()

	local_doc_set := make(map[string]local.Document)
	for doc := range local_docs {
		local_doc_set[doc.ID] = doc
	}
	fmt.Printf("collect %d documents from Kakao\n", len(local_doc_set))

	details := make(chan detail)
	for _, doc := range local_doc_set {
		syncer.Add(1)
		go parseDetails(doc, details, syncer)
	}
	go func() {
		syncer.Wait()
		close(details)
	}()

	for d := range details {
		detail_infos = append(detail_infos, d)
	}
	sort.Slice(detail_infos, func(i, j int) bool {
		return detail_infos[i].AveragePrice > detail_infos[j].AveragePrice
	})
	return
}

func parsePage(tag string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()

	parser := instagram.PageParserGenerator(tag)

	for total < 1e5 {
		page, err := parser()
		if err != nil {
			return
		}
		total += len(page.GraphQL.Hashtag.EdgeHashtagToMedia.Edges)

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

func parseMapInfo(location string, ch chan<- local.Document, syncer *sync.WaitGroup) {
	defer syncer.Done()

	map_info, err := kakao.LocalAPIParserGenerator(location)()
	if err != nil {
		return
	}
	if len(map_info.Documents) < 1 {
		return
	}
	ch <- map_info.Documents[0]
}

func parseDetails(doc local.Document, ch chan<- detail, syncer *sync.WaitGroup) {
	defer syncer.Done()

	price, err := averagePrice(doc.PlaceName, doc.ID)
	if err != nil {
		return
	}
	var price_str string
	if price > 0 {
		price_str = strconv.Itoa(price)
	}

	ch <- detail{
		PlaceName:    doc.PlaceName,
		Phone:        doc.Phone,
		Address:      doc.RoadAddressName,
		AveragePrice: price_str,
	}
}

func averagePrice(name, cid string) (average int, _ error) {
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
