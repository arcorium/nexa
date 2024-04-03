package model

import (
	"github.com/uptrace/bun"
	domain "nexa/services/user/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util/repo"
	"nexa/shared/variadic"
	"time"
)

type ProfileMapOption = repo.DataAccessModelMapOption[*domain.Profile, *Profile]

func FromProfileDomain(profile *domain.Profile, opts ...ProfileMapOption) Profile {
	pfl := Profile{
		UserId:    profile.Id.Underlying().String(),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		PhotoURL:  profile.PhotoURL.FileName(),
		Bio:       profile.Bio,
	}

	variadic.New(opts...).DoAtFirst(func(p *ProfileMapOption) {
		(*p)(profile, &pfl)
	})

	return pfl
}

type Profile struct {
	bun.BaseModel `bun:"table:profiles"`

	UserId    string `bun:",nullzero,pk"` // Profile is unique per user
	FirstName string `bun:",notnull"`
	LastName  string `bun:",nullzero"`
	PhotoURL  string `bun:",nullzero"`
	Bio       string `bun:",nullzero"`

	UpdatedAt time.Time `bun:",nullzero,notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}

func (p *Profile) ToDomain() domain.Profile {
	return domain.Profile{
		Id:        types.IdFromString(p.UserId),
		FirstName: p.FirstName,
		LastName:  p.LastName,
		PhotoURL:  types.FilePathFromString(p.PhotoURL),
		Bio:       p.Bio,
	}
}
