package whisper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
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

func NewGroupResolverHTTP(url string, apiRule string) GroupResolver {
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

func (gs *GroupResolverHTTP) GetMemberIDs(groupID string) ([]string, error) {

	url := gs.formatURL(fmt.Sprintf("%s%s", gs.url, gs.apiRule), map[string]string{
		"group_id": groupID,
	})

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {

		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrGroupNotFound
		}

		return nil, ErrOperationFailure
	}

	var gres GetGroupResponse
	err = json.Unmarshal(respData, &gres)
	if err != nil {
		return []string{}, nil
	}

	return gres.Members, nil
}
