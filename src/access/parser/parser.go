package parser

var BUFREADTOTALSIZE = 4096
var PACKETCHANSIZE = 100
var ERRORCHANSIZE = 30
var PACKETTOTALSIZE = 10240

//ES
var ESHEADSIZE = 8

type EsHead struct {
	Flag [2]byte
	Cmd  uint16
	Ver  uint16
	Len  uint16
}

//RPC
var RPCHEADSIZE = 6

type RpcHead struct {
	Flag [3]byte
	Ver  uint8
	Len  uint16
}