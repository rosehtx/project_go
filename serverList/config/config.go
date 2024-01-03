package config

var (
	MysqlIp string		= "127.0.0.1"
	MysqlPort int		= 3306
	MysqlUser string	= "root"
	MysqlPass string	= "root"

	NoticeUrl string    = "http://192.168.44.127/test/test.php"
)

//rabbitmq配置
const (
	RMQ_IP 			= "172.19.60.19"
	RMQ_PORT 		= 5672
	RMQ_VHOST 		= "Testscm2"
	RMQ_USER 		= "dbadmin"
	RMQ_PASS 		= "UatdbaduserPwd.8263"
	RMQ_CON_NUM 	= 5 //rmq 连接池数量
	RMQ_CHANNEL_NUM = 2//rmq 每个连接对应channel数量
	RMQ_QOS 		= 2//qos 一个消费者同时消费的消息数量
	RMQ_CONSUME_NUM = 1//每个队列对应消费者的数量,最好是1对1消费
)
