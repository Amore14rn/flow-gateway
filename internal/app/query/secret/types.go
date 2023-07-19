package secret

import "time"

type Secret struct {
	UUID    string    `db:"uuid"`
	Secret  string    `db:"secret"`
	Created time.Time `db:"created_at"`
}
