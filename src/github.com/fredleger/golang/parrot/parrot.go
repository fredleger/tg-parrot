package parrot

import (
    "regexp"
    "time"
    "math/rand"
    "log"
)

type parrot struct {
    Name string
    PreferedSentence string
    AnswersPrefix string
    AnswersRate int
    LastRepeat time
    OnUseridShoulder int
    LastShoulderSwitch time
}

func NewParrot(name string, sentence string, prefix string, repeatrate int) *parrot {
    p := new(parrot)
    p.Name = name
    p.OnUseridShoulder = 0
    p.PreferedSentence = sentence
    p.AnswersPrefix = prefix
    p.AnswersRate = repeatrate
    return p
}

func (p parrot) Dump() {
    log.Println("struct: parrot")
    log.Println("Name: " + p.Name)
    log.Printf("OnUseridShoulder: %d", p.OnUseridShoulder)
    log.Println("PreferedSentence: " + p.PreferedSentence)
    log.Println("AnswersPrefix: " + p.AnswersPrefix)
    log.Printf("AnswersRate: %d" , p.AnswersRate)
}

func (p parrot) Repeat(input string) string {
    r, _ := regexp.Compile("[oO]")
    return p.AnswersPrefix + " " + r.ReplaceAllString(input, "oooo")
}

func (p parrot) WillRepeat() bool {

    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    r2 := r1.Intn(p.AnswersRate)

    log.Printf("r2:%d / rate:%d" , r2, p.AnswersRate-1)
    switch  {
    case r2 >= p.AnswersRate-1:
        return true
    default:
        return false
    }
}

func threesholdExeded(time lastOccurence, float frequency) bool {
    var timeDelta = lastOccurence - time.Now()
    var chance    = timeDelta*frequency
    var s1 = rand.NewSource(time.Now().UnixNano())
    var r1 = rand.New(s1)

    log.Printf("r1:%d / rate:%d" , r1, chance)
    switch  {
    case r1 >= chance:
        return true
    default:
        return false
    }
}