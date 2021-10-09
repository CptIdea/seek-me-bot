package service

import (
	"reflect"
	"seek-me-bot/service/pkg"
	"testing"
)

func TestNewGameController(t *testing.T) {
	tests := []struct {
		name string
		want GameController
	}{
		{"basic", &basicGameController{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGameController(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGameController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_basicGameController_AddPetition(t *testing.T) {
	type fields struct {
		game pkg.Game
	}
	type args struct {
		petition pkg.Petition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "first petition", args: args{pkg.Petition{"user1", []string{"one", "two", "3"}, "user2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &basicGameController{
				game: tt.fields.game,
			}
			if err := b.AddPetition(tt.args.petition); (err != nil) != tt.wantErr {
				t.Errorf("AddPetition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_basicGameController_GetPetition(t *testing.T) {
	type fields struct {
		game pkg.Game
	}
	tests := []struct {
		name    string
		fields  fields
		want    pkg.Petition
		wantErr bool
	}{
		{name: "empty game", wantErr: true},
		{name: "one petition", fields: fields{game: pkg.Game{Petitions: []pkg.Petition{pkg.Petition{Author: "test"}}}}, want: pkg.Petition{Author: "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &basicGameController{
				game: tt.fields.game,
			}
			got, err := b.GetPetition()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPetition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPetition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_basicGameController_ResetGame(t *testing.T) {
	type fields struct {
		game pkg.Game
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "empty game", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &basicGameController{
				game: tt.fields.game,
			}
			if err := b.ResetGame(); (err != nil) != tt.wantErr {
				t.Errorf("ResetGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
