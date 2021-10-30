package gitlab

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gitlab-booster/config"
	"gitlab-booster/upstream"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
)

// ClientTraits gitlab client 接口
type ClientTraits interface {
	GrantOAuthToken(code string, redirectURI string) (*GrantedToken, error)
}

// Client  gitlab client 实现
type Client struct {
	endpoint   string
	httpClient *http.Client
	validate   *validator.Validate
}

// NewClient 创建一个 gitlab Client 实例
func NewClient() *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{endpoint: config.GetGitlabEndpoint(), httpClient: &http.Client{Transport: tr}, validate: validator.New()}
}

type grantTokenReq struct {
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
	Code         string `url:"code"`
	GrantType    string `url:"grant_type"`
	RedirectURI  string `url:"redirect_uri"`
}

// GrantedToken gitlab oauth/token 接口的响应数据
type GrantedToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// GrantOAuthToken 获取 OAuth token
func (c *Client) GrantOAuthToken(code string, redirectURI string) (*GrantedToken, error) {
	method := http.MethodPost
	url := fmt.Sprintf("%s/oauth/token", c.endpoint)
	reqBody := grantTokenReq{
		ClientID:     config.GetGitlabOAuthClientID(),
		ClientSecret: config.GetGitlabOAuthClientSecret(),
		Code:         code,
		GrantType:    "authorization_code",
		RedirectURI:  redirectURI,
	}
	v, err := query.Values(&reqBody)
	if err != nil {
		panic(err) // 不应该发生
	}

	req, err := http.NewRequest(method, url, strings.NewReader(v.Encode()))
	if err != nil {
		panic(err) // 不应该发生
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)

	// An error is returned if caused by client policy (such as
	// CheckRedirect), or failure to speak HTTP (such as a network
	// connectivity problem). A non-2xx status code doesn't cause an
	// error.
	if err != nil {
		panic(upstream.ErrUpstreamServiceConnectivity{Service: c.endpoint, Request: method + " " + url, Err: err})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(upstream.ErrUpstreamServiceConnectivity{Service: c.endpoint, Request: method + " " + url, Err: err})
	}

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("%s %s: responsed with unexpected status: %d %s", method, url, resp.StatusCode, string(body)))
	}

	log.Println(string(body))
	respBody := GrantedToken{}
	if err := json.Unmarshal(body, &respBody); err != nil {
		panic(err) // 不应该发生
	}
	return &respBody, nil
}
