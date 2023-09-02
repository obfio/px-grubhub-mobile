package px

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// UUIDSection generates the keys PX326, PX327, and PX328.
func (p *Payload) UUIDSection() error {
	if p.Cache.TimeStamp == p.Cache.PrevUUIDTime {
		p.Cache.TimeStamp++
		p.Cache.PrevUUIDTime++
	}
	// ! generate UUIDv1
	p.D.Px326 = fmt.Sprint(uuid.Must(uuid.NewV1(p.Cache.TimeStamp)))
	p.D.Px327 = strings.ToUpper(p.D.Px326[0:8])
	p.D.Px328 = strings.ToUpper(px328(p.D.Px320 + p.D.Px326 + p.D.Px327))
	return nil
}

func px328(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
