package game

import "github.com/mbaum0/euchrego/card"

type Player struct {
	hand         []card.Card
	tricksTaken  int
	name         string
	index        int
	playedCard   card.Card
	pointsEarned int
}

func InitPlayer(name string, index int) *Player {
	player := Player{}
	player.hand = make([]card.Card, 0)
	player.tricksTaken = 0
	player.pointsEarned = 0
	player.name = name
	player.index = index
	player.playedCard = card.EmptyCard()

	return &player
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetTricksTaken() int {
	return p.tricksTaken
}

func (p *Player) GiveCards(cards []card.Card) {
	p.hand = append(p.hand, cards...)
}

func (p *Player) GiveCard(c card.Card) {
	p.hand = append(p.hand, c)
}

func (p *Player) ReturnCards() []card.Card {
	var cards = make([]card.Card, 0)
	cards = append(cards, p.hand...)
	p.hand = p.hand[:0]
	return cards
}

// removes the card from the players hand and returns it
func (p *Player) ReturnCard(rc card.Card) card.Card {
	for i, c := range p.hand {
		if c == rc {
			p.hand = append(p.hand[:i], p.hand[i+1:]...)
			return c
		}
	}
	return card.EmptyCard()
}
