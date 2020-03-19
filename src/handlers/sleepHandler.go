package handlers

import (
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

type Feeling string

const (
	Sad 		Feeling = "sad"
	Happy 		Feeling = "happy"
	Tired		Feeling = "tired"
	Anxious 	Feeling = "anxious"
	Refreshed 	Feeling = "refreshed"
	Excited 	Feeling = "excited"
)

type sleepEvent struct {
	baseEvent 		BaseEvent
	sleepTime 		time.Time
	wakeupTime 		time.Time
	wakeupFeeling 	Feeling
	sleepFeeling 	Feeling
	qualityOfSleep 	int // [1, 10]
	alarmUsed 		bool
	ownBed 			bool
}

type SleepHandler struct {
	Session *mgo.Session
}

func (sh SleepHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

}

// Idea;
/* Have the generic event handler and pawn it off to the appropriate type manager

 */