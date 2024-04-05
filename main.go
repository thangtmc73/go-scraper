package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", pageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	subPath := ""
	if len(page) != 0 && page != "0" {
		subPath = fmt.Sprint("page/", page)
	}
	fmt.Println(subPath)
	fmt.Println(fmt.Sprint("https://scrapeme.live/shop/", subPath))

	doc := loadHTML(fmt.Sprint("https://scrapeme.live/shop/", subPath))
	listPokemons := make([]Pokemon, len(doc.Find(".product").Nodes))

	doc.Find(".product").Each(func(i int, s *goquery.Selection) {
		listPokemons[i] = Pokemon{
			Name:  s.Find(".woocommerce-loop-product__title").Text(),
			Image: s.Find(".attachment-woocommerce_thumbnail").AttrOr("src", ""),
		}
		return
	})
	if len(listPokemons) == 0 {
		w.Write([]byte("No pokemon found in this page"))
		return
	}
	tmpl, err := template.New("list.tmpl").ParseFiles("list.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, ListData{Page: page, List: listPokemons})
	if err != nil {
		panic(err)
	}
}

type ListData struct {
	Page string
	List []Pokemon
}

type Pokemon struct {
	Name  string
	Image string
}

func loadHTML(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Status code err %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}
