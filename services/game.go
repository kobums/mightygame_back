package services

import (
	"errors"
	"fmt"
	"math/rand"
	"mighty/global"
	"mighty/models"
	"strconv"
	"strings"
)

type Mode int

const (
	Bidding Mode = iota
)

type Type int

const (
	_ Type = iota

	Spade
	Heart
	Diamond
	Clover
	Joker
	Mighty
	JokerCall
	NoTrump
	NoType
)

type FriendType int

const (
	CardFriend FriendType = iota
	FirstFriend
	NoFriend
)

type Card struct {
	Type   Type
	Number int
}

type User struct {
	Id         int64
	Name       string
	Cards      []Card
	Draw       *Card
	PointCards []Card
	Chip       int
	Pass       bool
}

type Friend struct {
	Id   int
	Type FriendType
	Card Card
}

type Game struct {
	Id           int64
	Name         string
	Users        []User
	Master       int
	Friend       Friend
	Trump        Type
	TableCards   []Card
	Bidding      int
	BiddingType  Type
	BiddingUser  int
	BiddingCount int
	Winner       int
	DrawCount    int
	JokerCall    bool
	Round        int
}

func remove(slice []Card, s int) []Card {
	return append(slice[:s], slice[s+1:]...)
}

func getType(str string) Type {
	if str == "Spade" || str == "S" || str == "s" {
		return Spade
	} else if str == "Heart" || str == "H" || str == "h" {
		return Heart
	} else if str == "Diamond" || str == "D" || str == "d" {
		return Diamond
	} else if str == "Clover" || str == "C" || str == "c" {
		return Clover
	} else if str == "Joker" || str == "J" || str == "j" {
		return Joker
	} else if str == "Mighty" || str == "M" || str == "m" {
		return Mighty
	} else if str == "JokerCall" {
		return JokerCall
	} else if str == "NoTrump" || str == "N" || str == "n" {
		return NoTrump
	}

	return NoType
}

func getFriendType(str string) FriendType {
	if str == "no" {
		return NoFriend
	} else if str == "first" {
		return FirstFriend
	} else {
		return CardFriend
	}
}

func (c *Game) GetCard(str string) Card {
	typeid := str[:1]

	if typeid == "M" || typeid == "m" {
		if c.BiddingType == Spade {
			return Card{Type: Diamond, Number: 14}
		} else {
			return Card{Type: Spade, Number: 14}
		}
	}

	if typeid == "J" || typeid == "j" {
		return Card{Type: Joker, Number: 0}
	}

	var cardType Type

	if typeid == "S" || typeid == "s" {
		cardType = Spade
	} else if typeid == "H" || typeid == "h" {
		cardType = Heart
	} else if typeid == "D" || typeid == "d" {
		cardType = Diamond
	} else if typeid == "C" || typeid == "c" {
		cardType = Clover
	}

	nu := str[1:]
	number := 0

	if nu == "A" || nu == "a" {
		number = 14
	} else if nu == "K" || nu == "k" {
		number = 13
	} else if nu == "Q" || nu == "q" {
		number = 12
	} else if nu == "J" || nu == "j" {
		number = 11
	} else {
		number = global.Atoi(nu)
	}

	return Card{Type: cardType, Number: number}
}

func getNumber(value int) string {
	if value == 14 {
		return "A"
	} else if value == 13 {
		return "K"
	} else if value == 12 {
		return "Q"
	} else if value == 11 {
		return "J"
	} else {
		return strconv.Itoa(value)
	}
}

func getString(card Card) string {
	if card.Type == Spade {
		return "♠" + getNumber(card.Number)
	} else if card.Type == Heart {
		return "♥" + getNumber(card.Number)
	} else if card.Type == Diamond {
		return "◆" + getNumber(card.Number)
	} else if card.Type == Clover {
		return "♣" + getNumber(card.Number)
	} else if card.Type == Joker {
		return "JO"
	} else if card.Type == Mighty {
		return ""
	} else if card.Type == JokerCall {
		return ""
	} else if card.Type == NoTrump {
		return ""
	} else {
		return ""
	}
}

func (c *Game) Print() {
	fmt.Printf("\n\n")
	fmt.Println("==========================================")
	fmt.Printf("ROUND : %v round\n", c.Round)
	fmt.Println("")
	if c.Round == 0 {
		fmt.Printf("Bidding : %v\n", c.Bidding)
		fmt.Printf("BiddingType : %v\n", c.BiddingType)
		fmt.Printf("BiddingCount : %v\n", c.BiddingCount)
	}

	fmt.Println("------------------------------------------")
	for i, user := range c.Users {
		turn := ""
		winner := ""

		if c.Round == 0 {
			if i == c.BiddingCount {
				turn = "[TURN]"
			}
		} else {
			if i == c.Winner {
				winner = "[WINNER] "
			}

			if i == c.Master {
				if c.Round > 1 {
					winner = "[MASTER] " + winner
				} else {
					winner = "[MASTER] "
				}
			}

			if i == (c.Winner+c.DrawCount)%5 {
				turn = "[TURN]"
			}
		}

		fmt.Printf("[%v] %v (chip : %v) %v%v\n", i, user.Name, user.Chip, winner, turn)

		if user.Draw != nil {
			fmt.Printf("Draw : %v\n", getString(*user.Draw))
		}

		fmt.Printf("Hand : ")

		for i, v := range user.Cards {
			if i > 0 {
				fmt.Printf(", ")
			}

			fmt.Printf(getString(v))
		}

		fmt.Printf("\n")

		if len(user.PointCards) > 0 {
			fmt.Printf("Win  : ")

			for i, v := range user.PointCards {
				if i > 0 {
					fmt.Printf(", ")
				}

				fmt.Printf(getString(v))
			}

			fmt.Printf("\n")

		}
		fmt.Printf("\n")
	}

	fmt.Println("==========================================")
	fmt.Printf("\n\n")
}

func (c *Game) Command(str string) error {
	msg := strings.Split(str, " ")
	if len(msg) == 0 {
		return errors.New("empty command")
	}

	cmd := msg[0]
	user := 0
	if len(msg) > 1 {
		user = global.Atoi(msg[1])
	}

	if cmd == "init" {
		conn := models.NewConnection()
		defer conn.Close()

		manager := models.NewGameuserManager(conn)
		items := manager.FindByGame(c.Id)

		if len(*items) != 5 {
			return errors.New("not enough user")
		}

		a := []int{0, 1, 2, 3, 4}
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

		var users []User
		for i := 0; i < 5; i++ {
			item := (*items)[a[i]]
			extra := item.Extra.(map[string]interface{})["user"].(models.User)
			user := User{
				Id:         item.Id,
				Name:       extra.Name,
				Cards:      nil,
				Draw:       nil,
				PointCards: nil,
				Chip:       0,
			}

			users = append(users, user)
		}

		c.Init(users)
	} else if cmd == "bidding" {
		pass := false
		typeid := NoType
		count := 0

		if msg[2] == "pass" {
			pass = true
		} else {
			typeid = getType(msg[2])
			count = global.Atoi(msg[3])
		}

		fmt.Println(typeid, count)

		return c.DoBidding(user, typeid, count, pass)
	} else if cmd == "tablecards" {
		cards := make([]Card, 0)

		for i := 2; i <= 4; i++ {
			cards = append(cards, c.GetCard(msg[i]))
		}

		typeid := NoType
		count := 0

		if len(msg) >= 7 {
			typeid = getType(msg[5])
			count = global.Atoi(msg[6])
		} else {
			typeid = c.BiddingType
			count = c.Bidding
		}

		return c.DoSetTableCards(user, cards, typeid, count)
	} else if cmd == "friend" {
		friendType := NoFriend
		var card Card

		str := msg[2][:1]
		if str == "M" || str == "m" {
			friendType = CardFriend
			if c.BiddingType == Spade {
				card = Card{Type: Diamond, Number: 14}
			} else {
				card = Card{Type: Spade, Number: 14}
			}
		} else if str == "J" || str == "j" {
			friendType = CardFriend
			card = Card{Type: Joker, Number: 0}
		} else if str == "F" || str == "f" {
			friendType = FirstFriend
			card = Card{Type: NoType, Number: 0}
		} else if str == "N" || str == "n" {
			friendType = NoFriend
			card = Card{Type: NoType, Number: 0}
		} else {
			friendType = CardFriend
			card = c.GetCard(msg[2])
		}

		return c.DoSetFriend(user, friendType, card)
	} else if cmd == "draw" {
		var card Card
		jokerCall := false
		if len(msg[2]) == 1 || msg[2] == "10" {
			str := msg[2][:1]

			if str == "M" || str == "m" {
				if c.BiddingType == Spade {
					card = Card{Type: Diamond, Number: 14}
				} else {
					card = Card{Type: Spade, Number: 14}
				}
			} else if str == "J" || str == "j" {
				card = Card{Type: Joker, Number: 0}
			} else {
				card = c.Users[user].Cards[global.Atoi(msg[2])-1]
			}
		} else {
			card = c.GetCard(msg[2])
		}

		if len(msg) >= 4 {
			str := msg[3][:1]
			if msg[3] == "jokercall" || str == "Y" || str == "y" {
				jokerCall = true
			}
		}

		return c.DoDraw(user, card, jokerCall)
	}

	return nil
}

func (c *Game) Init(users []User) {
	c.Users = users
	c.Round = 0
	c.Bidding = 0
	c.BiddingType = NoTrump
	c.BiddingUser = 0
	c.Winner = 0

	c.ResetRound()
	c.DealCards()
}

func (c *Game) DealCards() {
	cards := make([]Card, 0)

	for i := 2; i <= 14; i++ {
		cards = append(cards, Card{Type: Spade, Number: i})
		cards = append(cards, Card{Type: Heart, Number: i})
		cards = append(cards, Card{Type: Diamond, Number: i})
		cards = append(cards, Card{Type: Clover, Number: i})
	}

	cards = append(cards, Card{Type: Joker, Number: 0})

	count := len(cards)

	for j := 0; j < 10; j++ {
		for i := 0; i < count; i++ {
			target := rand.Intn(count)

			if i == target {
				continue
			}

			typeid := cards[target].Type
			number := cards[target].Number

			cards[target].Type = cards[i].Type
			cards[target].Number = cards[i].Number

			cards[i].Type = typeid
			cards[i].Number = number
		}
	}

	pos := 0
	for i := 0; i < 5; i++ {
		c.Users[i].Chip = 50
		c.Users[i].Draw = nil
		c.Users[i].PointCards = make([]Card, 0)
		c.Users[i].Cards = make([]Card, 0)
		for j := 0; j < 10; j++ {
			c.Users[i].Cards = append(c.Users[i].Cards, cards[pos])
			pos++
		}
	}

	c.TableCards = make([]Card, 0)
	c.TableCards = append(c.TableCards, cards[pos])
	pos++
	c.TableCards = append(c.TableCards, cards[pos])
	pos++
	c.TableCards = append(c.TableCards, cards[pos])
}

// CheckDealMiss checks if a player has a deal miss (딜 미스)
// Returns true if the player has no point cards (except Mighty) or has Joker and 1 or less point cards
func (c *Game) CheckDealMiss(user int) bool {
	pointCardCount := 0
	hasMighty := false
	hasJoker := false

	for _, card := range c.Users[user].Cards {
		if c.IsMighty(card) {
			hasMighty = true
		} else if card.Type == Joker {
			hasJoker = true
		} else if card.Number >= 10 {
			pointCardCount++
		}
	}

	// If no point cards (except Mighty), it's a deal miss
	if pointCardCount == 0 && !hasMighty {
		return true
	}

	// If has Joker and 1 or less point cards, it's a deal miss
	if hasJoker && pointCardCount <= 1 {
		return true
	}

	return false
}

func (c *Game) DoBidding(user int, carttype Type, count int, pass bool) error {
	if user != c.BiddingUser {
		return errors.New("not turn")
	}

	if pass == true {
		if c.BiddingCount == 0 {
			return errors.New("must bidding")
		}
		c.Users[user].Pass = true
	} else {
		if count < 13 || count > 20 {
			return errors.New(fmt.Sprintf("bidding range error : %v", count))
		}

		if count <= c.Bidding {
			return errors.New("bidding value error")
		}

		if c.Users[user].Pass == true {
			return errors.New("already pass")
		}

		c.Bidding = count
		c.BiddingType = carttype
	}

	// 다른 사람이 있는지 검사
	flag := false
	next := -1
	remain := 0

	for i := user + 1; i < 5; i++ {
		if c.Users[i].Pass == false {
			if flag == false {
				next = i
			}

			flag = true
			remain++

			break
		}
	}

	if user != 0 {
		for i := 0; i < user; i++ {
			if c.Users[i].Pass == false {
				if flag == false {
					next = i
				}

				flag = true
				remain++

				break
			}
		}
	}

	if flag == true {
		c.BiddingUser = next
	}

	c.BiddingCount++

	if pass == true && remain == 1 {
		c.SetMaster(next)

		// send all setmaster
	}

	return nil
}

func (c *Game) ResetRound() {
	c.DrawCount = 0
	c.BiddingCount = 0
	c.BiddingType = NoTrump
	c.JokerCall = false

	for i, _ := range c.Users {
		c.Users[i].Draw = nil
		c.Users[i].Pass = false
	}
}

func (c *Game) SetWinner(user int) {
	c.Winner = user
}

func (c *Game) SetMaster(user int) {
	c.Master = user
	c.SetWinner(user)

	for i := 0; i < 3; i++ {
		c.Users[c.Master].Cards = append(c.Users[c.Master].Cards, c.TableCards[i])
	}

	c.TableCards = make([]Card, 0)
}

func (c *Game) DoSetTableCards(user int, cards []Card, trump Type, bidding int) error {
	if user != c.Master {
		return errors.New("not master")
	}

	if bidding > 20 || bidding < c.Bidding {
		return errors.New("bidding range error")
	}

	// Check trump change rules
	if trump != c.BiddingType {
		// From trump to NoTrump: +1 is allowed
		if c.BiddingType != NoTrump && trump == NoTrump {
			if bidding < c.Bidding+1 {
				return errors.New("bidding change error: trump to NoTrump requires +1")
			}
		} else if c.BiddingType == NoTrump && trump != NoTrump {
			// From NoTrump to trump: +1 is allowed
			if bidding < c.Bidding+1 {
				return errors.New("bidding change error: NoTrump to trump requires +1")
			}
		} else {
			// Other trump changes: +2 required
			if bidding < c.Bidding+2 {
				return errors.New("bidding change error: trump change requires +2")
			}
		}
	}

	if len(cards) != 3 {
		return errors.New("not 3 cards")
	}

	newCards := make([]Card, 0)

	for _, v := range c.Users[c.Master].Cards {
		flag := false
		for _, t := range cards {
			if v.Type == t.Type && v.Number == t.Number {
				c.TableCards = append(c.TableCards, v)
				flag = true
				break
			}
		}

		if flag == false {
			newCards = append(newCards, v)
		}
	}

	if len(newCards) != 10 {
		return errors.New("cards error")
	}

	c.Trump = trump
	c.Bidding = bidding

	c.Users[c.Master].Cards = newCards

	return nil
}

func (c *Game) DoSetFriend(user int, friendType FriendType, card Card) error {
	if user != c.Master {
		return errors.New("not master")
	}

	c.Round = 1

	c.Friend.Type = friendType
	c.Friend.Card = card

	if friendType == NoFriend {
		c.Friend.Id = c.Master
	} else if friendType == CardFriend {
		for i, v := range c.Users {
			if i == c.Master {
				continue
			}

			for _, item := range v.Cards {
				if card.Type == Mighty {
					if c.IsMighty(item) {
						c.Friend.Id = i
						return nil
					}
				} else if card.Type == JokerCall {
					if c.IsJokerCall(item) {
						c.Friend.Id = i
						return nil
					}
				} else if item.Type == card.Type && item.Number == card.Number {
					c.Friend.Id = i
					return nil
				}
			}
		}

		c.Friend.Id = c.Master
		c.Friend.Type = NoFriend
	}

	return nil
}

func (c *Game) IsMighty(card Card) bool {
	if card.Number != 14 {
		return false
	}

	if c.Trump == Spade {
		if card.Type == Diamond {
			return true
		}
	} else {
		if card.Type == Spade {
			return true
		}
	}

	return false
}

func (c *Game) IsJokerCall(card Card) bool {
	if card.Number != 3 {
		return false
	}

	if c.Trump == Clover {
		if card.Type == Spade {
			return true
		}
	} else {
		if card.Type == Clover {
			return true
		}
	}

	return false
}

func (c *Game) IsHaveJoker(user int) bool {
	for _, v := range c.Users[user].Cards {
		if v.Type == Joker {
			return true
		}
	}

	return false
}

func (c *Game) DoDraw(user int, card Card, jokerCall bool) error {
	if user != (c.Winner+c.DrawCount)%5 {
		return errors.New("not turn")
	}

	if c.Users[user].Draw != nil {
		return errors.New("already turn")
	}

	if c.DrawCount == 0 && c.IsJokerCall(card) && jokerCall == true {
		c.JokerCall = true
	}

	if c.JokerCall == true && c.DrawCount > 0 && c.IsHaveJoker(user) {
		if card.Type != Joker {
			return errors.New("not joker")
		}
	}

	if c.DrawCount > 0 {
		if !c.IsMighty(card) && card.Type != Joker {
			firstType := c.Users[c.Winner].Draw.Type
			if firstType != card.Type {
				for _, v := range c.Users[user].Cards {
					if v.Type == firstType {
						return errors.New("wrong card type")
					}
				}
			}
		}
	}

	for i, v := range c.Users[user].Cards {
		if v.Type == card.Type && v.Number == card.Number {
			c.Users[user].Draw = &v

			c.Users[user].Cards = remove(c.Users[user].Cards, i)

			c.DrawCount++

			if c.DrawCount == 5 {
				c.RoundEnd()
			}
			return nil
		}
	}

	return errors.New("not found")
}

func (c *Game) RoundEnd() {
	winner := -1

	// First, check for Mighty (always wins)
	for i, v := range c.Users {
		if c.IsMighty(*v.Draw) {
			winner = i
			break
		}
	}

	// Then check for Joker (but not in first/last round, and not when JokerCall is active)
	if winner == -1 {
		for i, v := range c.Users {
			// Joker wins except:
			// - First round (Round 1)
			// - Last round (Round 10)
			// - When JokerCall is active
			if v.Draw.Type == Joker && c.Round != 1 && c.Round != 10 && !c.JokerCall {
				winner = i
				break
			}
		}
	}

	if winner == -1 {
		trump := c.Trump
		if c.Trump == NoTrump {
			trump = c.Users[c.Winner].Draw.Type
		}

		max := 0
		pos := -1

		for i, v := range c.Users {
			if v.Draw.Type == trump {
				if v.Draw.Number > max {
					max = v.Draw.Number
					pos = i
				}
			}
		}

		if pos > -1 {
			winner = pos
		} else {
			max = 0
			pos = -1

			trump = c.Users[c.Winner].Draw.Type

			for i, v := range c.Users {
				if v.Draw.Type == trump {
					if v.Draw.Number > max {
						max = v.Draw.Number
						pos = i
					}
				}
			}

			winner = pos
		}
	}

	if c.Round == 1 {
		if c.Friend.Type == FirstFriend {
			c.Friend.Id = winner

			if winner == c.Master {
				c.Friend.Type = NoFriend
			}
		}
	}

	for _, v := range c.Users {
		if v.Draw.Number >= 10 {
			c.Users[winner].PointCards = append(c.Users[winner].PointCards, *v.Draw)
		}
	}

	c.SetWinner(winner)
	c.ResetRound()

	if c.Round < 10 {
		c.Round++
	} else {
		c.GameEnd()
	}
}

func (c *Game) GameEnd() {
	// Calculate opposition's point cards
	oppositionPoints := 0
	for i, v := range c.Users {
		if i != c.Master && i != c.Friend.Id {
			oppositionPoints += len(v.PointCards)
		}
	}

	// Ruling party's point cards (20 - opposition's cards)
	rulingPoints := 20 - oppositionPoints

	ret := 0

	// Check for Run (ruling party takes all 20 cards)
	isRun := (rulingPoints == 20)

	// Check for Back Run (opposition takes 11+ cards OR more than the bid)
	isBackRun := (oppositionPoints >= 11 || oppositionPoints >= c.Bidding)

	if rulingPoints >= c.Bidding {
		// Ruling party wins
		if isRun {
			// Run: all 20 cards
			if c.Bidding == 20 {
				ret = 40 // Perfect bid gets double
			} else {
				ret = 20 // Run base score
			}
			// NoTrump doubles even for Run
			if c.Trump == NoTrump {
				ret = ret * 2
			}
		} else {
			// Normal win: (Bidding - 13) * 2 + (RulingPoints - Bidding)
			// This rewards higher bids and excess points
			ret = (c.Bidding-13)*2 + (rulingPoints - c.Bidding)

			// NoTrump doubles the score
			if c.Trump == NoTrump {
				ret = ret * 2
			}
		}

		// NoFriend bonus is handled in distribution (Master gets 4x instead of 2x)
	} else {
		// Opposition wins
		shortage := c.Bidding - rulingPoints
		ret = -1 * shortage

		// Back Run: doubles the loss
		if isBackRun {
			ret = ret * 2
		}
	}

	// Distribute scores
	for i := range c.Users {
		if i == c.Master {
			if c.Friend.Type == NoFriend {
				// No friend: Master gets 4x (playing alone bonus)
				c.Users[i].Chip += ret * 4
			} else {
				// With friend: Master gets 2x
				c.Users[i].Chip += ret * 2
			}
		} else if i == c.Friend.Id {
			// Friend gets 1x
			c.Users[i].Chip += ret
		} else {
			// Opposition loses ret points each
			c.Users[i].Chip -= ret
		}
	}
}
