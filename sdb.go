/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package main

import (
    "SpaceDock"
     "SpaceDock/objects"
    _ "SpaceDock/routes"
    "strconv"
    "regexp"
)

/*
 The entrypoint for the spacedock application.
 Instead of running significant code here, we pass this task to the spacedock package
*/
func main() {
    if (SpaceDock.Settings.CreateDefaultDatabase) {
        CreateDefaultData()
    }
    SpaceDock.Run()
}

func CreateDefaultData() {

    // Setup users
    NewDummyUser("Administrator", "admin", "admin@example.com", true)
    NewDummyUser("SpaceDockUser", "user", "user@example.com", false)

    // Game 1
    ksp := NewDummyGame("Kerbal Space Program", "kerbal-space-program", "Squad MX")
    NewDummyVersion(*ksp, "1.2.1", false)
    NewDummyVersion(*ksp, "1.2.2", false)
    NewDummyVersion(*ksp, "1.2.9", true)

    // Game 2
    fac := NewDummyGame("Factorio", "factorio", "Wube Software")
    NewDummyVersion(*fac, "0.12", false)
}

func AddAbilityRe(role *objects.Role, expression string) {
    var abilities []objects.Ability
    SpaceDock.Database.Find(&abilities)
    for _,element := range abilities {
        if ok,_ := regexp.MatchString(expression, element.Name); ok {
            role.AddAbility(element.Name)
        }
    }
}

func NewDummyUser(name string, password string, email string, admin bool) *objects.User {
    user := objects.NewUser(name, email, password)
    SpaceDock.Database.Save(user)

    // Setup roles
    role := user.AddRole(user.Username)
    role.AddAbility("user-edit")
    role.AddAbility("mods-add")
    role.AddAbility("lists-add")
    role.AddAbility("logged-in")
    role.AddParam("user-edit", "userid", strconv.Itoa(int(user.ID)))
    role.AddParam("mods-add", "gameshort", ".*")
    role.AddParam("lists-add", "gameshort", ".*")
    SpaceDock.Database.Save(&role)

    // Admin roles
    if admin {
        admin_role := user.AddRole("admin")
        AddAbilityRe(admin_role, ".*")
        admin_role.AddAbility("mods-invite")
        admin_role.AddAbility("view-users-full")

        // Params
        admin_role.AddParam("admin-impersonate", "userid", ".*")
        admin_role.AddParam("game-edit", "gameshort", ".*")
        admin_role.AddParam("game-add", "pubid", ".*")
        admin_role.AddParam("game-remove", "short", ".*")
        admin_role.AddParam("mods-feature", "gameshort", ".*")
        //admin_role.AddParam("mods-edit", "gameshort", ".*")
        //admin_role.AddParam("mods-add", "gameshort", ".*")
        //admin_role.AddParam("mods-remove", "gameshort", ".*")
        admin_role.AddParam("lists-add", "gameshort", ".*")
        admin_role.AddParam("lists-edit", "gameshort", ".*")
        admin_role.AddParam("lists-remove", "gameshort", ".*")
        admin_role.AddParam("publisher-edit", "publid", ".*")
        admin_role.AddParam("token-edit", "tokenid", ".*")
        admin_role.AddParam("token-remove", "tokenid", ".*")
        admin_role.AddParam("user-edit", "userid", ".*")

        SpaceDock.Database.Save(&admin_role)
    }

    // Confirmation
    user.Confirmation = ""
    user.Public = true
    SpaceDock.Database.Save(user)
    return user
}

func NewDummyGame(name string, short string, publisher string) *objects.Game {
    pub := objects.NewPublisher(publisher)
    SpaceDock.Database.Save(pub)

    // Create the game
    game := objects.NewGame(name, *pub, short)
    game.Active = true
    SpaceDock.Database.Save(game)
    return game
}

func NewDummyVersion(game objects.Game, name string, beta bool) *objects.GameVersion {
    version := objects.NewGameVersion(name, game, beta)
    SpaceDock.Database.Save(version)
    return version
}
