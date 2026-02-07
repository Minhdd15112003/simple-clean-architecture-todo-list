package restapi

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type itemService struct {
	client     *resty.Client
	serviceURL string
}

func New(serviceURL string) *itemService {
	return &itemService{
		serviceURL: serviceURL,
		client:     resty.New(),
	}
}

func (s *itemService) GetLikeItem(ctx context.Context, ids []int) (map[int]int, error) {
	type requestBody struct {
		Ids []int `json:"ids"`
	}
	var response struct {
		Data map[int]int `json:"data"`
	}

	res, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody{Ids: ids}).
		SetResult(&response).
		Post(fmt.Sprintf("%s/%s", s.serviceURL, "v1/rpc/item-like"))

	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		log.Println(res.RawResponse)
		return nil, errors.New("cannot call api get item likes")
	}
	return response.Data, nil
}
