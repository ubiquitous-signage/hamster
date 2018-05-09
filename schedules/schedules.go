package schedules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"github.com/ubiquitous-signage/hamster/multiLanguageString"
	"github.com/ubiquitous-signage/hamster/panel"
	"github.com/ubiquitous-signage/hamster/util"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"gopkg.in/mgo.v2/bson"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, url.QueryEscape("hamster-calender-oauth.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func fetch() (panel.Panel, error) {
	var calendarID = viper.GetString("schedule.calendarID")
	ctx := context.Background()

	b, err := ioutil.ReadFile("schedules/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	t := time.Now()
	yy, mm, dd := t.Date()
	t_min := time.Date(yy, mm, dd, 0, 0, 0, 0, t.Location()).Format(time.RFC3339)
	t_max := time.Date(yy, mm, dd, 23, 59, 59, 0, t.Location()).Format(time.RFC3339)

	events, err := srv.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).TimeMin(t_min).TimeMax(t_max).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Printf("Unable to retrieve next 10 of the user's events. %v", err)
		return panel.Panel{}, err
	}

	schedules := &panel.Panel{
		Contents: []interface{}{},
	}
	schedules.Version = 0.0
	schedules.Type = "table"
	schedules.Title = *multiLanguageString.NewMultiLanguageString("本日の予定")
	schedules.Category = "internal"
	schedules.Date = time.Now()

	if len(events.Items) > 0 {
		fmt.Println("Today's events:")
		for _, item := range events.Items {
			var when string
			if item.Start.DateTime != "" {
				start, _ := time.Parse(time.RFC3339, item.Start.DateTime)
				end, _ := time.Parse(time.RFC3339, item.End.DateTime)
				when = start.Format("15:04") + "-" + end.Format("15:04")
			} else {
				when = "終日"
			}
			fmt.Printf("%s: %s\n", when, item.Summary)
			contentLine := []interface{}{
				*panel.NewStringContent(when),
				*panel.NewStringContent(item.Summary, true),
			}
			schedules.Contents = append(schedules.Contents.([]interface{}), contentLine)
		}
	} else {
		fmt.Printf("No events found today.\n")
		schedules.Contents = append([][]interface{}{{*panel.NewStringContent("本日、予定はありません。")}})
	}

	return *schedules, nil
}

func Run() {
	var startSecond = viper.GetDuration("schedule.startDelaySecond")
	var sleepSecond = viper.GetDuration("schedule.sleepSecond")

	time.Sleep(startSecond * time.Second)

	for {
		result, err := fetch()
		if err == nil {
			session, collection := util.GetPanel()
			defer session.Close()
			log.Println("Upsert schedules")
			collection.Upsert(
				bson.M{
					"version":  0.0,
					"type":     "table",
					"title.ja": "本日の予定",
					"category": "internal",
				},
				result,
			)
		} else {
			log.Println("Failed to get schedules from google calendar api: ", err.Error())
		}

		time.Sleep(sleepSecond * time.Second)
	}
}
