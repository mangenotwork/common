package grpcHelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/url"
)

func PostProto[T proto.Message](host, path string, req proto.Message, reps T, header http.Header) (*T, error) {
	parse, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	parse.Path = path
	body, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("POST", parse.String(), bytes.NewBuffer(body))
	request.Header.Add("Content-Type", "application/x-protobuf")
	for key, values := range header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	do, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}
	data, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(data, reps)
	if err != nil {
		return nil, err
	}
	return &reps, nil
}

func ProtoToString(val proto.Message) string {
	marshal, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		return ""
	}
	return string(marshal)
}
