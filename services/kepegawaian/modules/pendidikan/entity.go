package pendidikan

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type pendidikan struct {
	ID                  string      `json:"id"`
	Nama                pgtype.Text `json:"nama"`
	TingkatPendidikan   pgtype.Text `json:"tingkat_pendidikan,omitempty"`
	TingkatPendidikanID pgtype.Int2 `json:"tingkat_pendidikan_id"`
}

// generateID generates a unique ID for a pendidikan based on legacy system
func generateID() (string, error) {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	source := append([]byte(time.Now().Format(time.RFC3339Nano)), randomBytes...)
	sum := md5.Sum(source)
	return hex.EncodeToString(sum[:]), nil
}
