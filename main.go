package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/felixge/fgprof"
	"github.com/go-chi/chi"
)

func main() {
	http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// cpuを止める処理
		cpuIntensiveTask()
		// 良くわからない関数
		weirdFunction()
		// なんらかの network request
		slowNetworkRequest(r.Context())
	})

	fmt.Fprintf(os.Stdout, "start server\n")
	panic(http.ListenAndServe(":3000", r))
}

func slowNetworkRequest(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.jokes.one/jod", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-JokesOne-Api-Secret", "api_key")

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var joke Joke
	err = json.NewDecoder(resp.Body).Decode(&joke)
	if err != nil {
		panic(err)
	}
	fmt.Printf("title: %+v\n", joke.Contents.Jokes[0].Joke.Title)
	fmt.Printf("joke: %+v\n", joke.Contents.Jokes[0].Joke.Text)
}

func cpuIntensiveTask() {
	fmt.Println("sleeping...")
	time.Sleep(time.Second)
}

func weirdFunction() {
	fmt.Println("looping...")
	var row []int
	for i := 0; i < 1000*1000; i++ {
		row = append(row, i)
	}
	_ = row
}

type Joke struct {
	Contents struct {
		Copyright string `json:"copyright"`
		Jokes     []struct {
			Background  string `json:"background"`
			Category    string `json:"category"`
			Date        string `json:"date"`
			Description string `json:"description"`
			Joke        struct {
				Clean  string `json:"clean"`
				Date   string `json:"date"`
				ID     string `json:"id"`
				Lang   string `json:"lang"`
				Length string `json:"length"`
				Racial string `json:"racial"`
				Text   string `json:"text"`
				Title  string `json:"title"`
			} `json:"joke"`
			Language string `json:"language"`
		} `json:"jokes"`
	} `json:"contents"`
	Success struct {
		Total int64 `json:"total"`
	} `json:"success"`
}
