package rpc

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

type testService struct{}

func (testService) HelloAction(body []byte) (proto.Message, error) {
	return &RespBody{Msg: "test hello"}, nil
}

func TestDialWebsocket(t *testing.T) {
	srv := NewServer()
	srv.Register("test", testService{})
	s := httptest.NewServer(srv.WebsocketHandler([]string{"*"}))
	defer s.Close()
	defer srv.Close()

	wsURL := "ws:" + strings.TrimPrefix(s.URL, "http:")
	fmt.Println(wsURL)
	client, err := DialWebsocket(context.Background(), wsURL, "")
	if err != nil {
		t.Fatal(err)
	}

	req := &ReqMsg{
		Service: "test",
		Method:  "hello",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.Call(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	fmt.Println(resp)
}
