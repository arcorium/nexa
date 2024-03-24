package pg

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"nexa/services/user/internal/infra/model"
	"nexa/services/user/shared/domain/entity"
	"nexa/shared/database"
	"nexa/shared/types"
	"nexa/shared/util"
	"nexa/shared/wrapper"
	"testing"
	"time"
)

var Db *bun.DB

var Users = []entity.User{
	newRandomUser(),
	newRandomUser(),
}

var Profiles = []entity.Profile{
	newRandomProfile(Users[0].Id),
}

func TestMain(m *testing.M) {
	config, err := database.LoadConfig("../../../dev")
	if err != nil {
		panic(err)
	}

	Db, err = database.OpenPostgres(config, false, (*model.User)(nil), (*model.Profile)(nil))
	defer Db.Close()

	err = seed()
	if err != nil {
		panic(err)
	}

	m.Run()
}

func newRandomUser() entity.User {
	user := entity.User{
		Id:         types.Id(wrapper.SomeF(uuid.NewUUID).Data),
		Username:   util.RandomString(15),
		Email:      wrapper.Some(types.EmailFromString(util.RandomString(12))).Data,
		Password:   types.PasswordFromString(util.RandomString(12)),
		IsVerified: false,
		IsDeleted:  false,
	}
	return user
}

func newRandomProfile(userId types.Id) entity.Profile {
	return entity.Profile{
		Id:        userId,
		FirstName: util.RandomString(6),
		LastName:  util.RandomString(8),
		Bio:       util.RandomString(20),
		PhotoURL:  types.FilePathFromString(util.RandomString(12)),
	}
}

func seed() error {
	models := util.CastSlice(Users, func(from *entity.User) model.User {
		return model.FromUserDomain(from, func(domain *entity.User, user *model.User) {
			user.CreatedAt = time.Now()
		})
	})
	return database.Seed(Db, models)
}
