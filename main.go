package main

import (
	//"fmt"
	"log"
	"net/http"
  "os"
  "io/ioutil"
  "encoding/json"
  "html/template"
)

var accountNames = [...]string {
  "AussieZulu",
  "BigRudy",
  "Shammah",
  "Cheergirl",
  "PacmanCloudy",
  "AznBeast42",
  "ExtraLettuce",
  "KBA_Allstar",
  "Millidavids",
  "Armadillyo",
  "Gregsaw",
  "Molpg",
  "Wolv3r1n3",
}

type Player struct {
  Name string
  Rating float64 // "Rating"
  BestRating float64 // "Best Rating"
  RoundsPlayed int // "Rounds Played"
  Wins int // "Wins"
  Losses int // "Losses"
  TopTens int // "Top 10s"
  Kills int // "Kills"
  Assists int // "Assists"
  Kd float64  // "K/D Ratio"
  HeadshotKills int // "Headshot Kills"
  LongestKill float64 // "Longest Kill"
  Revives int // "Revives"
  DamageDealt float64 // "Damage Dealt"
  KnockOuts int // "Knock Outs"
}

type PUBGResponse struct {
  SelectedRegion string `json:"selectedRegion"`
  DefaultSeason string `json:"defaultSeason"`
  Stats []struct {
    Region string `json:"Region"`
  	Season string `json:"Season"`
  	Match string `json:"Match"`
  	Stats []struct {
  		Label string `json:"label"`
  		ValueInt int `json:"valueInt"`
  		ValueDec float64 `json:"valueDec"`
  	} `json:"Stats"`
  } `json:"Stats"`
}

func main() {
  http.HandleFunc("/", handle)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

  client := &http.Client{}
  var players = [len(accountNames)]Player{}

  for index, accountName := range accountNames {
    players[index] = generatePlayer(accountName, client)
  }

  template, _ := template.ParseFiles("table.html")
  template.Execute(w, players)
}

func generatePlayer(name string, client *http.Client) Player {
  req, _ := http.NewRequest("GET", "https://pubgtracker.com/api/profile/pc/" + name, nil)
  req.Header.Add("TRN-API-KEY", os.Getenv("TRN_API_KEY"))
  resp, _ := client.Do(req)
  body, _ := ioutil.ReadAll(resp.Body)
  bodyBytes := []byte(body)

  var jsonResponse PUBGResponse
  json.Unmarshal(bodyBytes, &jsonResponse)

  player := Player{Name: name}

  for _, outerstat := range jsonResponse.Stats {
    if (outerstat.Region == jsonResponse.SelectedRegion && outerstat.Season == jsonResponse.DefaultSeason && outerstat.Match == "squad") {
      for _, innerstat := range outerstat.Stats {
        switch innerstat.Label {
          case "Rating":
            player.Rating = innerstat.ValueDec
          case "Best Rating":
            player.BestRating = innerstat.ValueDec
          case "Rounds Played":
            player.RoundsPlayed = innerstat.ValueInt
          case "Wins":
            player.Wins = innerstat.ValueInt
          case "Losses":
            player.Losses = innerstat.ValueInt
          case "Top 10s":
            player.TopTens = innerstat.ValueInt
          case "Kills":
            player.Kills = innerstat.ValueInt
          case "Assists":
            player.Assists = innerstat.ValueInt
          case "K/D Ratio":
            player.Kd = innerstat.ValueDec
          case "Headshot Kills":
            player.HeadshotKills = innerstat.ValueInt
          case "Longest Kill":
            player.LongestKill = innerstat.ValueDec
          case "Revives":
            player.Revives = innerstat.ValueInt
          case "Damage Dealt":
            player.DamageDealt = innerstat.ValueDec
          case "Knock Outs":
            player.KnockOuts = innerstat.ValueInt
          default:
            continue
        }
      }
    }
  }

  return player
}
