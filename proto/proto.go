package proto

type User struct {
	UserId        string `json:"userid"`
	DisplayName   string `json:"displayname"`
	PictureUrl    string `json:"pictureurl"`
	StatusMessage string `json:"statusmessage,omitempty"`
	Language      string `json:"language,omitempty"`
	Message       string `json:"message"`
}

type Output struct {
	UserId      string `json:"userid"`
	DisplayName string `json:"displayname"`
	PictureUrl  string `json:"pictureurl"`
}
