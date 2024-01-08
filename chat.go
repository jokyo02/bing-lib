package binglib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/Harry-zklcdc/bing-lib/lib/hex"
	"github.com/Harry-zklcdc/bing-lib/lib/request"
	"golang.org/x/net/websocket"
)

const (
	PRECISE          = "Precise"          // 精准
	BALANCED         = "Balanced"         // 平衡
	CREATIVE         = "Creative"         // 创造
	PRECISE_OFFLINE  = "Precise-offline"  // 精准, 不联网搜索
	BALANCED_OFFLINE = "Balanced-offline" // 平衡, 不联网搜索
	CREATIVE_OFFLINE = "Creative-offline" // 创造, 不联网搜索
)

const (
	bingCreateConversationUrl = "https://%s/turing/conversation/create?bundleVersion=1.1418.13"
	sydneyChatHubUrl          = "wss://%s/sydney/ChatHub?sec_access_token=%s"

	spilt = "\x1e"
)

func NewChat(cookies string) *Chat {
	return &Chat{
		cookies:       cookies,
		BingBaseUrl:   bingBaseUrl,
		SydneyBaseUrl: sydneyBaseUrl,
	}
}

func (chat *Chat) SetCookies(cookies string) *Chat {
	chat.cookies = cookies
	return chat
}

func (chat *Chat) SetStyle(style string) *Chat {
	chat.GetChatHub().SetStyle(style)
	return chat
}

func (chat *Chat) SetBingBaseUrl(bingBaseUrl string) *Chat {
	chat.BingBaseUrl = bingBaseUrl
	return chat
}

func (chat *Chat) SetSydneyBaseUrl(sydneyBaseUrl string) *Chat {
	chat.SydneyBaseUrl = sydneyBaseUrl
	return chat
}

func (chat *Chat) GetCookies() string {
	return chat.cookies
}

func (chat *Chat) GetChatHub() *ChatHub {
	return chat.chatHub
}

func (chat *Chat) GetStyle() string {
	return chat.GetChatHub().GetStyle()
}

func (chat *Chat) GetBingBaseUrl() string {
	return chat.BingBaseUrl
}

func (chat *Chat) GetSydneyBaseUrl() string {
	return chat.SydneyBaseUrl
}

func (chat *Chat) NewConversation() error {
	c := request.NewRequest()
	c.SetUrl(fmt.Sprintf(bingCreateConversationUrl, chat.BingBaseUrl)).
		SetHeader("Cookie", chat.cookies).
		SetHeader("Origin", "https://www.bing.com").
		SetHeader("Referer", "https://www.bing.com/search?q=Bing+AI&showconv=1&FORM=hpcodx&wlexpsignin=1&wlexpsignin=1").
		SetHeader("User-Agent", userAgent).
		SetHeader("X-Ms-Useragent", "azsdk-js-api-client-factory/1.0.0-beta.1 core-rest-pipeline/1.12.0 OS/Windows").
		Do()

	var resp ChatReq
	err := json.Unmarshal(c.GetBody(), &resp)
	if err != nil {
		return err
	}
	resp.ConversationSignature = c.GetHeader("X-Sydney-Conversationsignature")
	resp.EncryptedConversationSignature = c.GetHeader("X-Sydney-Encryptedconversationsignature")

	chat.chatHub = newChatHub(resp)

	return nil
}

func (chat *Chat) MsgComposer(msg string) string {
	// TODO
	return ""
}

func (chat *Chat) optionsSetsHandler() []string {
	optionsSets := []string{
		"nlu_direct_response_filter",
		"deepleo",
		"disable_emoji_spoken_text",
		"responsible_ai_policy_235",
		"enablemm",
		"dv3sugg",
		"autosave",
		"iyxapbing",
		"iycapbing",
		"rai289",
		"enflst",
		"enpcktrk",
		"rcaldictans",
		"rcaltimeans",
		"eredirecturl",
	}

	tone := chat.GetStyle()
	if tone == PRECISE || tone == PRECISE_OFFLINE {
		optionsSets = append(optionsSets, "h3precise", "clgalileo", "gencontentv3")
	} else if tone == BALANCED || tone == BALANCED_OFFLINE {
		optionsSets = append(optionsSets, "galileo", "saharagenconv5")
	} else if tone == CREATIVE || tone == CREATIVE_OFFLINE {
		optionsSets = append(optionsSets, "h3imaginative", "clgalileo", "gencontentv3")
	}
	return optionsSets
}

func (chat *Chat) pluginHandler(optionsSets *[]string) []Plugins {
	plugins := []Plugins{}
	tone := chat.GetStyle()
	if tone == PRECISE || tone == BALANCED || tone == CREATIVE {
		plugins = append(plugins, Plugins{Id: "c310c353-b9f0-4d76-ab0d-1dd5e979cf68"})
	} else {
		*optionsSets = append(*optionsSets, "nosearchall")
	}
	return plugins
}

func (chat *Chat) systemContextHandler(prompt string) []SystemContext {
	systemContext := []SystemContext{}
	if prompt != "" {
		systemContext = append(systemContext, SystemContext{
			Author:      "user",
			Description: prompt,
			ContextType: "WebPage",
			MessageType: "Context",
			MessageId:   "discover-web--page-ping-mriduna-----",
		})
	}
	return systemContext
}

func (chat *Chat) requestPayloadHandler(msg string, optionsSets []string, plugins []Plugins, systemContext []SystemContext) map[string]any {
	msgId := hex.NewUUID()
	tone := chat.GetStyle()

	data := map[string]any{
		"arguments": []any{
			map[string]any{
				"source":      "cib",
				"optionsSets": optionsSets,
				"allowedMessageTypes": []string{
					"ActionRequest",
					"Chat",
					"ConfirmationCard",
					"Context",
					"InternalSearchQuery",
					"InternalSearchResult",
					"Disengaged",
					"InternalLoaderMessage",
					"InvokeAction",
					"Progress",
					"RenderCardRequest",
					"RenderContentRequest",
					"AdsQuery",
					"SemanticSerp",
					"GenerateContentQuery",
					"SearchQuery",
				},
				"sliceIds": []string{
					"techpillscf",
					"gbaa",
					"gba",
					"gbapa",
					"codecreator",
					"dlidcf",
					"specedge",
					"preall15",
					"suppsm240-t",
					"translref",
					"ardsw_1_9_9",
					"fluxnosearchc",
					"fluxnosearch",
					"1115rai289",
					"1119backoss0",
					"124multi2t",
					"1129gpt4ts0",
					"kchero50cf",
					"cacfastapis",
					"cacdupereccf",
					"cacmuidarb",
					"cacfrwebt2cf",
					"sswebtop2cf",
				},
				"isStartOfSession": true,
				"verbosity":        "verbose",
				"scenario":         "SERP",
				"plugins":          plugins,
				"previousMessages": systemContext,
				"traceId":          strings.ReplaceAll(hex.NewUUID(), "-", ""),
				"conversationHistoryOptionsSets": []string{
					"autosave",
					"savemem",
					"uprofupd",
					"uprofgen",
				},
				"requestId": msgId,
				"message": map[string]any{
					"author":      "user",
					"inputMethod": "Keyboard",
					"text":        msg,
					"messageType": "Chat",
					"requestId":   msgId,
					"messageId":   msgId,
				},
				// "conversationSignature": chat.GetChatHub().GetConversationSignature(),
				"tone":           strings.ReplaceAll(tone, "-offline", ""),
				"spokenTextMode": "None",
				"participant": map[string]any{
					"id": chat.GetChatHub().GetClientId(),
				},
				"conversationId": chat.GetChatHub().GetConversationId(),
			},
		},
		"invocationId": "0",
		"target":       "chat",
		"type":         4,
	}

	return data
}

func (chat *Chat) wsHandler(data map[string]any) (*websocket.Conn, error) {
	wsConfig, _ := websocket.NewConfig(fmt.Sprintf(sydneyChatHubUrl, chat.SydneyBaseUrl, url.QueryEscape(chat.GetChatHub().GetEncryptedConversationSignature())), "https://"+chat.BingBaseUrl)
	wsConfig.Header.Add("Accept-Encoding", "gzip, deflate, br")
	wsConfig.Header.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	wsConfig.Header.Add("User-Agent", userAgent)
	wsConfig.Header.Add("Upgrade", "websocket")
	wsConfig.Header.Add("Connection", "Upgrade")
	wsConfig.Header.Add("Host", "sydney.bing.com")
	// wsConfig.Header.Add("Sec-Websocket-Extensions", "permessage-deflate; client_max_window_bits")
	wsConfig.Header.Add("Sec-WebSocket-Version", "13")

	ws, err := websocket.DialConfig(wsConfig)
	if err != nil {
		return nil, err
	}

	var buf = make([]byte, 1024)

	_, err = ws.Write([]byte("{\"protocol\":\"json\",\"version\":1}" + spilt))
	if err != nil {
		return nil, err
	}

	_, err = ws.Read(buf)
	if err != nil {
		return nil, err
	}

	_, err = ws.Write([]byte("{\"type\":6}" + spilt))
	if err != nil {
		return nil, err
	}

	req, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	_, err = ws.Write(append(req, []byte(spilt)...))
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (chat *Chat) Chat(prompt, msg string) (string, error) {
	optionsSets := chat.optionsSetsHandler()
	plugins := chat.pluginHandler(&optionsSets)
	systemContext := chat.systemContextHandler(prompt)
	data := chat.requestPayloadHandler(msg, optionsSets, plugins, systemContext)

	ws, err := chat.wsHandler(data)
	if err != nil {
		return "", err
	}
	defer ws.Close()

	buf := make([]byte, 1024)
	tmp := ""
	text := ""
	var resp ResponsePayload

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				return "", err
			}
		}
		if strings.Contains(string(buf[:n]), "\x1e") {
			t := strings.Split(string(buf[:n]), "\x1e")
			tmp += t[0]
			json.Unmarshal([]byte(tmp), &resp)
			if resp.Type == 2 {
				break
			} else if resp.Type == 1 {
				if len(resp.Arguments) > 0 {
					if len(resp.Arguments[0].Messages) > 0 {
						text = resp.Arguments[0].Messages[0].Text
						// fmt.Println(resp.Arguments[0].Messages[0].Text + "\n\n")
					}
				}
			}

			tmp = t[1]
		} else {
			tmp += string(buf[:n])
		}
	}

	return text, nil
}

func (chat *Chat) ChatStream(prompt, msg string, c chan string) (string, error) {
	optionsSets := chat.optionsSetsHandler()
	plugins := chat.pluginHandler(&optionsSets)
	systemContext := chat.systemContextHandler(prompt)
	data := chat.requestPayloadHandler(msg, optionsSets, plugins, systemContext)

	ws, err := chat.wsHandler(data)
	if err != nil {
		return "", err
	}
	defer ws.Close()

	buf := make([]byte, 1024)
	tmp := ""
	text := ""
	var resp ResponsePayload

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				return "", err
			}
		}
		if strings.Contains(string(buf[:n]), "\x1e") {
			t := strings.Split(string(buf[:n]), "\x1e")
			tmp += t[0]
			json.Unmarshal([]byte(tmp), &resp)
			if resp.Type == 2 {
				break
			} else if resp.Type == 1 {
				if len(resp.Arguments) > 0 {
					if len(resp.Arguments[0].Messages) > 0 {
						if len(resp.Arguments[0].Messages[0].Text) > len(text) {
							c <- strings.ReplaceAll(resp.Arguments[0].Messages[0].Text, text, "")
						}
						text = resp.Arguments[0].Messages[0].Text
						// fmt.Println(resp.Arguments[0].Messages[0].Text + "\n\n")
					}
				}
			}

			tmp = t[1]
		} else {
			tmp += string(buf[:n])
		}
	}

	c <- "EOF"

	return text, nil
}
