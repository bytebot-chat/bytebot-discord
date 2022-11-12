package model

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

// MessageSend is the struct that is used to pass messages from the Redis pubsub to the Discord Gateway (outbound messages)
// Because the discordgo.Session.ChannelMessageSend() method only accepts channel ID and content as a string, our struct limits iteslef to those two fields as well.
// Future work may expand this to include more fields or expand metadata to include more information that can be used to forumlate more complex responses.
type MessageSend struct {
	ChannelID string   `json:"channel_id,omitempty"` // ChannelID is the ID of the discord channel to send the message to
	Content   string   `json:"content,omitempty"`    // Content is the text body of the message to send
	Metadata  Metadata `json:"metadata,omitempty"`
}

// Deprecated in favor of newer methods that consume the entire model.Message struct
// MarshalReply converts the message to JSON and adds the metadata from the original message
// MarshalReply sends a response to the originating channel or direct message but does not do a "discord reply"
func (m *Message) MarshalReply(meta Metadata, dest string, s string) ([]byte, error) {
	reply := &MessageSend{
		Content:   s,
		ChannelID: dest,
		Metadata:  meta,
	}
	return json.Marshal(reply)
}

// RespondToChannelOrThread generates a MessageSend struct that can be used to respond to a channel or thread
// It optionally allows the message to reply or mention the user that sent the original message
func (m *Message) RespondToChannelOrThread(sourceApp, content string, shouldReply, shouldMention bool) *MessageSend {
	meta := Metadata{
		Source:      sourceApp,
		Dest:        m.Metadata.Source,
		ID:          uuid.NewV4(),
		Reply:       shouldReply,
		InReplyTo:   m.ID,
		MentionUser: shouldMention,
	}

	return &MessageSend{
		ChannelID: m.ChannelID,
		Content:   content,
		Metadata:  meta,
	}
}

// Unmarshal converts the JSON (in bytes) to a message
// This method is deprecated in favor of the UnmarshalJSON method and will be removed in a future release
// Correct behavior from this method is not guaranteed
// Example:
// 	msg := &model.MessageSend{}
// 	if err := msg.Unmarshal([]byte(`{"content":"hello world"}`)); err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(msg.Content)
func (m *MessageSend) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, m); err != nil {
		return err
	}
	return nil
}
