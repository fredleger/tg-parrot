package parrot

import (
    "regexp"
    "time"
    "math/rand"
    //"log"
    "github.com/davecgh/go-spew/spew"
)

type parrot struct {

    Name string
    Debug bool
    PreferedSentence string
    Users map[string]time.Time

    RepeatPrefix string
    RepeatFrequency float64
    LastRepeat time.Time

    OnUseridShoulder string
    LastShoulderSwitch time.Time

}

// constructor
func NewParrot(name string, sentence string, repeatprefix string, repeatfreq float64) parrot {

    spew.Config.Indent = "\t"

    p := new(parrot)
    p.Name              = name
    p.OnUseridShoulder  = ""
    p.PreferedSentence  = sentence
    p.RepeatPrefix      = repeatprefix
    p.RepeatFrequency   = repeatfreq
    p.Users             = make(map[string]time.Time)
    p.LastShoulderSwitch= time.Now()
    p.LastRepeat        = time.Now()
    return *p
}

func (p parrot) Dump() {
    spew.Dump(p)
}

func (p *parrot) ToString() string {
    return spew.Sprintf("%#+v", p)
}

// Users list management functions
func (p *parrot) AddUser(userid string) bool {
    p.Users[userid] = time.Now()
    return true
}

// repeat functions
func (p *parrot) Repeat(input string) string {
    r, _ := regexp.Compile("[oO]")
    return p.RepeatPrefix + " " + r.ReplaceAllString(input, "oooo")
}

func (p *parrot) WillRepeat() bool {
    var shouldI bool = p.threesholdExeded(p.LastRepeat, p.RepeatFrequency)

    //spew.Printf("LastRepeat     : %v\n", p.LastRepeat)
    //spew.Printf("RepeatFrequency: %v\n", p.RepeatFrequency)
    //spew.Printf("shouldI        : %v\n", shouldI)

    if shouldI {
        p.LastRepeat = time.Now()
        return true
    }
    return false
}

// shoulder functions

// utilities
func (p *parrot) threesholdExeded(lastOccurence time.Time, frequency float64) bool {
    var timeDelta = time.Now().Sub(lastOccurence)
    var chance    = timeDelta.Minutes()*frequency
    var s1 = rand.NewSource(time.Now().UnixNano())
    var r1 = rand.New(s1).Float64()

    //spew.Printf("lastOccurence: %v\ntime.Now(): %v\n" , lastOccurence, time.Now())
    //spew.Printf("timeDelta    : %v\nfrequency : %v\n" , timeDelta, frequency)
    //spew.Printf("chance       : %v\nrandom    : %v\n", chance, r1)

    switch  {
    case r1 <= chance:
        return true
    default:
        return false
    }
}