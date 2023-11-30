package models

import "time"

////=======DashBoard
type Users struct {
	FirstName   string    `bson:"Имя"`
	LastName    string    `bson:"Фамилия"`
	MiddleName  string    `bson:"Отчество"`
	DateOfBirth time.Time `bson:"Дата рождения"`
	Email       string    `bson:"E-mail"`
	Unsubscribe bool      `bson:"unsubscribe"`
}

type Logs struct {
	FirstName   string    `bson:"Имя"`
	LastName    string    `bson:"Фамилия"`
	MiddleName  string    `bson:"Отчество"`
	DateOfBirth time.Time `bson:"Дата рождения"`
	Email       string    `bson:"E-mail"`
	DateCreate  time.Time `bson:"dateCreate"`
}
type Dashboard_Params struct {
	UUID         string `json:"uuid"`
	SendTo       string `json:"sendTo"`
	SendTemplate string `json:"sendTemplate"`
	Page         int    `json:"page"`
	Templates    bool   `json:"templates"`
	Seach        string `json:"seach"`
}

type DashboardGetResponse struct {
	Templates     []Templates `json:"templates"`
	UsersCount    int64       `json:"usersCount"`
	LogsCount     int64       `json:"logsCount"`
	CountBirtdays int         `json:"countBirtdays"`
	CountLogs     int64       `json:"todaySent"`
	SendEmail     string      `json:"sendEmailresult"`
	SendAutoAt    int         `json:"sendAutoAt"`
}
type DashboardPostResponse struct {
	Err               string `json:"err"`
	DocumentsInserted int64  `json:"documentsInserted"`
	DocumentsModified int64  `json:"documentsModified"`
}
type UsersUpload struct {
	Last_name   string `json:"Фамилия" bson:"Фамилия"`
	First_name  string `json:"Имя" bson:"Имя"`
	Middle_name string `json:"Отчество" bson:"Отчество"`
	Date_birth  string `json:"Дата рождения" bson:"Дата рождения"`
	Email       string `json:"E-mail" bson:"E-mail"`
}

////=======Settings
type SettingsUpload struct {
	Records    []Events    `json:"records"`
	Templates  []Templates `json:"templates"`
	EmailLogin string      `json:"emailLogin" bson:"emailLogin"`
	EmailPass  string      `json:"emailPass" bson:"emailPass"`
	Smtp       string      `json:"smtp" bson:"smtp"`
	Port       string      `json:"port" bson:"port"`
	UUID       string      `json:"uuid"`
}

///========DataBase
type GetDataBaseResponse struct {
	Records    []Users `json:"records"`
	TotalFound int64   `json:"totalFound"`
	UsersCount int64   `json:"usersCount"`
	PageNumber int     `json:"pageNumber"`
}

//======History
type GetHistoryResponse struct {
	Records            []Logs `json:"records"`
	LogsCount          int64  `json:"logsCount"`
	TotalFound         int64  `json:"totalFound"`
	TodayLogsCount     int64  `json:"todayLogsCount"`
	TommorowLogsCount  int64  `json:"tommorowLogsCount"`
	YesterdayLogsCount int64  `json:"yesterdayLogsCount"`
}

//Auth
type Auth struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	UUID     string `json:"uuid"`
}

//Unsubscribe
type Unsubscribe struct {
	Email string `json:"email"`
}

//Events
type Events struct {
	Name         string    `json:"name" bson:"name"`
	IsDaily      bool      `json:"isDaily" bson:"isDaily"`
	Type         string    `bson:"type"`
	IsSent       bool      `json:"isSent" bson:"isSent"`
	Subject      string    `json:"subject" bson:"subject"`
	Day          int64     `json:"day" bson:"day"`
	Month        int64     `json:"month" bson:"month"`
	TemplateName string    `json:"templateName" bson:"templateName"`
	MustSend     time.Time `json:"mustSend" bson:"mustSend"`
	SendAt       int64     `json:"sendAt" bson:"sendAt"`
	From         string    `json:"from" bson:"from"`
	Active       bool      `json:"active"`
	UUID         string    `json:"uuid"`
}

//Templates
type Templates struct {
	Name      string `bson:"name" json:"name"`
	IndexHTML string `bson:"indexHTML" json:"indexHtml"`
}
type TemplatesList struct {
	Name string `bson:"name"`
}
