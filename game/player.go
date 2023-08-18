package game

type Player struct {
	hand         []*Card
	tricksTaken  int
	name         string
	index        int
	playedCard   *Card
	pointsEarned int
}

func InitPlayer(name string, index int) Player {
	player := Player{}
	player.hand = make([]*Card, 0)
	player.tricksTaken = 0
	player.pointsEarned = 0
	player.name = name
	player.index = index
	player.playedCard = nil

	return player
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetTricksTaken() int {
	return p.tricksTaken
}

func (p *Player) GiveCards(cards []*Card) {
	p.hand = append(p.hand, cards...)
}

func (p *Player) GiveCard(card *Card) {
	p.hand = append(p.hand, card)
}

func (p *Player) ReturnCards() []*Card {
	var cards = make([]*Card, 0)
	cards = append(cards, p.hand...)
	p.hand = p.hand[:0]
	return cards
}

// removes the card from the players hand and returns it
func (p *Player) ReturnCard(card *Card) *Card {
	for i, c := range p.hand {
		if c == card {
			p.hand = append(p.hand[:i], p.hand[i+1:]...)
			return c
		}
	}
	return nil
}
