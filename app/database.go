/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package app

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mssql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "log"
    "sync"
)

/*
 The database handler that provides the API to interact with the the DB
 */
var Database *gorm.DB

/*
 Counter for recursive reference fetching
 */
var DBRecursion map[uint64]int = map[uint64]int{}
var DBRecursionMax int = 2
var DBRecursionLock sync.RWMutex = sync.RWMutex{}

/*
 Establishes the connection to the database
 */
func LoadDatabase() {
    var db, err = gorm.Open(Settings.Dialect, Settings.ConnectionData)
    if err != nil {
        log.Fatalf("* Failed to connect to the database: %s", err)
    }
    Database = db
    log.Print("* Database connection successfull")
    Database.LogMode(Settings.Debug)
}

/*
 Creates a table only if it doesn't exist
 */
func CreateTable(models interface{}) {
    if !Database.HasTable(models) {
        Database.CreateTable(models)
    }
}

func NoAssociations(callback func()) {
    Database.InstantSet("gorm:save_associations", false)
    callback()
    Database.InstantSet("gorm:save_associations", true)
}
