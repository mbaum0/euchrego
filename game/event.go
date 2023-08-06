package game

import "fmt"

// an event is some input to the state machine

type Event interface {
	GetPlayer() *Player
}

type DefaultEvent struct {
	player *Player // the player who performed the event
	name   string
}

func (event DefaultEvent) GetPlayer() *Player {
	return event.player
}

type EndGameEvent struct {
	DefaultEvent
}

func NewEndGameEvent() EndGameEvent {
	event := EndGameEvent{}
	event.name = "END_GAME_EVENT"
	fmt.Println(event.name)
	return event
}

type EmptyEvent struct {
	DefaultEvent
}

func NewEmptyEvent() EmptyEvent {
	event := EmptyEvent{}
	event.name = "EMPTY_EVENT"
	return event
}

type DrawnCardEvent struct {
	DefaultEvent
	DrawnCard *Card
}

func NewDrawnCardEvent(card *Card, player *Player) DrawnCardEvent {
	event := DrawnCardEvent{}
	event.name = "DRAWN_CARD_EVENT"
	event.DrawnCard = card
	event.player = player
	fmt.Printf("%s: %s drew a %s\n", event.name, player.name, card.ToString())
	return event
}

type DrawnJackEvent struct {
	DefaultEvent
}

func NewDrawnJackEvent() DrawnJackEvent {
	event := DrawnJackEvent{}
	event.name = "DRAWN_JACK_EVENT"
	fmt.Printf("%s: Jack was drawn!\n", event.name)
	return event
}

type DealtCardsEvent struct {
	DefaultEvent
	DealtCards []*Card
}

func NewDealtCardsEvent(cards []*Card, player *Player) DealtCardsEvent {
	event := DealtCardsEvent{}
	event.name = "DEALT_CARDS_EVENT"
	event.DealtCards = cards
	event.player = player
	fmt.Printf("%s: %s was dealt %d cards\n", event.name, player.name, len(cards))
	return event
}

type FinishedDealingEvent struct {
	DefaultEvent
}

func NewFinishedDealingEvent() FinishedDealingEvent {
	event := FinishedDealingEvent{}
	event.name = "FINISHED_DEALING_EVENT"
	fmt.Println(event.name)

	return event
}

type TrumpSelectedEvent struct {
	DefaultEvent
	TrumpSuite Suite
	player     *Player
}

func NewTrumpSelectedEvent(suite Suite, player *Player) TrumpSelectedEvent {
	event := TrumpSelectedEvent{}
	event.name = "TRUMP_SELECTED_EVENT"
	event.player = player
	event.TrumpSuite = suite
	fmt.Printf("%s: %s selected %s as trump!\n", event.name, player.name, suite.ToString())
	return event
}

type TrumpPassedEvent struct {
	DefaultEvent
	PassedSuite Suite
	player      *Player
	anySuite    bool
}

func NewTrumpPassedEvent(passedSuite Suite, player *Player, anySuite bool) TrumpPassedEvent {
	event := TrumpPassedEvent{}
	event.name = "TRUMP_PASSED_EVENT"
	event.player = player
	event.anySuite = anySuite
	fmt.Printf("%s: %s passed picking trump\n", event.name, player.name)
	return event
}

type AskPlayerForTrumpEvent struct {
	DefaultEvent
	trumpCard *Card
	player    *Player
	anySuite  bool
}

func NewAskPlayerForTrumpEvent(trumpCard *Card, player *Player, anySuite bool) AskPlayerForTrumpEvent {
	event := AskPlayerForTrumpEvent{}
	event.name = "ASK_PLAYER_FOR_TRUMP_EVENT"
	event.player = player
	event.trumpCard = trumpCard
	event.anySuite = anySuite
	fmt.Printf("%s: %s requested to pick trump", event.name, player.name)
	if trumpCard != nil {
		fmt.Printf(". %s is shown", trumpCard.ToString())
	}
	if anySuite {
		fmt.Printf(". You may pick any suite")
	}
	fmt.Println()
	return event
}

type MisdealEvent struct {
	DefaultEvent
}

func NewMisdealEvent() MisdealEvent {
	event := MisdealEvent{}
	event.name = "MISDEAL_EVENT"
	fmt.Printf("%s: Misdeal!\n", event.name)
	return event
}
