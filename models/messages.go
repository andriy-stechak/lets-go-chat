package models

import "time"

type Message struct {
	Id          string `bson:"_id"`
	RecipientId string `bson:"recipientId"`
	SenderId    string `bson:"senderId"`
	SenderName  string `bson:"senderName"`
	Payload     string `bson:"payload"`
	Time        int64  `bson:"time"`
}

func NewMessage(id, sId, sName, rId, payload string) *Message {
	return &Message{
		Id:          id,
		RecipientId: rId,
		SenderId:    sId,
		SenderName:  sName,
		Payload:     payload,
		Time:        time.Now().Unix(),
	}
}
