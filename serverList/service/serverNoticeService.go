package service

import (
	"fmt"
	"serverList/dao"
	"serverList/enum"
	"serverList/model"
	"time"
)

type Notice struct {
	ServerId int  `json:"serverId"`
	Notice string `json:"notice"`
}

type NoticeOperation struct {
	Action int8  //记录行为
	ServerId int
	Notice string
}
type NoticeOperationResult struct {
	Result int8  //结果
	ServerId int
	Notice string
}
// NoticeOperationChan 记录map行为的chan
var NoticeOperationChan 		= make(chan NoticeOperation)
// NoticeOperationResultChan 获取map行为返回的chan
var NoticeOperationResultChan 	= make(chan NoticeOperationResult)

var ServerNoticeMap 			= make(map[int]string)
// AllNotice 便于监控设置为全局
var AllNotice []*model.ServerNotice

func InitServerNotice() error {
	var err error
	err, AllNotice = dao.GetAllNotEndServerNoticeData()
	if err != nil {
		return err
	}
	fmt.Println("start init notice")
	fmt.Println(AllNotice)
	fmt.Println("serviceNoticeMap处理协程启动")
	go OperateNoticeMap()
	if AllNotice != nil{
		var noticeOperation NoticeOperation
		for _,noticeItem  := range AllNotice{
			noticeOperation.Action 		= enum.MAP_ADD_OR_UPDATE
			noticeOperation.ServerId 	= noticeItem.ServerId
			noticeOperation.Notice 		= noticeItem.Notice
			NoticeOperationChan <- noticeOperation
			<-NoticeOperationResultChan
		}
	}
	fmt.Println("end init notice")

	//监控notice是否过期了
	go monitorExpireNotice()

	return nil
}

func EndServerNotice(serverId int){
	for i:=0 ; i < len(AllNotice); i++  {
		if serverId == AllNotice[i].ServerId {
			fmt.Printf("外部请求删除公告:serverId%d\n",serverId)
			AllNotice[i].IsEnd = enum.NOTICE_IS_END_YES
			goto endFor
		}
	}
endFor:
	return
}

func OperateNoticeMap(){
	var noticeOperationResult NoticeOperationResult
	for  {
		select {
		case noticeOperation := <-NoticeOperationChan:
			noticeOperationResult.Result 	= enum.STATUS_SUCC
			noticeOperationResult.ServerId 	= 0
			noticeOperationResult.Notice 	= ""
			if noticeOperation.Action == enum.MAP_ADD_OR_UPDATE {
				fmt.Printf("新增编辑公告:serverId=>%d,公告=>%s\n",noticeOperation.ServerId,noticeOperation.Notice)
				noticeOperationResult.ServerId 	= noticeOperation.ServerId
				noticeOperationResult.Notice 	= noticeOperation.Notice
				//NoticeOperation struck
				ServerNoticeMap[noticeOperation.ServerId] = noticeOperation.Notice
			}
			if noticeOperation.Action == enum.MAP_GET {
				noticeOperationResult.ServerId 	 = noticeOperation.ServerId
				notice, checkNotice 			:= ServerNoticeMap[noticeOperation.ServerId]
				if checkNotice != false{
					noticeOperationResult.Notice = notice
				}
				fmt.Println(ServerNoticeMap)
				fmt.Printf("获取公告:serverId=>%d,公告=>%s\n",noticeOperationResult.ServerId,noticeOperationResult.Notice)
			}
			if noticeOperation.Action == enum.MAP_DELETE {
				fmt.Printf("删除公告:serverId%d\n",noticeOperation.ServerId)
				delete(ServerNoticeMap,noticeOperation.ServerId)
			}
			NoticeOperationResultChan <- noticeOperationResult
		}
	}
}

func monitorExpireNotice()  {
	tickerDb 			:= time.NewTicker(2 * time.Second)
	for _,noticeItem  	:= range AllNotice{
		go func(notice *model.ServerNotice) {
			for  {
				select {
				case <-tickerDb.C:
					now := time.Now().Unix()
					if notice.EndTime < uint64(now) || notice.IsEnd == enum.NOTICE_IS_END_YES{
						//请求后台更新状态
						//args := make([]string, 0)
						//args  = append(args,"ServerId=1")
						//response,_     := http.Get(config.NoticeUrl + "?" + strings.Join(args,"&"))
						//body, _        := ioutil.ReadAll(response.Body)
						//fmt.Println(string(body))
						//这边也不用删带，定义结构体的时候带上结束时间即可
						noticeOperation := NoticeOperation{
							Action 		: enum.MAP_DELETE,
							ServerId 	: notice.ServerId,
						}
						NoticeOperationChan <- noticeOperation
						<-NoticeOperationResultChan
						fmt.Println("公告删除结果")
						fmt.Println(ServerNoticeMap)
						goto endNotice
					}
				default:
				}
			}
		endNotice:
			fmt.Printf("server : %v notice is end\n",notice.ServerId)
		}(noticeItem)
	}
}


