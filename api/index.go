package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chyroc/lark"
)

type requestBody struct {
	MessageID string `json:"message_id"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var requestBody requestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	appID := os.Getenv("LARK_APP_ID")
	appSecret := os.Getenv("LARK_APP_SECRET")
	receiveIDType := os.Getenv("LARK_RECEIVE_ID_TYPE")
	receiveID := os.Getenv("LARK_RECEIVE_ID")

	client := lark.New(lark.WithAppCredential(appID, appSecret))
	messageResp, _, err := client.Message.GetMessage(r.Context(), &lark.GetMessageReq{
		MessageID: requestBody.MessageID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sender := messageResp.Items[0].Sender
	senderResp, _, err := client.Contact.GetUser(r.Context(), &lark.GetUserReq{
		UserIDType: &sender.IDType,
		UserID:     sender.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	senderName := senderResp.User.Name
	_, _, err = client.Message.SendRawMessage(r.Context(), &lark.SendRawMessageReq{
		ReceiveIDType: lark.IDType(receiveIDType),
		ReceiveID:     receiveID,
		MsgType:       lark.MsgTypeText,
		Content:       fmt.Sprintf(`{"text": "收到来自 %s 的消息"}`, senderName),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _, err = client.Message.ForwardMessage(r.Context(), &lark.ForwardMessageReq{
		MessageID:     requestBody.MessageID,
		ReceiveIDType: lark.IDType(receiveIDType),
		ReceiveID:     receiveID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}
