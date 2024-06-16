package pg

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"nexa/services/user/internal/domain/entity"
	model2 "nexa/services/user/internal/infra/repository/model"
	"nexa/shared/database"
	"nexa/shared/env"
	"nexa/shared/types"
	"nexa/shared/util"
	"nexa/shared/wrapper"
	"testing"
	"time"
)

var Db *bun.DB

var UserIds = []string{}

var Users = []entity.User{}

var Profiles = []entity.Profile{}

const USER_SIZE = 5
const PROFILE_SIZE = 3

func TestMain(m *testing.M) {
	for i := 0; i < USER_SIZE; i += 1 {
		id := types.Id(uuid.New())
		UserIds = append(UserIds, id.Underlying().String())
		Users = append(Users, newRandomUserWithId(id))

		if i >= PROFILE_SIZE {
			continue
		}
		Profiles = append(Profiles, newRandomProfile(id))
	}

	err := env.LoadEnvs("../../../../dev.env")
	if err != nil {
		panic(err)
	}

	config, err := database.LoadConfig()
	if err != nil {
		panic(err)
	}

	Db, err = database.OpenPostgres(config, false)
	if err != nil {
		panic(err)
	}
	defer Db.Close()

	model2.RegisterBunModels(Db)
	err = model2.CreateTables(Db)
	if err != nil {
		panic(err)
	}

	err = seed()
	if err != nil {
		panic(err)
	}
	defer func() {
		res, err := Db.NewDelete().
			Model(util.Nil[model2.User]()).
			Where("id IN (?)", bun.In(UserIds)).
			Exec(context.Background())

		util.DoNothing(res, err)
	}()

	m.Run()
}

func newRandomUser() entity.User {
	user := entity.User{
		Id:         types.Id(wrapper.SomeF(uuid.NewUUID).Data),
		Username:   util.RandomString(15),
		Email:      types.Email(util.RandomString(12)),
		Password:   wrapper.Some(types.PasswordFromString(util.RandomString(12))).Data,
		IsVerified: util.RandomBool(),
		IsDeleted:  util.RandomBool(),
	}
	return user
}

func newRandomUserWithId(id types.Id) entity.User {
	user := entity.User{
		Id:         id,
		Username:   util.RandomString(15),
		Email:      wrapper.DropError(types.EmailFromString(util.RandomString(12) + "@gmail.com")),
		Password:   wrapper.DropError(types.PasswordFromString(util.RandomString(12))),
		IsVerified: util.RandomBool(),
		IsDeleted:  util.RandomBool(),
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
	profileModels := util.CastSlice(Profiles, func(from *entity.Profile) model2.Profile {
		return model2.FromProfileDomain(from, func(profile *entity.Profile, db *model2.Profile) {
			db.UpdatedAt = time.Now()
		})
	})

	userModels := util.CastSlice(Users, func(from *entity.User) model2.User {
		return model2.FromUserDomain(from, func(domain *entity.User, user *model2.User) {
			user.CreatedAt = time.Now()
			user.UpdatedAt = time.Now()
		})
	})

	err := database.Seed(Db, userModels...)
	if err != nil {
		return err
	}

	return database.Seed(Db, profileModels...)
}
