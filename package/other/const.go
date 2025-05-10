package other

const DefaultPort = "7540"
const DefaultNameDB = "scheduler.db"

var Install bool

const Schema = `CREATE TABLE IF NOT EXISTS scheduler(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT(128))`

const DBindex = `CREATE  INDEX IF NOT EXISTS dateIndex ON scheduler(date)`

const TimeFormat = "20060102"
