package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hypefast.io/services/config/env"
	shorterModel "hypefast.io/services/domain/shorter/model"
	"hypefast.io/services/pkg/utils/generate"
	"hypefast.io/services/pkg/utils/logging"
)

func (i *ShorterUseCaseImpl) SetURL(payload shorterModel.Payload) (interface{}, error) {

	rdb := i.Conf.RedisConnect

	key := generate.Encode(uint64(time.Now().UnixNano()))

	result := shorterModel.Return{
		URL:      payload.URL,
		ShortURL: fmt.Sprintf("%s/shorter/%s", env.Conf.DomainName, key),
	}

	data, err := json.Marshal(result)

	if err != nil {
		logging.Errorf("error: %v", err)
		return nil, err
	}

	rdbSet := rdb.Set(context.Background(), key, data, time.Duration(time.Hour*24))
	if err := rdbSet.Err(); err != nil {
		logging.Errorf("unable to SET data. error: %v", err)
		return nil, err
	}

	return result, nil
}
