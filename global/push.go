package global

import (
	"database/sql"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"firebase.google.com/go"
	"context"
	"server/message"
	"server/event"
)

func GetTokens(db *sql.DB, userIds ...int64) (map[int64][]string) {
	tokens := make(map[int64][]string, 0)

	for _, userId := range userIds {
		userTokens, err := GetTokensForUserId(db, userId)
		if err != nil {
			print(userId, err.Error())
			continue
		}

		tokens[userId] = userTokens
	}

	return tokens
}

func GetTokensForUserId(db *sql.DB, userId int64) ([]string, error) {
	query := "SELECT token FROM tokens WHERE userId = ?"

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	tokens := make([]string, 0)

	for rows.Next() {
		var token string
		rows.Scan(&token)

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func SendNotifications(notifications ...*messaging.Message) {
	ctx := context.Background()

	// Initialize the Firebase app
	opt := option.WithCredentialsFile("upfit-serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		println(err.Error())
		return
	}

	// Obtain a messaging client from the Firebase app
	client, err := app.Messaging(ctx)

	// Send notifications
	for _, notification := range notifications {
		_, err := client.Send(ctx, notification)
		if err != nil {
			print(err.Error())
			continue
		}
	}
}

func BaseNotification(token string) (*messaging.Message) {
	notification := &messaging.Message{
		Notification: &messaging.Notification{},
		Token: token,
	}

	customData := make(map[string]interface{})
	badgeIcon := 1

	notification.APNS = &messaging.APNSConfig{
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				CustomData: customData,
				Badge: &badgeIcon,
			},
		},
	}

	return notification
}

func DebugNotification(token string) (*messaging.Message) {
	notification := BaseNotification(token)

	notification.Notification.Body = "Notification de test envoyée."
	notification.APNS.Payload.Aps.CustomData["type"] = "debug"

	return notification
}

func MessageNotification(token string, message message.Info) (*messaging.Message) {
	notification := BaseNotification(token)

	notification.Notification.Title = message.Sender.FirstName + " " + message.Sender.LastName
	notification.Notification.Body = message.Content
	notification.APNS.Payload.Aps.CustomData["type"] = "message"

	return notification
}

func EventNotification(token string, event event.Info) (*messaging.Message) {
	notification := BaseNotification(token)

	notification.Notification.Title = event.Name
	notification.APNS.Payload.Aps.CustomData["type"] = "event"

	if event.Created == event.Updated {
		notification.Notification.Body = "Nouvel événement inséré dans votre calendrier."
	} else {
		notification.Notification.Body = "Événement modifié dans votre calendrier."
	}

	return notification
}

func CancelEventNotification(token string, event event.Info) (*messaging.Message) {
	notification := BaseNotification(token)

	notification.Notification.Title = event.Name
	notification.APNS.Payload.Aps.CustomData["type"] = "event"
	notification.Notification.Body = "Événement annulé."

	return notification
}