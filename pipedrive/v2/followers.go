package v2

import "time"

type Follower struct {
	UserID  UserID     `json:"user_id,omitempty"`
	AddTime *time.Time `json:"add_time,omitempty"`
}

type FollowerChangelog struct {
	Action         string     `json:"action,omitempty"`
	ActorUserID    UserID     `json:"actor_user_id,omitempty"`
	FollowerUserID UserID     `json:"follower_user_id,omitempty"`
	Time           *time.Time `json:"time,omitempty"`
}

type FollowerDeleteResult struct {
	UserID UserID `json:"user_id"`
}
