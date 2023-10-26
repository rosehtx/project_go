package model

type UserHero struct {
	UserId uint64
	HeroId uint64
	HeroName string
}

type AllUserHero []*UserHero

func (uh *AllUserHero) GetHeroByUserId(userId uint64) {
	*uh = append(*uh,&UserHero{
		UserId:userId,
		HeroId:1,
		HeroName:"须佐",
	})
	*uh = append(*uh,&UserHero{
		UserId:userId,
		HeroId:2,
		HeroName:"神荒",
	})
}
