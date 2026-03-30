package callsignal

import (
	"encoding/json"
	"ws/src/chat"
)

func SignalMessages(c *chat.Client) {
	defer func() {
		// Khi client đóng, xoá userid được khai báo ra khỏi map
		// và đóng kết n
		delete(signalingClients, c.UserID)
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse JSON từ client -> SignalPayload
		var signal SignalPayload
		if err := json.Unmarshal(msg, &signal); err != nil {
			continue
		}

		targetClient, ok := signalingClients[signal.ToUserID]

		if ok {
			targetMsg := map[string]any{
				"from": c.UserID,
				"type": signal.Type,
				"data": signal.Data,
			}

			out, _ := json.Marshal(targetMsg)
			targetClient.Send <- out
		} else {
			out, _ := json.Marshal(map[string]any{
				"to_user_id is not connected": signal.ToUserID,
			})
			c.Send <- out
		}
	}
}
