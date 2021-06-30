package post

type ListRequest struct {
	UserID       int
	LatestPostID int
}

type PostsUser struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
}

type Post struct {
	ID            int       `json:"id"`
	Text          string    `json:"text"`
	CreatedAt     int64     `json:"created_at"`
	CommentsCount uint      `json:"comments_count"`
	LikesCount    uint      `json:"likes_count"`
	RepostsCount  uint      `json:"reposts_count"`
	User          PostsUser `json:"user"`
}

type ListResponse struct {
	Total int    `json:"total"`
	Posts []Post `json:"posts"`
}

type CreateRequest struct {
	Text string `json:"text"`
}

type CreateResponse Post

type FeedRequest struct {
	LatestPostID int
}

type FeedResponse struct {
	Total int    `json:"total"`
	Posts []Post `json:"posts"`
}

type SearchRequest struct {
	Query   string // query is used for initial search. subsequent requests should use queryID
	QueryID int
	Offset  int // it's used for pagination
}

type SearchResponse struct {
	// total found records. there is some sane maximum value (1k for example)
	Total int `json:"total"`
	// search query id. created after first search request. subsequent requests for the same search query should use the same id
	QueryID int    `json:"query_id"`
	Posts   []Post `json:"posts"`
}
