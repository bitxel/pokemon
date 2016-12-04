package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitxel/crawlee"
	"log"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

var (
	pokeid   = map[string]string{
		"4":"charmander",
		"5":"charmander2",
		"6":"charmander3",
		"38":"vulpix2",
		"63":"abra",
		"64":"abra2",
		"65":"abra3",
		"66":"machop",
		"67":"machop2",
		"68":"machop3",
		"83":"83dachong",
		"86":"seel",
		"87":"seal2",
		"88":"grimer",
		"89":"grimer2",
		"106":"hitmonlee",
		"107":"hitmonlee2",
		"108":"lickitung",
		"115":"115daishu",
		"122":"122mr mime",
		"128":"tauros",
		"131":"lapras",
		"132":"ditto",
		"137":"porygon",
		"138":"omanyte",
		"139":"omanyte2",
		"140":"kabuto",
		"141":"kabuto2",
		"143":"snorlax",
		"144":"144",
		"145":"145",
		"146":"146",
		"147":"dratini",
		"148":"dratini2",
		"149":"dratini3",
	}

	url      = "https://sgpokemap.com/query2.php?since=0&mons="
	position = []float64{1.318563, 103.774169}
	distance = 0.25
)

type RespStruct struct {
	//Meta []string `json:"meta"`
	Pokemons []Pokemon `json:"pokemons"`
}

type Pokemon struct {
	Attack  string `json:"attack"`
	Defense string `json:"defence"`
	Despawn string `json:"despawn"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
	Move1   string `json:"move1"`
	Move2   string `json:"move2"`
	Id      string `json:"pokemon_id"`
	Stamina string `json:"stamina"`
}

func sendmail(content string) {
	cmd := exec.Command("mail", "-s", "pokemon comming", "pokemon-sg@googlegroups.com")
	cmd.Stdin = strings.NewReader(content)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func getdist(lati1, long1, lati2, long2 float64) float64 {
	C := math.Sin(90.0-lati1)*math.Sin(90.0-lati2)*math.Cos(long1-long2) + math.Cos(90.0-lati1)*math.Cos(90.0-lati2)
	R := 6371.004
	Pi := math.Pi
	Dist := R * math.Acos(C) * Pi / 180
	return Dist
}

func check(mon *Pokemon) (float64, bool) {
	lat, _ := strconv.ParseFloat(mon.Lat, 10)
	lng, _ := strconv.ParseFloat(mon.Lng, 10)
	dist := getdist(lat, lng, position[0], position[1])
	if dist < distance {
		log.Printf("id:%s distance:%f lat:%f lng:%f", mon.Id, dist, lat, lng)
		return dist, true
	}
	return 0, false
}

func main() {
	var mons string
	for k, _ := range pokeid {
		mons += "," + k
	}
	mons = mons[1:]
	url += mons
	log.Printf("start to fetch url:%s", url)

	headers := make(http.Header)
	headers.Set("x-requested-with", "XMLHttpRequest")
	headers.Set("Referer", "https://sgpokemap.com/")
	resp, err := crawlee.GETX(url, headers)
	if err != nil {
		log.Printf("http err:%s", err)
		return
	}
	//log.Println(string(resp))
	res := &RespStruct{}
	err = json.Unmarshal(resp, res)
	if err != nil {
		log.Printf("unmarshal err:%s", err)
		return
	}
	pokemons := res.Pokemons

	var result string
	for _, mon := range pokemons {
		if dist, ok := check(&mon); ok {
			result += fmt.Sprintf("pokemon:%s %f km away, attack %s, defense %s, stamina %s\n",
				pokeid[mon.Id], dist, mon.Attack, mon.Defense, mon.Stamina)
		}
	}
	if len(result) > 0 {
		sendmail(result)
	}
}
