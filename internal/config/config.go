package config

type Bootstrap struct {
	Server Server
	Data
}

type Server struct {
	GB28181 GB28181
	Media Media
	HttpPort uint16
}

type GB28181 struct {
	Port  uint16
	SipId string
	SipDomain string
}

type Media struct{
	Addr string
	Port uint16
	StreamRecvPort uint16
}

type Data struct {
	Database Database
}

type Database struct {
	Address, UserName, Password, DBName, Driver string
	Timeout int
}
