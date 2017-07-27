package parrot

import (
    "regexp"
    "time"
    "math/rand"
    "github.com/davecgh/go-spew/spew"
)

type parrot struct {
    Name string
    PreferedSentence string
    RepeatPrefix string
    RepeatFrequency float64
    LastRepeat time.Time
    OnUseridShoulder string
    LastShoulderSwitch time.Time
    Users map[string]time.Time
}

// constructor
func NewParrot(name string, sentence string, repeatprefix string, repeatfreq float64) *parrot {
    p := new(parrot)
    p.Name              = name
    p.OnUseridShoulder  = ""
    p.PreferedSentence  = sentence
    p.RepeatPrefix      = repeatprefix
    p.RepeatFrequency   = repeatfreq
    p.Users             = make(map[string]time.Time)
    p.LastShoulderSwitch= time.Now()
    p.LastRepeat        = time.Now()
    return p
}

func (p parrot) Dump() {
    spew.Dump(p)
}

// Users list management functions
func (p parrot) AddUser(userid string) bool {
    p.Users[userid] = time.Now()
    return true
}

// repeat functions
func (p parrot) Repeat(input string) string {
    r, _ := regexp.Compile("[oO]")
    return p.RepeatPrefix + " " + r.ReplaceAllString(input, "oooo")
}

func (p parrot) WillRepeat() bool {
    var shouldI bool = p.threesholdExeded(p.LastRepeat, p.RepeatFrequency)
    spew.Printf("LastRepeat: %v, RepeatFrequency: %v, shouldI: %v\n", p.LastRepeat, p.RepeatFrequency, shouldI)
    if shouldI {
        return true
    }
    return false
}

// shoulder functions


// utilities
func (p parrot) threesholdExeded(lastOccurence time.Time, frequency float64) bool {
    var timeDelta = time.Now().Sub(lastOccurence)
    var chance    = timeDelta.Minutes()*frequency
    var s1 = rand.NewSource(time.Now().UnixNano())
    var r1 = rand.New(s1).Float64()

    spew.Printf("timeDelta: %v, frequency: %v\n" , timeDelta, frequency)
    spew.Printf("chance: %v, random: %v\n", chance, r1)
    switch  {
    case r1 <= chance:
        return true
    default:
        return false
    }
}