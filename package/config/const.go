package config

const (
	defaultPort = "7540"
	EnvPort     = "TODO_PORT"
	EnvDBFile   = "TODO_DBFILE"
	EnvDbPath   = "../"
	DbName      = "scheduler.db"
)
const Schema = `CREATE TABLE IF NOT EXISTS scheduler(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT,
	title TEXT,
	comment TEXT,
	repeat TEXT(128))`

const DBindex = `CREATE  INDEX IF NOT EXISTS dateIndex ON scheduler(date)`

const TimeFormat = "20060102"

var Install bool
