package duckchat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

type Agent struct {
	chat Chat
	vqd  string
}

func NewAgent(model string) Agent {
	a := Agent{
		chat: NewChat(model),
		vqd:  "",
	}

	a.GetVqd()

	return a
}

func (a *Agent) Send(message string) (string, error) {
	a.chat.AddMessage(NewMessage("user", message))
	return a.MakeRequest()
}

func (a *Agent) GetVqd() error {
	r, err := http.NewRequest("GET", ENDPOINT_STATUS, nil)
	if err != nil {
		return err
	}

	r.Header.Add("x-vqd-accept", "1")
	r.Header.Add("User-Agent", USER_AGENT)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	a.vqd = resp.Header.Get("x-vqd-4")
	return nil
}

func (a *Agent) MakeRequest() (string, error) {
    if a.vqd == "" {
        err := a.GetVqd()
        if err != nil {
            return "", err
        }
    }

    buf, err := a.chat.Json()
    if err != nil {
        return "", err
    }

    r, err := http.NewRequest("POST", ENDPOINT_CHAT, bytes.NewBuffer(buf))
    if err != nil {
        return "", err
    }

    r.Header.Add("Content-Type", "application/json")
    r.Header.Add("x-vqd-4", a.vqd)
    r.Header.Add("User-Agent", USER_AGENT)

    client := &http.Client{}
    resp, err := client.Do(r)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    a.vqd = resp.Header.Get("x-vqd-4")

    ai_message := ""
    reader := bufio.NewReader(resp.Body)

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", fmt.Errorf("error reading response: %v", err)
        }

        line, _ = strings.CutSuffix(line, "\n")
        line, _ = strings.CutPrefix(line, "data: ")
        if line == "" {
            continue
        }

        var res map[string]string
        json.Unmarshal([]byte(line), &res)

        if text, ok := res["message"]; ok {
            ai_message += text
        }
    }

    a.chat.AddMessage(NewMessage("assistant", ai_message))

    return ai_message, nil
}

