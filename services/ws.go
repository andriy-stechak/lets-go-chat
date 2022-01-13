package services

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/andriystech/lgc/db/repositories"
	"github.com/andriystech/lgc/facilities/ws"
	"github.com/andriystech/lgc/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketService interface {
	NewConnection(http.ResponseWriter, *http.Request, *models.User) error
	GetActiveConnectionsCount(context.Context) (int, error)
	GetActiveUsers(context.Context) ([]string, error)
	SendMessageToAllConnections(context.Context, string, *models.User) error
	LoadUserMessages(context.Context, *models.User, ws.ConnHelper) error
	SaveUnreadMessages(context.Context, *models.User, string) error
}

type webSocketService struct {
	connections repositories.ConnectionsRepository
	messages    repositories.MessagesRepository
	upgrader    ws.UpgraderHelper
	users       repositories.UsersRepository
}

func NewWebSocketService(
	cr repositories.ConnectionsRepository,
	mr repositories.MessagesRepository,
	ur repositories.UsersRepository,
	wu ws.UpgraderHelper,
) WebSocketService {
	return &webSocketService{
		connections: cr,
		messages:    mr,
		upgrader:    wu,
		users:       ur,
	}
}

func (svc *webSocketService) NewConnection(w http.ResponseWriter, r *http.Request, user *models.User) error {
	c, err := svc.upgrader.Upgrade(w, r)
	if err != nil {
		log.Printf("Unable to establish web socket connection. Reason: %s", err.Error())
		return err
	}

	id := uuid.NewString()
	if svc.connections.AddConnection(r.Context(), id, c, user); err != nil {
		return err
	}

	defer func() {
		c.Close()
		if err = svc.connections.DeleteConnection(r.Context(), id); err != nil {
			log.Printf("Unable to delete connection. Reason: %s", err.Error())
		}
	}()

	if err = svc.LoadUserMessages(r.Context(), user, c); err != nil {
		log.Printf("Unable to read messages. Reason: %s", err.Error())
		return err
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("web socket read error:", err)
			break
		}
		if err = svc.SendMessageToAllConnections(r.Context(), string(message), user); err != nil {
			log.Println("web socket write error:", err)
			break
		}
		if err = svc.SaveUnreadMessages(r.Context(), user, string(message)); err != nil {
			log.Println("save unread messages error:", err)
			break
		}
	}

	return nil
}

func (svc *webSocketService) LoadUserMessages(ctx context.Context, usr *models.User, conn ws.ConnHelper) error {
	messages, err := svc.messages.FindUserMessages(ctx, usr.Id)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if err = conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			return err
		}
	}

	return nil
}

func (svc *webSocketService) SaveUnreadMessages(ctx context.Context, sender *models.User, message string) error {
	cs, err := svc.connections.GetAllConnections(ctx)
	if err != nil {
		return err
	}

	var activeUsrIds []string
	for usrId := range cs {
		activeUsrIds = append(activeUsrIds, usrId)
	}

	notActiveUsers, err := svc.users.FindUsersNotInIdList(ctx, activeUsrIds)
	if err != nil {
		return err
	}

	for _, usr := range notActiveUsers {
		msg := models.NewMessage(
			uuid.NewString(),
			sender.Id,
			sender.UserName,
			usr.Id,
			message,
		)
		if _, err = svc.messages.SaveMessage(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}

func (svc *webSocketService) sendMessage(
	ctx context.Context,
	wg *sync.WaitGroup,
	ch chan error,
	conn ws.ConnHelper,
	msg *models.Message,
) {
	defer wg.Done()

	if _, err := svc.messages.SaveMessage(ctx, msg); err != nil {
		ch <- err
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
		ch <- err
		return
	}
}

func (svc *webSocketService) SendMessageToAllConnections(
	ctx context.Context,
	payload string,
	sender *models.User,
) error {
	var wg sync.WaitGroup

	cs, err := svc.connections.GetAllConnections(ctx)
	if err != nil {
		return err
	}

	ch := make(chan error, len(cs))

	for rId, conn := range cs {
		wg.Add(1)
		if sender.Id == rId {
			wg.Done()
			continue
		}

		msg := models.NewMessage(
			uuid.NewString(),
			sender.Id,
			sender.UserName,
			rId,
			payload,
		)
		go svc.sendMessage(ctx, &wg, ch, conn, msg)
	}

	wg.Wait()
	close(ch)

	var errs []error
	for err := range ch {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		log.Printf("Warn: Unable to send message to some users. Reasons: %v", errs)
	}

	return nil
}

func (svc *webSocketService) GetActiveConnectionsCount(ctx context.Context) (int, error) {
	return svc.connections.CountConnections(ctx)
}

func (svc *webSocketService) GetActiveUsers(ctx context.Context) ([]string, error) {
	return svc.connections.ConnectedClients(ctx)
}
