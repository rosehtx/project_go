package service

import (
	"context"
	"errors"
	"heroServer/model"
	"heroServer/proto/rpcPb/equip"
	"heroServer/proto/rpcPb/hero"
)

//用来实现rpc里的 GetUserHero 服务
type GetUserHeroServer struct {

}

func (server *GetUserHeroServer)  GetUserHero(ctx context.Context, requestData *hero.GetUserHeroRequest) (*hero.GetUserHeroResponse, error){
	userId := requestData.GetId()

	//返回的参数
	var response *hero.GetUserHeroResponse
	response = new(hero.GetUserHeroResponse)

	if userId == 0 {
		return response,errors.New("用户不能为空")
	}

	response.Name = "get user name"
	response.Age  = userId

	//用户的英雄数据
	var respHero []*hero.Hero
	//获取用户的英雄数据
	uh 			:= &model.AllUserHero{}
	uh.GetHeroByUserId(uint64(userId))

	//英雄数值
	var heroNumericalValue []*hero.NumericalValue
	//英雄装备
	var heroEquip []*equip.Equip
	for _, v := range *uh {
		//数值赋值 根据业务自行处理
		heroNumericalValue = append(heroNumericalValue,&hero.NumericalValue{
			Type: hero.NumericalValueType_NV_OFFENCE,
			Value:1,
		})
		heroNumericalValue = append(heroNumericalValue,&hero.NumericalValue{
			Type: hero.NumericalValueType_NV_DEFENSE,
			Value:2,
		})
		//装备
		heroEquip= append(heroEquip,&equip.Equip{
			Id: 1,
			Name:"天丛云",
			EquipType:equip.EquipType_UP,
		})
		respHero = append(respHero, &hero.Hero{
			Id: (*v).HeroId,
			Name: (*v).HeroName,
			UserId: (*v).UserId,
			HeroType: hero.HeroType_OFFENCE,
			NumericalValue:heroNumericalValue,
			Equip: heroEquip,
		})
		//置为空
		heroNumericalValue  = []*hero.NumericalValue{}
		heroEquip			= []*equip.Equip{}
	}
	response.Hero = respHero
	return response,nil
}


