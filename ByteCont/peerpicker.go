package ByteCont

import pd "sancache/grpc"

type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

//type PeerGetter interface {
//	Get(group string, key string) ([]byte, error)
//}

type PeerGetter interface {
	Get(group *pd.Request, key *pd.Response) error
}
