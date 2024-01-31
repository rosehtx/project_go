package service

import (
	"serverList/enum"
	"serverList/model"
)

type Notice struct {
	ServerId int  `json:"serverId"`
	Notice string `json:"notice"`
}

var ServerNoticeMap = make(map[int]string)

var AllNotice []model.ServerNotice

//func InitServerNotice() (bool, string) {
//	var result *gorm.DB
//	result, AllNotice = daoNotice.GetAllNotEndServerNoticeData()
//	if result.Error != nil {
//		return false, result.Error.Error()
//	}
//	fmt.Println("start init notice")
//	if AllNotice != nil{
//		for i := 0; i < len(AllNotice); i++ {
//			//先注册到map里
//			ServerNoticeMap[AllNotice[i].ServerId] = AllNotice[i].Notice
//			go func(i int) {
//				tickerDb := time.NewTicker(2 * time.Second)
//				for  {
//					select {
//					case <-tickerDb.C:
//						now := time.Now().Unix()
//						if AllNotice[i].EndTime < uint64(now) || AllNotice[i].IsEnd == enum.NOTICE_IS_END_YES{
//							//请求后台更新状态
//							args := make([]string, 0)
//							args = append(args,"ServerId=1")
//							response,_     := http.Get(config.NoticeUrl + "?" + strings.Join(args,"&"))
//							body, _        := ioutil.ReadAll(response.Body)
//							fmt.Println(string(body))
//							delete(ServerNoticeMap,AllNotice[i].ServerId)
//							goto endNotice
//						}
//					default:
//					}
//				}
//			endNotice:
//				fmt.Printf("server : %v notice is end",AllNotice[i].ServerId)
//			}(i)
//		}
//	}
//	fmt.Println("end init notice")
//	return true, ""
//}

func EndServerNotice(serverId int) bool{
	for i:=0 ; i < len(AllNotice); i++  {
		if serverId == AllNotice[i].ServerId {
			AllNotice[i].IsEnd = enum.NOTICE_IS_END_YES
			goto endFor
		}
	}
endFor:
	return true
}


