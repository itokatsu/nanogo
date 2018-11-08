package catplugin

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/itokatsu/nanogo/botutils"
)

var baseUrl = "http://thecatapi.com/api/images/get?format=xml"

type catPlugin struct {
	name string
}

func (p *catPlugin) Name() string {
	return "cat"
}

func New() (*catPlugin, error) {
	var pInstance catPlugin
	return &pInstance, nil
}

type Response struct {
	Images []CatImg `xml:"data>images>image"`
}

type CatImg struct {
	Url        string `xml:"url"`
	Id         string `xml:"id"`
	Source_url string `xml:"source_url"`
}

func (p *catPlugin) HandleMsg(cmd *botutils.Cmd, s *discordgo.Session, m *discordgo.MessageCreate) {
	switch strings.ToLower(cmd.Name) {
	case "cat":
		reqUrl := fmt.Sprintf("%s&results_per_page=5", baseUrl)
		res := Response{}
		botutils.FetchXML(reqUrl, &res)
		for _, img := range res.Images {
			_, err := botutils.Client.Get(img.Url)
			if err != nil {
				continue
			}
			s.ChannelMessageSend(m.ChannelID, img.Url)
			return
		}
	}
}

func (p *catPlugin) Help() string {
	return "Get a random cat picture"
}
