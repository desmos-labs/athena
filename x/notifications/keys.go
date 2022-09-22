package notifications

const (
	NotificationTypeKey    = "type"
	TransactionHashKey     = "tx_hash"
	TransactionErrorKey    = "tx_error"
	SubspaceIDKey          = "subspace_id"
	RelationshipCreatorKey = "relationship_creator"
	PostIDKey              = "post_id"
	PostAuthorKey          = "post_author"
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
	TypeReply              = "comment"
	TypeRepost             = "repost"
	TypeQuote              = "quote"
	TypeMention            = "mention"
	TypeReaction           = "reaction"
)
