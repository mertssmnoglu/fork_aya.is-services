package arcade

import "time"

type ExecuteToolResponse struct {
	ID            string             `json:"id"`
	ExecutionID   string             `json:"execution_id"`
	ExecutionType string             `json:"execution_type"`
	FinishedAt    time.Time          `json:"finished_at"`
	Duration      float64            `json:"duration"`
	Status        string             `json:"status"`
	Output        ToolResponseOutput `json:"output"`
	Success       bool               `json:"success"`
	RunAt         *time.Time         `json:"run_at,omitempty"`
}

type ToolResponseOutput struct {
	Value         ArcadeOutputValue `json:"value"` // NOTE: Spec defines 'value' as any (`{}`), this struct is specific.
	Error         *ToolError        `json:"error,omitempty"`
	Logs          []*ToolLog        `json:"logs,omitempty"`
	Authorization *AuthResponse     `json:"authorization,omitempty"`
}

// ArcadeOutputValue contains the main data, includes, and metadata.
// NOTE: This corresponds to the 'value' field in ToolResponseOutput.
// The spec defines 'value' as any (`{}`), meaning its structure depends on the tool executed.
// This specific struct (ArcadeOutputValue) is likely tailored for a specific tool (e.g., Twitter).
type ArcadeOutputValue struct {
	Data     []ArcadeTweetData `json:"data"`
	Includes ArcadeIncludes    `json:"includes"`
	Meta     ArcadeMeta        `json:"meta"`
}

// ArcadeTweetData represents a single tweet's data.
type ArcadeTweetData struct {
	AuthorID            string             `json:"author_id"`
	AuthorName          string             `json:"author_name,omitempty"`     // Optional, seems derived
	AuthorUsername      string             `json:"author_username,omitempty"` // Optional, seems derived
	EditHistoryTweetIDs []string           `json:"edit_history_tweet_ids"`
	ID                  string             `json:"id"`
	Text                string             `json:"text"`
	TweetURL            string             `json:"tweet_url"`
	Attachments         *ArcadeAttachments `json:"attachments,omitempty"` // Pointer for optional field
}

// ArcadeAttachments holds media keys associated with a tweet.
type ArcadeAttachments struct {
	MediaKeys []string `json:"media_keys"`
}

// ArcadeIncludes contains additional data like media and users.
type ArcadeIncludes struct {
	Media []ArcadeMedia `json:"media"`
	Users []ArcadeUser  `json:"users"`
}

// ArcadeMedia represents a media item (photo, video, etc.).
type ArcadeMedia struct {
	Height   int    `json:"height"`
	MediaKey string `json:"media_key"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Width    int    `json:"width"`
}

// ArcadeUser represents a user mentioned or included in the response.
type ArcadeUser struct {
	Entities *ArcadeUserEntities `json:"entities,omitempty"` // Pointer for optional field
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Username string              `json:"username"`
}

// ArcadeUserEntities holds entity information for a user, like URLs.
type ArcadeUserEntities struct {
	URL ArcadeUserURL `json:"url"`
}

// ArcadeUserURL contains URLs associated with a user profile.
type ArcadeUserURL struct {
	URLs []ArcadeUserURLEntry `json:"urls"`
}

// ArcadeUserURLEntry represents a single URL entity.
type ArcadeUserURLEntry struct {
	DisplayURL  string `json:"display_url"`
	End         int    `json:"end"`
	ExpandedURL string `json:"expanded_url"`
	Start       int    `json:"start"`
	URL         string `json:"url"`
}

// ArcadeMeta contains pagination and result count information.
type ArcadeMeta struct {
	NewestID    string `json:"newest_id"`
	NextToken   string `json:"next_token"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
}

// Based on spec: tool.Error.
type ToolError struct {
	AdditionalPromptContent string `json:"additional_prompt_content,omitempty"`
	CanRetry                bool   `json:"can_retry,omitempty"`
	DeveloperMessage        string `json:"developer_message,omitempty"`
	Message                 string `json:"message"` // Required
	RetryAfterMs            int    `json:"retry_after_ms,omitempty"`
}

// Based on spec: tool.Log.
type ToolLog struct {
	Level   string `json:"level"`   // Required
	Message string `json:"message"` // Required
	Subtype string `json:"subtype,omitempty"`
}

// Based on spec: auth.AuthorizationContext.
type AuthorizationContext struct {
	Token    string         `json:"token,omitempty"`
	UserInfo map[string]any `json:"user_info,omitempty"`
}

// Based on spec: auth.AuthorizationStatus.
type AuthorizationStatus string

const (
	StatusPending   AuthorizationStatus = "pending"
	StatusCompleted AuthorizationStatus = "completed"
	StatusFailed    AuthorizationStatus = "failed"
)

// Based on spec: auth.AuthorizationResponse.
type AuthResponse struct {
	Context    *AuthorizationContext `json:"context,omitempty"`
	ID         string                `json:"id,omitempty"`
	ProviderID string                `json:"provider_id,omitempty"`
	Scopes     []string              `json:"scopes,omitempty"`
	Status     AuthorizationStatus   `json:"status,omitempty"`
	URL        string                `json:"url,omitempty"`
	UserID     string                `json:"user_id,omitempty"`
}

type ExecuteToolRequest struct {
	Input       ExecuteToolInput `json:"input"`
	ToolName    string           `json:"tool_name"`
	UserID      string           `json:"user_id"`
	ToolVersion *string          `json:"tool_version,omitempty"`
}

// ExecuteToolInput represents the input part of the request payload.
// NOTE: The spec defines input as tool.RawInputs (any JSON object `{}`).
// This struct assumes a specific input structure required by a particular tool.
type ExecuteToolInput struct {
	Username   string `json:"username"`
	MaxResults string `json:"max_results"`
}
