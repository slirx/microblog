package user

type UpdateRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type MeResponse struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	PhotoURL  string `json:"photo_url"`
	Bio       string `json:"bio"`
	Following uint   `json:"following"`
	Followers uint   `json:"followers"`
}

type GetResponse struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	Name       string `json:"name"`
	PhotoURL   string `json:"photo_url"`
	Bio        string `json:"bio"`
	Following  uint   `json:"following"`
	Followers  uint   `json:"followers"`
	IsFollowed bool   `json:"is_followed"` // whether this user is followed by the current one
}

type FollowRequest struct {
	UserID int `json:"user_id"`
}

type UnfollowRequest struct {
	UserID int `json:"user_id"`
}

type FollowersRequest struct {
	Login            string
	LatestFollowerID int
}

type FollowersResponse struct {
	Total     int        `json:"total"`
	Followers []Follower `json:"followers"`
}

type Follower struct {
	FollowerID int    `json:"follower_id,omitempty"`
	UserID     int    `json:"user_id,omitempty"`
	Login      string `json:"login,omitempty"`
	Name       string `json:"name,omitempty"`
	PhotoURL   string `json:"photo_url,omitempty"`
	Bio        string `json:"bio,omitempty"`
}

type FollowingRequest struct {
	Login            string
	LatestFollowerID int
}

type FollowingResponse struct {
	Total     int        `json:"total"`
	Following []Followed `json:"following"`
}

type Followed struct {
	FollowerID int    `json:"follower_id,omitempty"`
	UserID     int    `json:"user_id,omitempty"`
	Login      string `json:"login,omitempty"`
	Name       string `json:"name,omitempty"`
	PhotoURL   string `json:"photo_url,omitempty"`
	Bio        string `json:"bio,omitempty"`
}

// User is used for response in InternalUsers service.
type User struct {
	ID        int    `json:"id,omitempty"`
	Login     string `json:"login,omitempty"`
	Name      string `json:"name,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Following uint   `json:"following,omitempty"`
	Followers uint   `json:"followers,omitempty"`
}
