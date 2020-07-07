package utils

import (
	"errors"
	"time"
)

type Storage struct {
	Server Server
}

type FObject struct {
	Name      string
	Location  string
	LocalPath string

	IsDir    bool
	IsBinary bool

	Lastwritten time.Time
	Lastpulled  time.Time

	Version int16
}

func (s *Storage) SaveObject(r *UserRequestPackage, data *[]byte) error {
	u, ok := s.Server.Users[r.ClientUser.Username]
	if !ok {
		return errors.New("User Not Found")
	}

	if u.Key != r.ClientUser.Key {
		return errors.New("key doesn'y match")
	}

	return nil
}

func (s *Storage) VerifyUser(u *string, p *string) error {
	return errors.New("User not Found")
}
