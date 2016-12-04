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
	//pokeid   = []int{4, 86, 108, 132, 137, 143, 42}
	pokeid   = []int{35}
	url      = "https://sgpokemap.com/query2.php?since=0&mons="
	position = []float64{1.1234, 103.1234}
	distance = 1.0
)

type RespStruct struct {
	//Meta []string `json:"meta"`
	Pokemons []Pokemon `json:"pokemons"`
}

type Pokemon struct {
	Attack  string `json:"attack"`
	Defense string `json:"defense"`
	Despawn string `json:"despawn"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
	Move1   string `json:"move1"`
	Move2   string `json:"move2"`
	Id      string `json:"pokemon_id"`
	Stamina string `json:"stamina"`
}

func sendmail(content string) {
	cmd := exec.Command("mail", "-s", "pokemon comming", "email@email.com", "email@email.com")
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
	mons := strconv.Itoa(pokeid[0])
	for _, v := range pokeid {
		mons += "," + strconv.Itoa(v)
	}
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
			result += fmt.Sprintf("pokemon id:%s %f km away, attack %s, defense %s, stamina %s\n",
				mon.Id, dist, mon.Attack, mon.Defense, mon.Stamina)
		}
	}
	if len(result) > 0 {
		sendmail(result)
	}
}
