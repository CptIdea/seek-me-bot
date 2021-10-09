package service

import (
	"math/rand"
	"seek-me-bot/service/pkg"
)

type GameController interface {
	AddPetition(petition pkg.Petition) error
	GetPetition() (pkg.Petition, error)
	ResetGame() error
}

type basicGameController struct {
	game pkg.Game
}

func NewGameController() GameController {
	return &basicGameController{}
}

func (b *basicGameController) AddPetition(petition pkg.Petition) error {
	b.game.Petitions = append(b.game.Petitions, petition)
	return nil
}

func (b *basicGameController) GetPetition() (pkg.Petition, error) {
	if len(b.game.Petitions) == 0 {
		return pkg.Petition{}, EmptyGameError
	}

	i := rand.Intn(len(b.game.Petitions))
	ans := b.game.Petitions[i]
	b.game.Petitions = append(b.game.Petitions[:i], b.game.Petitions[i+1:]...)

	return ans, nil
}

func (b *basicGameController) ResetGame() error {
	b.game.Petitions = []pkg.Petition{}
	return nil
}
