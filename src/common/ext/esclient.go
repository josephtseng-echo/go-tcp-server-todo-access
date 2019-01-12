package ext

import (
	"common/service"
	"strconv"
)

type EsClientManager struct {
	Clients map[string]service.EasyClientInterface
}

var EsClient = &EsClientManager{
	Clients: make(map[string]service.EasyClientInterface),
}

func AddEsClient(ip string, port int32) service.EasyClientInterface {
	k := ip + ":" + strconv.Itoa(int(port))
	cli := service.NewTcpClient(k)
	EsClient.Clients[k] = cli
	return cli
}

func GetEsClient(ip string, port int32) service.EasyClientInterface {
	k := ip + ":" + strconv.Itoa(int(port))
	ret, ok := EsClient.Clients[k]
	if ok {
		return ret
	}
	return nil
}

func RemoveEsClient(ip string, port int32) {
	k := ip + ":" + strconv.Itoa(int(port))
	_, ok := EsClient.Clients[k]
	if ok {
		delete(EsClient.Clients, k)
	}
}
