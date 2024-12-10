package rest

// Import resty into your code and refer it as `resty`.
import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Rest struct {
	http   *resty.Client
	token  *Token
	config map[string]interface{}
}

func (r *Rest) getHttp() *resty.Client {
	return r.http
}

func (r *Rest) getToken() (*Token, error) {
	if r.token == nil {
		return nil, errors.New("missing authentication token")
	}
	if !r.token.IsValid() {
		r.token = nil
		return nil, errors.New("invalid authentication token")
	}
	return r.token, nil
}

// defines a token to be used in the requests
func (r *Rest) SetToken(token *Token) error {
	if !token.IsValid() {
		return errors.New("token is invalid")
	}
	r.token = token
	return nil
}

// sets variables used in the requests
func (r *Rest) SetConfig(key string, value string) {
	r.config[key] = value
}

// gets variables used in the requests
func (r *Rest) GetConfig(key string) string {
	return r.config[key].(string)
}

// sets variables used in the requests
func (r *Rest) GetConfigData() map[string]interface{} {
	return r.config
}

// posts request to the given link, using the defined token
func (r *Rest) Post(payload map[string]interface{}, link string) (*Response, error) {
	token, err := r.getToken()
	if err != nil {
		return nil, err
	}
	resp, err := r.getHttp().R().SetBody(payload).SetAuthToken(token.GetKey()).Post(link)
	if err != nil {
		return nil, err
	}
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

// posts request to the given link, using the defined token
func (r *Rest) PostArray(payload []map[string]interface{}, link string) (*Response, error) {
	token, err := r.getToken()
	if err != nil {
		return nil, err
	}
	resp, err := r.getHttp().R().SetBody(payload).SetAuthToken(token.GetKey()).Post(link)
	if err != nil {
		return nil, err
	}
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

// posts request to the given link, using the defined token and context
func (r *Rest) PostWithContext(payload map[string]interface{}, link string, ctx context.Context) (*Response, error) {
	token, err := r.getToken()
	if err != nil {
		return nil, err
	}
	resp, err := r.getHttp().R().SetContext(ctx).SetBody(payload).SetAuthToken(token.GetKey()).Post(link)
	if err != nil {
		return nil, err
	}
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

// posts request to the given link, using the defined token and specific header
func (r *Rest) PostWithHeader(payload map[string]interface{}, link string, header map[string]string) (*Response, error) {
	token, err := r.getToken()
	if err != nil {
		return nil, err
	}
	resp, err := r.getHttp().R().SetBody(payload).SetHeaders(header).SetAuthToken(token.GetKey()).Post(link)
	if err != nil {
		return nil, err
	}
	resp.Time()
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

// posts request to the given link, without token and specific header
func (r *Rest) PostWithHeaderNoAuth(payload map[string]interface{}, link string, header map[string]string) (*Response, error) {
	resp, err := r.getHttp().R().SetBody(payload).SetHeaders(header).Post(link)
	if err != nil {
		return nil, err
	}
	resp.Time()
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

// gets request to the given link, using the defined token
func (r *Rest) Get(payload map[string]interface{}, link string) (*Response, error) {
	token, err := r.getToken()
	if err != nil {
		return nil, err
	}
	data := r.preparePayload(payload)
	resp, err := r.getHttp().R().SetQueryParams(data).SetAuthToken(token.GetKey()).Get(link)
	if err != nil {
		return nil, err
	}
	resp.Time()
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

// gets request to the given link, using the defined token and specific header
func (r *Rest) GetWithHeader(payload map[string]interface{}, link string, header map[string]string) (*Response, error) {
	token, err := r.getToken()
	if err != nil {
		return nil, err
	}
	data := r.preparePayload(payload)
	resp, err := r.getHttp().R().SetQueryParams(data).SetHeaders(header).SetAuthToken(token.GetKey()).Get(link)
	if err != nil {
		return nil, err
	}
	resp.Time()
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

// gets request to the given link, without token and specific header
func (r *Rest) GetWithHeaderNoAuth(payload map[string]interface{}, link string, header map[string]string) (*Response, error) {
	data := r.preparePayload(payload)
	resp, err := r.getHttp().R().SetQueryParams(data).SetHeaders(header).Get(link)
	if err != nil {
		return nil, err
	}
	resp.Time()
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

// deletes request to the given link, using the defined token
func (r *Rest) Delete(link string) (*Response, error) {
	token, err := r.getToken()
	if err != nil {
		return nil, err
	}
	resp, err := r.getHttp().R().SetAuthToken(token.GetKey()).Delete(link)
	if err != nil {
		return nil, err
	}
	return &Response{
		code: resp.StatusCode(),
		raw:  resp.String(),
	}, nil
}

func (r *Rest) preparePayload(payload map[string]interface{}) map[string]string {
	result := map[string]string{}
	for k, v := range payload {
		switch t := v.(type) {
		case string:
			result[k] = v.(string)
		case bool:
			if v.(bool) {
				result[k] = "true"
				continue
			}
			result[k] = "false"
		default:
			result[k] = fmt.Sprintf("%v", t)
		}
	}
	return result
}

// gets a Rest struct with the given config
// if InsecureSkipVerify is set to true, the client will skip the verification of the server's certificate
func NewRest(config map[string]interface{}) *Rest {
	client := resty.New()
	if config["InsecureSkipVerify"] != nil && config["InsecureSkipVerify"].(bool) {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: config["InsecureSkipVerify"].(bool)})
	}
	rest := &Rest{
		http:   client,
		config: config,
		token:  &Token{},
	}
	rest.http.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	rest.getHttp().SetTimeout(1 * time.Minute)
	return rest
}
