/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "SpaceDock"
    "SpaceDock/utils"
    "errors"
    "time"
)

type Game struct {
    Model

    Name string `gorm:"size:1024;unique_index;not null" json:"name"`
    Active bool `json:"active"`
    Fileformats string `gorm:"size:1024" json:"fileformats" spacedock:"json"`
    Altname string `gorm:"size:1024" json:"altname"`
    Rating float32 `json:"rating"`
    Releasedate time.Time `json:"releasedate"`
    Short string `gorm:"size:1024" json:"short"`
    publisherID uint `json:"publisher"`
    Description string `gorm:"size:100000" json:"description"`
    ShortDescription string `gorm:"size:1000" json:"short_description"`
    // Mods []Mod
    // Modlists []ModList
    // Versions []GameVersion

}

func NewGame(name string, publisher Publisher, short string) *Game {
    game := &Game {
        Name: name,
        Active: false,
        Altname: "",
        Rating: 0,
        Releasedate: time.Now(),
        Short: short,
        Description: "",
        ShortDescription: "",
        publisherID: publisher.ID,
    }
    game.Meta = "{}"
    return game
}

func (game *Game) GetById(id interface{}) error {
    SpaceDock.Database.First(&game, id)
    if game.Name != "" {
        return errors.New("Invalid ability ID")
    }
    return nil
}

func (game Game) GetPublisher() *Publisher {
    pub := &Publisher{}
    err := pub.GetById(game.publisherID)
    if err != nil {
        return nil
    }
    return pub
}

func (game Game) Format() map[string]interface{} {
    return map[string]interface{} {
        "id": game.ID,
        "name": game.Name,
        "active": game.Active,
        "rating": game.Rating,
        "releasedate": game.Releasedate,
        "short": game.Short,
        "publisher": game.publisherID,
        "description": game.Description,
        "short_description": game.ShortDescription,
        "created": game.CreatedAt,
        "updated": game.UpdatedAt,
        "meta": utils.LoadJSON(game.Meta),
    }
}