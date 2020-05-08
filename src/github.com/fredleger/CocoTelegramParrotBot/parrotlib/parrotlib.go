package parrotlib

import (
	"log"
	"math/rand"
	"regexp"
	"time"

	//"log"
	"github.com/davecgh/go-spew/spew"
)

type Parrot struct {
	Name  string
	Debug bool

	PreferedSentence string
	Users            map[string]time.Time

	repeatPrefix      string
	repeatFrequency   float64
	RepeatMultiplier  float64
	RepeatAccumulator float64
	LastRepeat        time.Time

	OnUseridShoulder   string
	LastShoulderSwitch time.Time
}

// constructor
func NewParrot(name string, sentence string, repeatPrefix string, repeatFrequency float64, repeatMultiplier float64) Parrot {

	spew.Config.Indent = "\t"

	p := new(Parrot)
	p.Name = name
	p.OnUseridShoulder = ""
	p.PreferedSentence = sentence
	p.repeatPrefix = repeatPrefix
	p.repeatFrequency = repeatFrequency
	p.Users = make(map[string]time.Time)
	p.LastShoulderSwitch = time.Time{}
	p.LastRepeat = time.Time{}
	p.RepeatAccumulator = 0
	p.RepeatMultiplier = repeatMultiplier
	return *p
}

func (p Parrot) Dump() {
	spew.Dump(p)
}

func (p *Parrot) ToString() string {
	return spew.Sprintf("%#+v", p)
}

// Users list management functions
func (p *Parrot) AddUser(userid string) bool {
	p.Users[userid] = time.Now()
	return true
}

// repeat functions
func (p *Parrot) Repeat(input string) string {
	r, _ := regexp.Compile("[oO]")
	return p.repeatPrefix + " " + r.ReplaceAllString(input, "oooo")
}

func (p *Parrot) WillRepeat() bool {

	if p.isThreesholdExeded() {
		p.LastRepeat = time.Now()
		p.RepeatAccumulator = 0
		return true
	}
	return false
}

// shoulder functions
func (p *Parrot) SwitchShoulder(userId string) {
	log.Printf("I Switched to %v shoulder", userId)
}

// utilities
func (p *Parrot) isThreesholdExeded() bool {
	var r1 = rand.New(rand.NewSource((time.Now().UnixNano()))).Float64()

	if p.Debug {
		spew.Printf("frequency: %v\n", p.repeatFrequency)
		spew.Printf("treeshold: %v\n", (1-p.repeatFrequency)*p.RepeatMultiplier)
		spew.Printf("accumulator: %v\n", p.RepeatAccumulator)
		spew.Printf("random: %v\n", r1)
	}

	p.RepeatAccumulator += r1

	switch {
	case (p.RepeatAccumulator + r1) > (1-p.repeatFrequency)*p.RepeatMultiplier:
		return true
	default:
		return false
	}
}
