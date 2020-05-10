package parrotlib

import (
	"encoding/json"
	"log"
	"math/rand"
	"regexp"
	"time"
)

type UserStats struct {
	UserId        uint64
	MessagesCount int64
	LastActivity  time.Time
}

type Parrot struct {
	Name  string
	Debug bool

	chats []int64

	preferedSentence string
	Users            map[string]*UserStats

	repeatPrefix      string
	repeatFrequency   float64
	RepeatMultiplier  float64
	RepeatAccumulator float64
	LastRepeat        time.Time

	onUserShoulder     string
	LastShoulderSwitch time.Time
}

// constructor
func NewParrot(name string, sentence string, repeatPrefix string, repeatFrequency float64, repeatMultiplier float64) Parrot {

	p := new(Parrot)
	p.Name = name
	p.onUserShoulder = ""
	p.preferedSentence = sentence
	p.repeatPrefix = repeatPrefix
	p.repeatFrequency = repeatFrequency
	p.Users = make(map[string]*UserStats)
	p.LastShoulderSwitch = time.Time{}
	p.LastRepeat = time.Time{}
	p.RepeatAccumulator = 0
	p.RepeatMultiplier = repeatMultiplier
	return *p
}

func (p Parrot) Dump() string {
	return p.ToString()
}

func (p *Parrot) ToString() string {
	obj, _ := json.Marshal(p)
	return string(obj)
}

// Say functions
func (p *Parrot) GetPreferedSentence() string {
	return p.preferedSentence
}

// Chats function
func (p *Parrot) AddChat(chatId int64) {

	if !sliceContainsInt64(p.chats, chatId) {
		p.chats = append(p.chats, chatId)
	}
}

func (p *Parrot) GetChats() []int64 {
	return p.chats
}

// Users list management functions
func (p *Parrot) AddUser(UserName string) bool {
	if user, ok := p.Users[UserName]; ok {
		user.LastActivity = time.Now()
	} else {
		p.Users[UserName] = new(UserStats)
		p.Users[UserName].LastActivity = time.Now()
		p.Users[UserName].MessagesCount = 0
	}
	return true
}

func (p *Parrot) RandomUser() (string, bool) {
	for k := range p.Users {
		return k, true
	}
	return "", false
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
	p.onUserShoulder = userId
}

func (p *Parrot) GetCurrentShoulder() string {
	if p.onUserShoulder != "" {
		return p.onUserShoulder
	} else {
		return "nobody"
	}
}

// utilities
func (p *Parrot) isThreesholdExeded() bool {
	var r1 = rand.New(rand.NewSource((time.Now().UnixNano()))).Float64()

	p.RepeatAccumulator += r1

	switch {
	case (p.RepeatAccumulator + r1) > (1-p.repeatFrequency)*p.RepeatMultiplier:
		return true
	default:
		return false
	}
}

// TODO : move that in another package
func sliceContainsInt64(s []int64, searchterm int64) bool {
	for _, k := range s {
		if k == searchterm {
			return true
		}
	}
	return false
}
