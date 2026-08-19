package tests

import (
	"os"
	"path/filepath"
)

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODY1NDc4OTIsImhhc2giOiJhNjY1YTQ1OTIwNDIyZjlkNDE3ZTQ4NjdlZmRjNGZiOGEwNGExZjNmZmYxZmEwN2U5OThlODZmN2Y3YTI3YWUzIn0.heLV1WBpYH9NEEcgISZ-uTSURUUU4x0DEvxiGKJDJ8E"

// init resolves the database file path to an absolute path upon package initialization
// and synchronizes TODO_DBFILE to ensure tests and the server target the exact same database file.
func init() {
	if abs, err := filepath.Abs(DBFile); err == nil {
		DBFile = abs
		os.Setenv("TODO_DBFILE", abs)
	}
}
