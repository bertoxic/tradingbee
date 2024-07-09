package models

type JsonResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorJson  `json:"error,omitempty"`
}

type ErrorJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayLoad  `json:"log,omitempty"`
	Mail   MailPayload `jon:"mail,omitempty"`
}

type AuthPayload struct {
	UserDetails UserDetails `json:"user"`
}

type LogPayLoad struct {
	Name    string `json:"name"`
	Data    string `json:"data"`
	Message string `json:"message,omitempty"`
}

type MailPayload struct {
	FROM    string `json:"from"`
	TO      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
