package main

import (
	"fmt"
	cr "github.com/go-openapi/runtime/client"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/naoina/genmai"
	apiclient "github.com/raa0121/go-kokoro-io/client"
	"github.com/raa0121/go-kokoro-io/client/bot"
	"github.com/raa0121/go-kokoro-io/models"
	"github.com/raa0121/onunume-bot/core"
	"github.com/yosssi/ace"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	tpl, err := ace.Load("layout", "index", &ace.Options{BaseDir: "view", DynamicReload: true})
	if err != nil {
		log.Fatal(err)
	}
	var results []core.Onunume
	if err = core.Db.Select(&results); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if err := tpl.Execute(w, map[string]interface{}{"results": results}); err != nil {
		log.Fatal(err)
	}
}

func kokoro(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "kokoro.io bot EntryPoint")
	} else if r.Method == http.MethodPost {
		length, err := strconv.Atoi(r.Header.Get("Content-Length"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body := make([]byte, length)
		length, err = r.Body.Read(body)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var post models.MessageEntity
		err = post.UnmarshalBinary(body[:length])
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		if result := onunumeCallback(post); result == "" {
			fmt.Fprintln(w, onunumeGreatestHit(post))
		} else {
			displayName := "ｵﾇﾇﾒ"
			expandEmbedContents := true
			nsfw := false
			param := &bot.PostV1BotChannelsChannelIDMessagesParams{
				ChannelID:           post.Channel.ID,
				DisplayName:         &displayName,
				ExpandEmbedContents: &expandEmbedContents,
				Message:             result,
				Nsfw:                &nsfw,
				HTTPClient:          &http.Client{Timeout: cr.DefaultTimeout},
			}
			param = param.WithTimeout(cr.DefaultTimeout)
			client := apiclient.New(httptransport.New("", "", nil), strfmt.Default)
			apiKeyHeaderAuth := httptransport.APIKeyAuth("X-Access-Token", "header", os.Getenv("API_KEY"))
			// resp, err := client.Bot.PostV1BotChannelsChannelIDMessages(param)
			log.Println(result)
			fmt.Fprintln(w, result)
		}
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func onunumeCallback(m models.MessageEntity) string {
	var result string
	r := regexp.MustCompile(`^/Onunume`)
	if r.MatchString(m.RawContent) {
		contents := strings.Split(m.RawContent, " ")
		switch contents[1] {
		case "add":
			result = core.Add(m)
		case "remove":
			result = core.Remove(m)
		case "help":
			result = core.Help(m)
		}

		return result
	}
	return ""
}

func main() {
	var err interface{}
	core.Db, err = genmai.New(&genmai.SQLite3Dialect{}, ":memory:")
	if err != nil {
		panic(err)
	}
	defer core.Db.Close()
	if err = core.Db.CreateTable(&core.Onunume{}); err != nil {
		panic(err)
	}
	sample := &core.Onunume{
		Name:      "hi",
		Regexp:    "^hi!$",
		Content:   "### hi",
		CreatedBy: "raa0121",
		Enable:    true,
		Expect:    "ねこ",
	}
	_, err = core.Db.Insert(sample)
	if err != nil {
		panic(err)
	}
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/kokoro", kokoro)
	http.ListenAndServe("0.0.0.0:8090", nil)
}

func onunumeGreatestHit(message models.MessageEntity) string {
	var results []core.Onunume
	if err := core.Db.Select(&results); err != nil {
		log.Fatal(err)
	}
	for _, val := range results {
		if val.Expect != message.Channel.ChannelName {
			continue
		}
		QContent := strconv.Quote(val.Regexp)
		r := regexp.MustCompile(QContent)
		if false == r.MatchString(message.RawContent) {
			continue
		}
		return val.Content
	}
	return ""
}
