package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const ApiLogsEndpoint string = "api_logs"

type ApiLogRequest struct {
	client *Client
}

type HttpMethod string

const (
	PostMethod		HttpMethod = "post"
	PutMethod			HttpMethod = "put"
	DeleteMethod	HttpMethod = "delete"
)

type ApiLog struct {
	RequestId       uuid.UUID              	`json:"request_id"`
	ApiVersion      string                 	`json:"api_version"`
	Client         	string     							`json:"client"`
	HttpMethod 			HttpMethod 							`json:"http_method"`
	HttpStatus			int 										`json:"http_status"`
	RequestBody  		map[string]interface{} 	`json:"request_body,omitempty"`
	RequestOrigin		string                 	`json:"request_origin"`
	RequestPath     string              		`json:"request_path"`
	RequestResponse	map[string]interface{} 	`json:"request_response,omitempty"`
	LoggedAt        time.Time              	`json:"logged_at"`
	CreatedAt       time.Time              	`json:"created_at"`
}

type ApiLogListInput struct {
	PerPage 			int 		 `json:"per_page,omitempty,string"`
	Page    			int 		 `json:"page,omitempty,string"`
	FromDate      string   `json:"from_date,omitempty"`
	ToDate        string   `json:"to_date,omitempty"`
	HttpMethods   []string `json:"http_methods,omitempty"`
	HttpStatuses	[]string `json:"http_statuses,omitempty"`
	ApiVersion		[]string `json:"api_version,omitempty"`
	RequestPaths  []string `json:"request_paths,omitempty"`
}

type ApiLogResult struct {
	ApiLog  	*ApiLog  `json:"api_log,omitempty"`
	ApiLogs 	[]ApiLog `json:"api_logs,omitempty"`
	Meta      Metadata `json:"meta,omitempty"`
}

func (c *Client) ApiLog() *ApiLogRequest {
	return &ApiLogRequest{
		client: c,
	}
}

func (alr *ApiLogRequest) Get(ctx context.Context, RequestId string) (*ApiLog, *Error) {
	subPath := fmt.Sprintf("%s/%s", ApiLogsEndpoint, RequestId)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ApiLogResult{},
	}

	result, err := alr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	apiLogResult, ok := result.(*ApiLogResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return apiLogResult.ApiLog, nil
}

func (alr *ApiLogRequest) GetList(ctx context.Context, apiLogListInput *ApiLogListInput) (*ApiLogResult, *Error) {
	jsonQueryParams, err := json.Marshal(apiLogListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        ApiLogsEndpoint,
		QueryParams: queryParams,
		Result:      &ApiLogResult{},
	}

	result, clientErr := alr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	apiLogResult, ok := result.(*ApiLogResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return apiLogResult, nil
}
