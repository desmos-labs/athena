package types

const (
	NotificationTypeKey    = "type"
	TransactionHashKey     = "tx_hash"
	TransactionErrorKey    = "tx_error"
	RecipientKey           = "recipient"
	SubspaceIDKey          = "subspace_id"
	RelationshipCreatorKey = "relationship_creator"
	PostIDKey              = "post_id"
	PostAuthorKey          = "post_author"
	RepostIDKey            = "repost_id"
	RepostAuthorKey        = "repost_author"
	CommentIDKey           = "comment_id"
	CommentAuthorKey       = "comment_author"
	ReplyIDKey             = "reply_id"
	ReplyAuthorKey         = "reply_author"
	QuoteIDKey             = "quote_id"
	QuoteAuthorKey         = "quote_author"
	ReactionIDKey          = "reaction_id"
	ReactionAuthorKey      = "reaction_author"

	// Actions

	ClickActionKey   = "click_action"
	ClickActionValue = "open"

	NotificationActionKey = "action"
	ActionOpenPost        = "open_post"
	ActionOpenProfile     = "open_profile"

	// Notification types

	TypeTransactionSuccess = "transaction_success"
	TypeTransactionFailed  = "transaction_fail"
	TypeFollow             = "follow"
	TypeReply              = "reply"
	TypeComment            = "comment"
	TypeRepost             = "repost"
	TypeQuote              = "quote"
	TypeMention            = "mention"
	TypeReaction           = "reaction"
)
