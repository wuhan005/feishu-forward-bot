package handler

import (
	"context"
	"encoding/json"
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
	_, _, err := client.Message.ForwardMessage(context.Background(), &lark.ForwardMessageReq{
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
