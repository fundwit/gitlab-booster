package servehttp

import (
	"errors"
	"gitlab-booster/config"
	"gitlab-booster/upstream/gitlab"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type gitlabV4AuthHandler struct {
	validator *validator.Validate
	client    gitlab.ClientTraits
}

// RegisterGitlabV4AuthHandlers 注册路由
func RegisterGitlabV4AuthHandlers(root *gin.Engine, middlewares ...gin.HandlerFunc) {
	r := root.Group("/v1/gitlabv4/")
	r.Use(middlewares...)

	handler := &gitlabV4AuthHandler{
		validator: validator.New(),
		client:    gitlab.NewClient(),
	}

	r.POST("approves", handler.newGitlabV4Approve)
	r.POST("sessions", handler.newGitlabV4Session)
}

const stateExpiration = 24 * time.Hour

var stateCache = cache.New(stateExpiration, 1*time.Minute)

type newApproveReq struct {
	RedirectURI string `json:"redirectUri" validate:"required"`
}
type newApproveResp struct {
	State        string `json:"state"`
	ProviderURI  string `json:"providerUri"`
	ClientID     string `json:"clientId"`
	ResponseType string `json:"responseType"`

	Endpoint string `json:"endpoint"`
}
type newSessionReq struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}
type newSessionResp struct {
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
}

func (m *gitlabV4AuthHandler) newGitlabV4Approve(c *gin.Context) {
	reqBody := newApproveReq{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		panic(err)
	}
	if err := m.validator.Struct(reqBody); err != nil {
		panic(err)
	}

	state := uuid.New().String()
	q, _ := url.ParseQuery("redirectUri=" + reqBody.RedirectURI)
	stateCache.Set(state, q.Get("redirectUri"), cache.DefaultExpiration)
	c.JSON(201, &newApproveResp{
		State:        state,
		ProviderURI:  config.GetGitlabEndpoint() + "/oauth/authorize",
		ClientID:     config.GetGitlabOAuthClientID(),
		ResponseType: "code",
		Endpoint:     config.GetGitlabEndpoint(),
	})
}

func (m *gitlabV4AuthHandler) newGitlabV4Session(c *gin.Context) {
	reqBody := newSessionReq{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		panic(err)
	}
	if err := m.validator.Struct(reqBody); err != nil {
		panic(err)
	}

	value, found := stateCache.Get(reqBody.State)
	if !found {
		panic(errors.New("state " + reqBody.State + " is invalid"))
	}

	redirectURL, _ := value.(string)
	token, err := m.client.GrantOAuthToken(reqBody.Code, redirectURL)
	if err != nil {
		panic(err)
	}
	c.JSON(201, &newSessionResp{Token: token.AccessToken, Endpoint: config.GetGitlabEndpoint()})
}
