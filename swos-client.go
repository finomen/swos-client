package swos_client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/icholy/digest"
)

type SwOsClient struct {
	client http.Client
	url    string

	Links LinkPage
	Sfp   SfpPage
	Sys   SysPage
	Rstp  RstpPage
	Fwd   FwdPage
	Vlan  VlanPage
}

func get(client *SwOsClient, path string) (*http.Response, error) {
	res, err := client.client.Get(client.url + path)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Request failed")
	}

	return res, err
}

func post(client *SwOsClient, path string, request string) (*http.Response, error) {
	res, err := client.client.Post(client.url+path, "text/plain", strings.NewReader(request))

	if res.StatusCode != 200 {
		return nil, errors.New("Request failed")
	}

	return res, err
}

func getPage[I any, P swOsPage[I]](client *SwOsClient, page P) error {
	resp, err := get(client, page.url())
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	dec := json.NewDecoder(strings.NewReader(fixJson(string(body))))
	var inp I
	err = dec.Decode(&inp)
	if err != nil {
		return err
	}

	err = page.load(inp)
	if err != nil {
		return err
	}

	return nil

}

func savePage[I any, O any, P writebleSwOsPage[O, I]](client *SwOsClient, page P) error {
	_, err := post(client, page.url(), anyToMikrotik(page.store()))
	if err != nil {
		return err
	}
	return getPage(client, page)
}

func NewSwOsClient(url string, user string, password string) (*SwOsClient, error) {
	client := &SwOsClient{
		client: http.Client{
			Transport: &digest.Transport{
				Username: user,
				Password: password,
			},
		},
		url: url,
	}

	res, err := get(client, "/link.b")
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("Auth failed")
	}

	return client, nil
}

func (client *SwOsClient) Fetch() error {
	err := getPage(client, &client.Links)
	if err != nil {
		return err
	}
	err = getPage(client, &client.Sfp)
	if err != nil {
		return err
	}
	client.Sys.numPorts = len(client.Links.Links)
	err = getPage(client, &client.Sys)
	if err != nil {
		return err
	}
	client.Rstp.numPorts = len(client.Links.Links)
	err = getPage(client, &client.Rstp)
	if err != nil {
		return err
	}
	client.Fwd.numPorts = len(client.Links.Links)
	err = getPage(client, &client.Fwd)
	if err != nil {
		return err
	}
	client.Vlan.numPorts = len(client.Links.Links)
	err = getPage(client, &client.Vlan)
	if err != nil {
		return err
	}

	return nil
}

func (client *SwOsClient) Save() error {
	err := savePage(client, &client.Links)
	if err != nil {
		return err
	}
	err = savePage(client, &client.Sys)
	if err != nil {
		return err
	}
	err = savePage(client, &client.Rstp)
	if err != nil {
		return err
	}
	err = savePage(client, &client.Fwd)
	if err != nil {
		return err
	}
	err = savePage(client, &client.Vlan)
	if err != nil {
		return err
	}
	return nil
}
