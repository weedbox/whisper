package whisper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	DefaultAPIRule = "/apis/v1/group/:group_id/members/id"
)

type GetGroupResponse struct {
	Members []string `json:"members"`
}

type GroupResolverHTTP struct {
	w       Whisper
	url     string
	apiRule string
}

func NewGroupResolverHTTP(url string, apiRule string) *GroupResolverHTTP {
	return &GroupResolverHTTP{
		url:     url,
		apiRule: apiRule,
	}
}

func (gs *GroupResolverHTTP) Init(w Whisper) error {
	gs.w = w
	return nil
}

func (gs *GroupResolverHTTP) formatURL(template string, params map[string]string) string {
	for k, v := range params {
		placeholder := fmt.Sprintf(":%s", k)
		template = strings.Replace(template, placeholder, v, 1)
	}
	return template
}

func (gs *GroupResolverHTTP) GetGroupRule(groupID string) GroupRule {

	url := gs.formatURL(fmt.Sprintf("%s%s", gs.url, gs.apiRule), map[string]string{
		"group_id": groupID,
	})

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer req.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("group manager service is not available")
		return nil
	}

	var gres GetGroupResponse
	err = json.Unmarshal(respData, &gres)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	r := NewGroupRule()
	r.AddMembers(gres.Members)

	return r
}
