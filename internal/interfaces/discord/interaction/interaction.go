package interaction

import (
	"fmt"
	"strings"

	"tsniper/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

type HandleFunc func(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error)

func GetUserID(i *discordgo.InteractionCreate) string {
	if i.Member != nil {
		return i.Member.User.ID
	}
	return i.User.ID
}

func ParseInteractionOptions(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.ApplicationCommandData().Options))
	for _, opt := range i.ApplicationCommandData().Options {
		options[opt.Name] = opt
	}
	return options
}

func EncodeCustomID[T any](customID string, data ...utils.KeyValue[string, T]) string {
	if len(data) > 0 {
		params := make([]string, 0, len(data))
		for _, pair := range data {
			params = append(params, pair.Key+"="+fmt.Sprint(pair.Value))
		}
		customID += "?" + strings.Join(params, "&")
	}

	return customID
}

func DecodeCustomID(customID string) (string, map[string]string) {
	routingKey := customID
	params := make(map[string]string)
	if key, query, ok := strings.Cut(customID, "?"); ok {
		routingKey = key
		for pair := range strings.SplitSeq(query, "&") {
			if k, v, ok := strings.Cut(pair, "="); ok {
				params[k] = v
			}
		}
	}

	return routingKey, params
}
