package fcapiprovider

import (
	"context"
	"log"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
)

type LotusApiAccessLevel int

const (
	Admin     LotusApiAccessLevel = 0
	Read                          = 1
	Write                         = 2
	ReadWrite                     = 3
)

//some interface here?

type LotusClient struct {
	id          int
	name        string
	ipAddr      string
	port        string
	token       string
	accessLevel LotusApiAccessLevel
	closer      jsonrpc.ClientCloser
}

type LotusFullNodeClient struct {
	LotusClient
	FullNodeApi apistruct.FullNodeStruct
}

func NewLotusFullNodeClient(ctx context.Context, name string, ipAddr string, port string, token string, accessLevel LotusApiAccessLevel) (*LotusFullNodeClient, error) {
	//log.Println("XXXX - " + name)
	lfnc := LotusFullNodeClient{
		FullNodeApi: apistruct.FullNodeStruct{},
	}
	lfnc.ipAddr = ipAddr
	lfnc.port = port
	lfnc.token = token
	lfnc.name = name
	lfnc.accessLevel = accessLevel
	closer, err := jsonrpc.NewMergeClient(ctx, "ws://"+lfnc.ipAddr+":"+lfnc.port+"/rpc/v0", "Filecoin", []interface{}{&lfnc.FullNodeApi.Internal, &lfnc.FullNodeApi.CommonStruct.Internal}, http.Header{"Authorization": []string{"Bearer " + lfnc.token}})
	lfnc.closer = closer

	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	//log.Println("XXXX - " + closer)
	return &lfnc, err
}

func (lfnc *LotusFullNodeClient) CloseClient() {
	lfnc.closer()
}
