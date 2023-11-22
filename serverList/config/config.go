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
	RMQ_IP 		= "172.19.60.19"
	RMQ_PORT 	= 5672
	RMQ_VHOST 	= "Testscm2"
	RMQ_USER 	= "dbadmin"
	RMQ_PASS 	= "UatdbaduserPwd.8263"
)
