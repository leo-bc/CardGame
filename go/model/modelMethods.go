package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/robertkrimen/otto"
)

// DrawCards :
func DrawCards(state *State, seed int) []int {
	rand.Seed(time.Now().UnixNano() + int64(seed))

	indices := []int{}
	for i := 0; i < 60; i++ {
		index := rand.Intn(len(state.Cards))
		indices = append(indices, index)
	}

	return indices
}

// DrawCardsByRank :
func DrawCardsByRank(state *State, seed int) []int {
	var drawnLegends int
	var drawnRares int
	var drawnMinions int
	var result []int

	legendsMax := 5
	raresMax := 10
	minionsMax := 20

	rand.Seed(time.Now().UnixNano() + int64(seed))

	for drawnLegends < legendsMax || drawnRares < raresMax || drawnMinions < minionsMax {
		index := rand.Intn(len(state.Cards))
		drawnCard := &state.Cards[index]
		if drawnLegends < legendsMax && drawnCard.Rank.Ranking == "Legend" {
			drawnLegends++
			result = append(result, index)
		} else if drawnRares < raresMax && drawnCard.Rank.Ranking == "Rare" {
			drawnRares++
			result = append(result, index)
		} else if drawnMinions < minionsMax && drawnCard.Rank.Ranking == "Minion" {
			drawnMinions++
			result = append(result, index)
		}
	}

	return result
}

// StartGame :
func StartGame(state *State, game *Game) {
	for i := 0; i < len(game.Info.Players); i++ {
		player := &game.Info.Players[i]
		cardIDs := DrawCards(state, i)

		for _, cardID := range cardIDs {
			player.CardIDs = append(player.CardIDs, cardID)
		}
	}

	side0 := BattleSide{PlayerID: 0, Info: BattleSideInfo{Cards: make(map[string]*CardsContainer)}}
	side0.Info.Cards["Bench"] = &CardsContainer{}
	side0.Info.Cards["Hand"] = &CardsContainer{}
	side0.Info.Cards["TakePile"] = &CardsContainer{}
	side0.Info.Cards["ThrowPile"] = &CardsContainer{}
	side0.Info.Cards["Prizes"] = &CardsContainer{}

	side1 := BattleSide{PlayerID: 2, Info: BattleSideInfo{Cards: make(map[string]*CardsContainer)}}
	side1.Info.Cards["Bench"] = &CardsContainer{}
	side1.Info.Cards["Hand"] = &CardsContainer{}
	side1.Info.Cards["TakePile"] = &CardsContainer{}
	side1.Info.Cards["ThrowPile"] = &CardsContainer{}
	side1.Info.Cards["Prizes"] = &CardsContainer{}

	battle := Battle{Sides: []BattleSide{side0, side1}}
	rand.Seed(time.Now().UnixNano())
	for i := range battle.Sides {
		side := &battle.Sides[i]

		cardIDsCopy := []int{}
		for _, id := range game.Info.Players[side.PlayerID].CardIDs {
			cardIDsCopy = append(cardIDsCopy, id)
		}
		rand.Shuffle(len(cardIDsCopy), func(i, j int) { cardIDsCopy[i], cardIDsCopy[j] = cardIDsCopy[j], cardIDsCopy[i] })
		for i, id := range cardIDsCopy {
			if i < 1 {
				side.Info.Cards["Prizes"].Append(CardSlot{CardID: id})
			} else {
				side.Info.Cards["TakePile"].Append(CardSlot{CardID: id})
			}
		}
	}

	game.Info.Battles = append(game.Info.Battles, battle)
}

func readJSON(state *State, create bool) {
	if create {
		jsonFile, err := os.Open("C:\\Users\\Leo\\Documents\\CardGameWebsite\\cards\\manipulatedPokemon.json")

		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var cards []Card
		json.Unmarshal(byteValue, &cards)

		for _, card := range cards {
			state.Cards = append(state.Cards, card)
		}
	} else {
		jsonFile, err := os.Open("C:\\Users\\Leo\\Documents\\CardGameWebsite\\cards\\exampleState.json")

		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var newState State
		json.Unmarshal(byteValue, &newState)
		*state = newState
	}
}

// StartBattle :
func StartBattle(battle *Battle) {
	for i := range battle.Sides {
		side := &battle.Sides[i]
		DrawCard(&side.Info, 7)
	}
}

// BattleStartTurn :
func BattleStartTurn(battle *Battle) {

}

// PlayCard :
func PlayCard(side *BattleSide, handID int) {
	card := side.Info.Cards["Hand"].Slots[handID]
	newSlots := []CardSlot{}
	for i := range side.Info.Cards["Hand"].Slots {
		if i != handID {
			newSlots = append(newSlots, side.Info.Cards["Hand"].Slots[i])
		}
	}
	side.Info.Cards["Hand"].Slots = newSlots
	side.Info.Cards["Bench"].Append(card)
}

func doDamage(target *CardSlot, damage int) {
	target.DamageTaken += damage
}

func setStatus(target *CardSlot, status string) {
	target.Status = status
}

// AttackCard :
func AttackCard(state *State, battle *Battle, playerSide *BattleSide, opponentSide *BattleSide, attackInfo AttackInfo, vm *otto.Otto) {
	cardID := playerSide.Info.Cards["Bench"].Slots[attackInfo.Source].CardID
	card := state.Cards[cardID]

	opponentCard := &opponentSide.Info.Cards["Bench"].Slots[attackInfo.Target]

	vm.Set("doDamage", func(call otto.FunctionCall) otto.Value {
		damage, _ := call.Argument(0).ToInteger()
		self, _ := call.Argument(0).ToBoolean()
		target := opponentCard
		if self {
			target = &playerSide.Info.Cards["Bench"].Slots[attackInfo.Source]
		}
		doDamage(target, int(damage))
		return otto.Value{}
	})

	vm.Set("setStatus", func(call otto.FunctionCall) otto.Value {
		status, _ := call.Argument(0).ToString()
		setStatus(opponentCard, status)
		return otto.Value{}
	})

	vm.Set("card", card)
	vm.Set("attack", card.Attacks[attackInfo.Attack])

	Standard := `
		doDamage(attack.Damage, false);
	`

	FoulGas := `
		doDamage(attack.Damage, false);
		var result = flipCoin()
		if (result) {
			setStatus(POISONED);
		} else {
			setStatus(CONFUSED);
		}
	`

	ThunderPunch := `
		var result = flipCoin()
		if (result) {
			doDamage(attack.Damage + 10, false);
		} else {
			doDamage(attack.Damage);
			doDamage(10, true);
		}
	`

	attacks := []string{Standard, FoulGas}
	fmt.Printf("%s\n", attacks)

	_, err := vm.Run(ThunderPunch)
	if err != nil {
		fmt.Printf("%s", err)
	}
	if opponentCard.DamageTaken >= state.Cards[opponentCard.CardID].Rank.HP {
		opponentSide.Info.Cards["Bench"].Replace([]CardSlot{})
		opponentSide.Info.Cards["ThrowPile"].Append(*opponentCard)

		newPrizes := []CardSlot{}
		for i, slot := range opponentSide.Info.Cards["Prizes"].Slots {
			if i > 0 {
				newPrizes = append(newPrizes, slot)
			}
		}
		opponentSide.Info.Cards["Prizes"].Replace(newPrizes)
	}

	if len(opponentSide.Info.Cards["Prizes"].Slots) == 0 {
		endBattle(battle, playerSide)
	}
}

func endBattle(battle *Battle, winner *BattleSide) {
	battle.Info.IsFinished = true
	battle.Info.Winner = winner.PlayerID
}

// BattleEndTurn :
func BattleEndTurn(battle *Battle) {
	side := &battle.Sides[battle.Info.Turn]
	DrawCard(&side.Info, 1)
	battle.Info.Turn = (battle.Info.Turn + 1) % len(battle.Sides)
}

// DrawCard :
func DrawCard(side *BattleSideInfo, amount int) {
	for i := 0; i < amount; i++ {
		if len(side.Cards["TakePile"].Slots) > 0 {
			drawCard := side.Cards["TakePile"].Slots[0]
			side.Cards["TakePile"].Replace(side.Cards["TakePile"].Slots[1:])
			side.Cards["Hand"].Append(drawCard)
		}
	}
}

// GetExampleState :
func GetExampleState(create bool) State {
	var state = State{Games: []Game{}, Players: []Player{}, Cards: []Card{}}

	readJSON(&state, create)

	if create {
		player1 := Player{Name: "dj leo"}
		state.Players = append(state.Players, player1)

		player2 := Player{Name: "funky kong"}

		state.Players = append(state.Players, player2)

		player3 := Player{Name: "rick"}

		state.Players = append(state.Players, player3)

		game1 := Game{Info: GameInfo{Players: []PlayerGameLink{}, Battles: []Battle{}}}
		game1.Info.Title = "Cool game!"

		var deck []int

		link1 := PlayerGameLink{PlayerID: 0, CardIDs: deck}
		link2 := PlayerGameLink{PlayerID: 1, CardIDs: deck, IsReady: true}
		link3 := PlayerGameLink{PlayerID: 2, CardIDs: deck, IsReady: true}

		game1.Info.Players = append(game1.Info.Players, link1, link2, link3)
		state.Games = append(state.Games, game1)
	}
	return state
}
