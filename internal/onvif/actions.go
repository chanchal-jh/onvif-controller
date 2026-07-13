package onvifservice

// Direction actions
const (
	ActionDirectionLeft  = "DIRECTION_LEFT"
	ActionDirectionRight = "DIRECTION_RIGHT"
	ActionDirectionUp    = "DIRECTION_UP"
	ActionDirectionDown  = "DIRECTION_DOWN"

	ActionZoomIn  = "ZOOM_IN"
	ActionZoomOut = "ZOOM_OUT"

	ActionSetPreset  = "SET_PRESET"
	ActionGotoPreset = "GOTO_PRESET"
)

type ActionRequest struct {
	IP          string   `json:"ip" binding:"required,ip"`
	Username    string   `json:"username" binding:"required"`
	Password    string   `json:"password" binding:"required"`
	Action      string   `json:"action" binding:"required,oneof=DIRECTION_LEFT DIRECTION_RIGHT DIRECTION_UP DIRECTION_DOWN ZOOM_IN ZOOM_OUT SET_PRESET GOTO_PRESET"`
	Speed       *float64 `json:"speed,omitempty"`
	PresetName  string   `json:"preset_name,omitempty"`
	PresetToken string   `json:"preset_token,omitempty"`
}
