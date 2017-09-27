// Nanobot Project
//
// Plugin for google searches

/*
@TODO: googleimg ~ collage
@TODO: cmd with more results
*/

package googleplugin

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/itokatsu/nanogo/botutils"
)

type googlePlugin struct {
	name     string
	apiKey   string
	lastReqs map[string]*url.URL
}

func New(apiKey string) *googlePlugin {
	var pInstance googlePlugin
	pInstance.apiKey = apiKey
	pInstance.lastReqs = make(map[string]*url.URL)
	return &pInstance
}

func (p *googlePlugin) Name() string {
	return "google"
}

func (p *googlePlugin) HasSaves() bool {
	return false
}

type SearchResults struct {
	Items []Result
}
type Result struct {
	Title   string
	Link    string
	Snipper string
}

func (p *googlePlugin) buildRequestURL(query string, n int) url.URL {
	numStr := strconv.Itoa(n)
	qs := url.Values{}
	qs.Set("key", p.apiKey)
	qs.Set("cx", "004895194701224026743:zdbrbrrm0bw")
	qs.Set("client", "google-csbe")
	qs.Set("num", numStr)
	qs.Set("ie", "utf8")
	qs.Set("oe", "utf8")
	qs.Set("fields", "items(title,link,snippet)")
	query = strings.Replace(query, "/", "", -1)
	query = strings.Replace(query, "&", "", -1)
	qs.Set("q", query)

	var reqUrl = url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "customsearch/v1",
		RawQuery: qs.Encode(),
	}
	return reqUrl
}

func (p *googlePlugin) HandleMsg(cmd *botutils.Cmd, s *discordgo.Session, m *discordgo.MessageCreate) {
	switch strings.ToLower(cmd.Name) {

	//first google result
	case "g", "google":
		if len(cmd.Args) == 0 {
			return
		}
		query := strings.Join(cmd.Args, " ")
		url := p.buildRequestURL(query, 1)

		p.lastReqs[m.ChannelID] = &url
		result := SearchResults{}
		err := botutils.FetchJSON(url.String(), &result)
		if err != nil {
			return
		}

		if len(result.Items) > 0 {
			s.ChannelMessageSend(m.ChannelID, result.Items[0].Link)
		} else {
			msg := fmt.Sprintf("No result found for %v", strings.Join(cmd.Args, " "))
			s.ChannelMessageSend(m.ChannelID, msg)
		}

	//google img
	case "gi":
		if len(cmd.Args) == 0 {
			return
		}
		query := strings.Join(cmd.Args, " ")
		url := p.buildRequestURL(query, 9)
		url.Query().Set("searchType", "image")

		p.lastReqs[m.ChannelID] = &url
		results := []SearchResults{}
		err := botutils.FetchJSON(url.String(), &results)
		if err != nil {
			return
		}

	//more results
	case "gm":
		if _, ok := p.lastReqs[m.ChannelID]; !ok {
			return
		}
	}
}

func (p *googlePlugin) Help() string {
	return `
	!g <term> - Return first search results
	`
}

func (p *googlePlugin) Save() []byte {
	return nil
}

func (p *googlePlugin) Load(data []byte) error {
	return nil
}

func (p *googlePlugin) Cleanup() {
}
