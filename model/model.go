package model

type User struct {
	UserId int `json:"id"`
}

type SegName struct {
	SegName string `json:"seg_name"`
}

type Segment struct {
	SegID   int    `json:"seg_id"`
	SegName string `json:"seg_name"`
}

type GetUserSegments struct {
	SegID   int    `json:"seg_id"`
	SegName string `json:"seg_name"`
}


type UserSegments struct {
	UserID           int   `json:"user_id"`
	SegmentsToAdd    []string `json:"segments_to_add"`
	SegmentsToRemove []string `json:"segments_to_remove"`
}