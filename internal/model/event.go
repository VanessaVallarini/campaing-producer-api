package model

const (
	EVENT_ACTION_CREATE            = "C"
	EVENT_ACTION_UPDATE_CLICK      = "UC"
	EVENT_ACTION_UPDATE_IMPRESSION = "UI"
	EVENT_ACTION_ACTIVATE          = "A"
	EVENT_ACTION_INACTIVATE        = "I"
)

type MessageBody struct {
	Message string `json:"Message"`
}

type CampaingCreatingEvent struct {
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
	Ip     string  `json:"ip"`
	Action string  `json:"action"`
}
