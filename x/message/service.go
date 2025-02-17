package message

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/totegamma/concurrent/x/core"
	"github.com/totegamma/concurrent/x/stream"
	"github.com/totegamma/concurrent/x/util"
)

// Service is the interface for message service
// Provides methods for message CRUD
type Service interface {
    Get(ctx context.Context, id string) (core.Message, error)
    PostMessage(ctx context.Context, objectStr string, signature string, streams []string) (core.Message, error)
    Delete(ctx context.Context, id string) (core.Message, error)
}

type service struct {
	rdb    *redis.Client
	repo   Repository
	stream stream.Service
}

// NewService creates a new message service
func NewService(rdb *redis.Client, repo Repository, stream stream.Service) Service {
	return &service{rdb, repo, stream}
}

// Get returns a message by ID
func (s *service) Get(ctx context.Context, id string) (core.Message, error) {
	ctx, span := tracer.Start(ctx, "ServiceGet")
	defer span.End()

	return s.repo.Get(ctx, id)
}

// PostMessage creates a new message
// It also posts the message to the streams
func (s *service) PostMessage(ctx context.Context, objectStr string, signature string, streams []string) (core.Message, error) {
	ctx, span := tracer.Start(ctx, "ServicePostMessage")
	defer span.End()

	var object SignedObject
	err := json.Unmarshal([]byte(objectStr), &object)
	if err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	if err := util.VerifySignature(objectStr, object.Signer, signature); err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	message := core.Message{
		Author:    object.Signer,
		Schema:    object.Schema,
		Payload:   objectStr,
		Signature: signature,
		Streams:   streams,
	}

	id, err := s.repo.Create(ctx, &message)
	if err != nil {
		span.RecordError(err)
		return message, err
	}

	for _, stream := range message.Streams {
		s.stream.Post(ctx, stream, id, "message", message.Author, "", "")
	}

	return message, nil
}

// Delete deletes a message by ID
// It also emits a delete event to the sockets
func (s *service) Delete(ctx context.Context, id string) (core.Message, error) {
	ctx, span := tracer.Start(ctx, "ServiceDelete")
	defer span.End()

	deleted, err := s.repo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		return core.Message{}, err
	}

	for _, deststream := range deleted.Streams {
		jsonstr, _ := json.Marshal(stream.Event{
			Stream: deststream,
			Type:   "message",
			Action: "delete",
			Body: stream.Element{
				ID: deleted.ID,
			},
		})
		err := s.rdb.Publish(context.Background(), deststream, jsonstr).Err()
		if err != nil {
			span.RecordError(err)
			return deleted, err
		}
	}

	return deleted, err
}
