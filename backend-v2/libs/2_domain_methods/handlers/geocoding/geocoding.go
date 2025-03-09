package geocoding

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oschwald/geoip2-golang"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

var db *geoip2.Reader

func GeoCodingHandler(sctx smart_context.ISmartContext, params types.ANY_DATA) (interface{}, error) {
	id, ok := params.GetStringValue("id")
	if !ok || id == "" {
		return nil, fmt.Errorf("missing device id")
	}

	var metric model.Metric
	err := sctx.GetDB().Where("device_id = ?", id).Order("created_at DESC").First(&metric).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении устройства: %w", err)
	}

	geoURL := fmt.Sprintf("http://ip-api.com/json/%s", metric.PublicIP)
	resp, err := http.Get(geoURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении устройства: %w", err)
	}

	var loc Location
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return nil, fmt.Errorf("ошибка декодирования геоданных: %w", err)
	}

	if loc.Lat == 0 || loc.Lon == 0 {
		return nil, fmt.Errorf("геокодирование не удалось: %w", err)
	}

	return loc, nil
}
