package usecase

import (
	"context"
	"encoding/json"

	shorterModel "hypefast.io/services/domain/shorter/model"
	"hypefast.io/services/pkg/utils/logging"
)

func (i *ShorterUseCaseImpl) GetURL(key string) (string, error) {
	rdb := i.Conf.RedisConnect
	var result shorterModel.Return

	rdbGet, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		logging.Errorf("unable to GET data. error: %v", err)
		return "", err
	}

	if err := json.Unmarshal([]byte(rdbGet), &result); err != nil {
		return "", err
	}

	return result.URL, nil
}
