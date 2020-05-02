package resolvers

import (
    "context"
    "log"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
)

type conversationResolver struct{ *Resolver }

func (r *Resolver) Conversation() graphql.ConversationResolver {
    return &conversationResolver{r}
}
func (r *conversationResolver) Conversations(ctx context.Context, conversation *models.Conversation) ([]*models.Conversation, error) {
    return r.ConversationsRepo.GetConversationMessages(conversation.ID)
}

type subscriptionResolver struct{ *Resolver }

func (r *Resolver) Subscription() graphql.SubscriptionResolver {
    return &subscriptionResolver{r}
}
func (r *subscriptionResolver) MessageAdded(ctx context.Context) (<-chan *models.Conversation, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    err := r.createUser(currentUser.ID)
    if err != nil {
        return nil, err
    }
    event := make(chan *models.Conversation, 1)
    sub, err := r.nClient.Subscribe("conversation", func(t *models.Conversation) {
        participants, err := r.ConversationsRepo.GetConversationParticipants(t.ParentID)

        if err != nil {
            log.Fatalln("couldn't find participants", err)
        }
        for _, participant := range participants {
            if currentUser.ID == participant.UserID {
                event <- t
            }
        }
    })
    if err != nil {
        return nil, err
    }

    go func() {
        <-ctx.Done()
        sub.Unsubscribe()
    }()
    return event, nil
}
func (r *subscriptionResolver) UserJoined(ctx context.Context) (<-chan string, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    err := r.createUser(currentUser.ID)
    if err != nil {
        return nil, err
    }

    // Create new channel for request
    users := make(chan string, 1)
    r.mutex.Lock()
    r.userChannels[currentUser.ID] = users
    r.mutex.Unlock()

    // Delete channel when done
    go func() {
        <-ctx.Done()
        r.mutex.Lock()
        delete(r.userChannels, currentUser.ID)
        r.mutex.Unlock()
    }()

    return users, nil
}
func (r *queryResolver) Conversation(ctx context.Context, id string) (*models.Conversation, error) {
    return r.ConversationsRepo.GetByID(id)
}

func (r *queryResolver) Conversations(ctx context.Context) ([]*models.Conversation, error) {
    return nil, nil
}
