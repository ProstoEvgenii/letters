package models

import "time"

////=======DashBoard
type Users struct {
	FirstName   string    `bson:"Имя"`
	LastName    string    `bson:"Фамилия"`
	MiddleName  string    `bson:"Отчество"`
	DateOfBirth time.Time `bson:"Дата рождения"`
	Email       string    `bson:"E-mail"`
}
type Templates struct {
	Name      string `bson:"name"`
	IndexHTML string `bson:"indexHTML"`
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
	SendTo     string `json:"sendTo"`
	SendAutoAt int    `json:"sendAutoAt"`
}

type DashboardGetResponse struct {
	UsersCount    int64  `json:"usersCount"`
	LogsCount     int64  `json:"logsCount"`
	CountBirtdays int    `json:"countBirtdays"`
	CountLogs     int64  `json:"todaySent"`
	SendEmail     string `json:"sendEmailresult"`
	SendAutoAt    int    `json:"sendAutoAt"`
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

type IsSent struct {
	Name       string    `bson:"name"`
	dateCreate time.Time `bson:"date"`
	IsSent     bool      `bson:"isSent"`
}

////=======Settings

type GetSettingsResponse struct {
	EventsRecords []Events       `json:"events_records"`
	Settings      SettingsUpload `json:"settings"`
}

type SettingsPost struct {
}
type SettingsUpload struct {
	Template   string `json:"template" bson:"template"`
	EmailLogin string `json:"emailLogin" bson:"emailLogin"`
	EmailPass  string `json:"emailPass" bson:"emailPass"`
	Smtp       string `json:"smtp" bson:"smtp"`
	Port       string `json:"port" bson:"port"`
	SendAutoAt int    `bson:"sendAutoAt"`
}
type Events struct {
	Name         string `json:"name" bson:"name"`
	TemplateName string `json:"template" bson:"templateName"`
	SendAt       int    `json:"sendAt" bson:"sendAt"`
	IsDaily      string `json:"isDaily" bson:"isDaily"`
	Date         int    `json:"date" bson:"date"`
	Active       bool   `json:"active" bson:"active"`
}

type EventUpload struct {
	Name         string `json:"name" bson:"name"`
	TemplateName string `json:"template" bson:"templateName"`
	SendAt       string `json:"sendAt" bson:"sendAt"`
	IsDaily      string `json:"isDaily" bson:"isDaily"`
	Date         string `json:"date"`
	Month        string `json:"month"`
	// Active       string `json:"active" bson:"active"`
}

///========DataBase

type UserRecord struct {
	User UsersUpload `json:"user"`
}
type GetDataBaseResponse struct {
	Records    []Users `json:"records"`
	UsersCount int64   `json:"usersCount"`
}

//======History

type GetHistoryResponse struct {
	Records            []Logs `json:"records"`
	LogsCount          int64  `json:"logsCount"`
	TodayLogsCount     int    `json:"todayLgsCount"`
	TommorowLogsCount  int    `json:"tommorowLogsCount"`
	YesterdayLogsCount int    `json:"yesterdayLogsCount"`
}

// type GetDataBaseResponse struct {
// 	Records    []byte
// 	UsersCount int64
// }

////===========================
// type Data struct {
// 	Last_name   string `json:"Фамилия" bson:"Фамилия"`
// 	First_name  string `json:"Имя" bson:"Имя"`
// 	Middle_name string `json:"Отчество" bson:"Отчество"`
// 	Date_birth  string `json:"Дата рождения" bson:"Дата рождения"`
// 	Email       string `json:"E-mail" bson:"E-mail"`
// }
