package game

type Player struct {
	hand        []*Card
	tricksTaken int
	name        string
	index       int
}

func InitPlayer(name string, index int) Player {
	player := Player{}
	player.hand = make([]*Card, 0)
	player.tricksTaken = 0
	player.name = name
	player.index = index

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

func (p *Player) ReturnCards() []*Card {
	var cards = make([]*Card, 0)
	cards = append(cards, p.hand...)
	p.hand = p.hand[:0]
	return cards
}
