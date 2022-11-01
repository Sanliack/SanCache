package cachehttp

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	"sancache/ByteCont"
	pd "sancache/grpc"
)

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(in *pd.Request, out *pd.Response) error {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body:%v", err)
	}
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body :%v", err)
	}
	return nil
}

var _ ByteCont.PeerGetter = (*httpGetter)(nil)

var _PeerGetter = (*httpGetter)(nil)
