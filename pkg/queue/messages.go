// worker package contains structs for worker messages (tasks). this structs are used for work queue.
package queue

const (
	JobEmailSend            = "job:email/send"
	JobUserCreate           = "job:user/create"
	JobAuthCreate           = "job:auth/create"
	JobPostFollow           = "job:post/follow"
	JobPostUnfollow         = "job:post/unfollow"
	JobAuthUpdateUserIDAuth = "job:auth/update_user_id"
)

// Email represents fields which email's worker fetches from the queue to handle.
// It sends emails specified in this struct.
type Email struct {
	RequestID      string
	RecipientEmail string
	Subject        string
	HTML           string
	Text           string
}

type UserCreate struct {
	RequestID string
	Login     string
	Email     string
}

type AuthCreate struct {
	RequestID string
	Login     string
	Email     string
	Password  string
}

type AuthUpdateUserID struct {
	RequestID string
	Login     string
	UserID    int
}

type PostFollow struct {
	RequestID    string
	UserID       int // current user id
	FollowUserID int // user id to follow
}

type PostUnfollow struct {
	RequestID      string
	UserID         int
	UnfollowUserID int
}
