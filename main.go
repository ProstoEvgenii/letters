package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type Data struct {
	Last_name   string `json:"Фамилия"`
	First_name  string `json:"Имя"`
	Middle_name string `json:"Отчество"`
	Date_birth  string `json:"Дата рождения"`
	Email       string `json:"E-mail"`
}

func main() {
	files_name := []string{"index.html", ".env", "users.json"}
	CheckFilesAndConnectToEmail(files_name)

	file, err := os.Open("users.json")
	if err != nil {
		fmt.Println("Не найден файл БД users.json")
		time.Sleep(10 * time.Second)
		log.Fatal()
	}
	defer file.Close()

	var records []Data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&records)
	if err != nil {
		if err.Error() == "EOF" && len(records) == 0 {
			fmt.Println("Файл users.json пуст")
			time.Sleep(10 * time.Second)
			log.Fatal()
		} else {
			fmt.Println("Убедитесь, что вы используете верную базу данных", err)
			time.Sleep(10 * time.Second)
			log.Fatal()
		}
	}

	todayMonthDate := time.Now().Format("01/02")
	foundBirthday := false
	for _, item := range records {
		if item.Date_birth[:5] == todayMonthDate {
			checkAndLog(item)
			foundBirthday = true
		}
	}

	if !foundBirthday {
		log.Println("Сегодня нет дней рождений среди пользователей.")
		time.Sleep(10 * time.Second)
		log.Fatal()
	}
}
func CheckFilesAndConnectToEmail(files_name []string) {
	fmt.Println("Проверяю файлы в папке...")
	for _, item := range files_name {
		if _, err := os.Stat(item); os.IsNotExist(err) {
			fmt.Printf("Файл %s не найден в корне проекта.\n", item)
			time.Sleep(10 * time.Second)
			log.Fatal()
		}
	}
	log.Println("Все файлы на присутствуют.")

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		time.Sleep(10 * time.Second)
		log.Fatal()
	}

	// d := gomail.NewDialer("smtp.yandex.ru", 465, "support@crypto-emergency.com", os.Getenv("EMAIL_PASS"))
	// if err := d.DialAndSend(); err != nil {
	// 	log.Printf("Не удалось отправить установить соединение с почтовым ящиком. Убедитесь ,что E-mail и пароль в файле .env указаны верно \n%v", err)
	// 	time.Sleep(10 * time.Second)
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Соединение с почтовым ящиком установлено.")
	// }
}
func checkAndLog(item Data) {
	log.Println("=1854d3=", item)
	log_name := time.Now().Format("01.02.2006")

	existingLogs, err := os.Open("./logs/" + log_name + ".json")
	if os.IsNotExist(err) { //Если файла log не существует создаю и записываю в него item
		create_log(item)
		SendEmailReg(item)
		return
	}
	defer existingLogs.Close()

	var logs []Data
	decoder := json.NewDecoder(existingLogs)
	err = decoder.Decode(&logs)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	itemAlreadyExists := false
	for _, log := range logs {
		if log == item {
			itemAlreadyExists = true
			break
		}
	}

	if !itemAlreadyExists { //Если нет в логах - поздравить и записать в логи
		SendEmailReg(item)
		logs = append(logs, item)
		logJson, err := json.Marshal(logs)
		overwriteLogs, err := os.Create("./logs/" + log_name + ".json")
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		_, err = overwriteLogs.Write(logJson)
		if err != nil {
			log.Println("=Ошибка записи в json=", err)
		}
		defer overwriteLogs.Close()
	}

}
func create_log(item Data) {
	log_name := time.Now().Format("01.02.2006")

	newLog := []Data{item} //Форматирую item в Json
	logJson, err := json.Marshal(newLog)
	if err != nil {

		// log.Fatal("=Ошибка форматирования лога в json=", err)
	}

	if err := os.MkdirAll("./logs", os.ModePerm); err != nil {
		log.Fatal("Ошибка при создании директории logs:", err)
	}
	newLogs, err := os.Create("./logs/" + log_name + ".json")
	if err != nil {
		fmt.Println("Unable to create file:", err)
	}
	_, err = newLogs.Write(logJson)
	if err != nil {
		log.Println("=Ошибка записи в json=", err)
	}
	defer newLogs.Close()
}

var confirmEmail = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD XHTML 1.0 Transitional //EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
<!--[if gte mso 9]>
<xml>
  <o:OfficeDocumentSettings>
    <o:AllowPNG/>
    <o:PixelsPerInch>96</o:PixelsPerInch>
  </o:OfficeDocumentSettings>
</xml>
<![endif]-->
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="x-apple-disable-message-reformatting">
  <!--[if !mso]><!--><meta http-equiv="X-UA-Compatible" content="IE=edge"><!--<![endif]-->
  <title></title>
  
    <style type="text/css">
      @media only screen and (min-width: 620px) {
  .u-row {
    width: 600px !important;
  }
  .u-row .u-col {
    vertical-align: top;
  }

  .u-row .u-col-100 {
    width: 600px !important;
  }

}

@media (max-width: 620px) {
  .u-row-container {
    max-width: 100% !important;
    padding-left: 0px !important;
    padding-right: 0px !important;
  }
  .u-row .u-col {
    min-width: 320px !important;
    max-width: 100% !important;
    display: block !important;
  }
  .u-row {
    width: 100% !important;
  }
  .u-col {
    width: 100% !important;
  }
  .u-col > div {
    margin: 0 auto;
  }
}
body {
  margin: 0;
  padding: 0;
}

table,
tr,
td {
  vertical-align: top;
  border-collapse: collapse;
}

p {
  margin: 0;
}

.ie-container table,
.mso-container table {
  table-layout: fixed;
}

* {
  line-height: inherit;
}

a[x-apple-data-detectors='true'] {
  color: inherit !important;
  text-decoration: none !important;
}

table, td { color: #000000; } </style>
  
  

<!--[if !mso]><!--><link href="https://fonts.googleapis.com/css?family=Raleway:400,700&display=swap" rel="stylesheet" type="text/css"><link href="https://fonts.googleapis.com/css?family=Raleway:400,700&display=swap" rel="stylesheet" type="text/css"><!--<![endif]-->

</head>

<body class="clean-body u_body" style="margin: 0;padding: 0;-webkit-text-size-adjust: 100%;background-color: #ffffff;color: #000000">
  <!--[if IE]><div class="ie-container"><![endif]-->
  <!--[if mso]><div class="mso-container"><![endif]-->
  <table style="border-collapse: collapse;table-layout: fixed;border-spacing: 0;mso-table-lspace: 0pt;mso-table-rspace: 0pt;vertical-align: top;min-width: 320px;Margin: 0 auto;background-color: #ffffff;width:100%" cellpadding="0" cellspacing="0">
  <tbody>
  <tr style="vertical-align: top">
    <td style="word-break: break-word;border-collapse: collapse !important;vertical-align: top">
    <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td align="center" style="background-color: #ffffff;"><![endif]-->
    
  
  
<div class="u-row-container" style="padding: 0px 10px;background-color: rgba(255,255,255,0)">
  <div class="u-row" style="margin: 0 auto;min-width: 320px;max-width: 600px;overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;background-color: #b5e0ff;">
    <div style="border-collapse: collapse;display: table;width: 100%;height: 100%;background-color: transparent;">
      <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding: 0px 10px;background-color: rgba(255,255,255,0);" align="center"><table cellpadding="0" cellspacing="0" border="0" style="width:600px;"><tr style="background-color: #b5e0ff;"><![endif]-->
      
<!--[if (mso)|(IE)]><td align="center" width="600" style="width: 600px;padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;" valign="top"><![endif]-->
<div class="u-col u-col-100" style="max-width: 320px;min-width: 600px;display: table-cell;vertical-align: top;">
  <div style="height: 100%;width: 100% !important;">
  <!--[if (!mso)&(!IE)]><!--><div style="box-sizing: border-box; height: 100%; padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;"><!--<![endif]-->
  
<table style="font-family:'Raleway',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
  <tbody>
    <tr>
      <td style="overflow-wrap:break-word;word-break:break-word;padding:16px 20px 8px;font-family:'Raleway',sans-serif;" align="left">
        
  <div style="font-size: 14px; color: #ffffff; line-height: 120%; text-align: center; word-wrap: break-word;">
    <p style="font-size: 14px; line-height: 120%;"><strong><span style="font-size: 48px; line-height: 57.6px; font-family: Raleway, sans-serif;">С днем роождения! ${first_name} ${last_name}</span></strong></p>
  </div>

      </td>
    </tr>
  </tbody>
</table>

<table style="font-family:'Raleway',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
  <tbody>
    <tr>
      <td style="overflow-wrap:break-word;word-break:break-word;padding:22px 0px 0px;font-family:'Raleway',sans-serif;" align="left">
        
<table width="100%" cellpadding="0" cellspacing="0" border="0">
  <tr>
    <td style="padding-right: 0px;padding-left: 0px;" align="center">
      
      <img align="center" border="0" src="images/image-1.png" alt="Image" title="Image" style="outline: none;text-decoration: none;-ms-interpolation-mode: bicubic;clear: both;display: inline-block !important;border: none;height: auto;float: none;width: 100%;max-width: 564px;" width="564"/>
      
    </td>
  </tr>
</table>

      </td>
    </tr>
  </tbody>
</table>

<table style="font-family:'Raleway',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
  <tbody>
    <tr>
      <td style="overflow-wrap:break-word;word-break:break-word;padding:15px 20px 14px;font-family:'Raleway',sans-serif;" align="left">
        
  <div style="font-size: 14px; color: #ffffff; line-height: 120%; text-align: center; word-wrap: break-word;">
    <p style="font-size: 14px; line-height: 120%;"><span style="font-size: 30px; line-height: 36px;">ляляля</span></p>
  </div>

      </td>
    </tr>
  </tbody>
</table>

<table style="font-family:'Raleway',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
  <tbody>
    <tr>
      <td style="overflow-wrap:break-word;word-break:break-word;padding:10px 20px;font-family:'Raleway',sans-serif;" align="left">
        
  <div style="font-size: 14px; color: #ffffff; line-height: 130%; text-align: center; word-wrap: break-word;">
    <p style="font-size: 14px; line-height: 130%;"><span style="font-size: 16px; line-height: 20.8px;">Your birthday comes around only once a year,</span></p>
<p style="font-size: 14px; line-height: 130%;"><span style="font-size: 16px; line-height: 20.8px;">so let’s give it the attention it deserves.</span></p>
  </div>

      </td>
    </tr>
  </tbody>
</table>

<table style="font-family:'Raleway',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
  <tbody>
    <tr>
      <td style="overflow-wrap:break-word;word-break:break-word;padding:10px 20px;font-family:'Raleway',sans-serif;" align="left">
        
  <div style="font-size: 14px; color: #ffffff; line-height: 140%; text-align: center; word-wrap: break-word;">
    <p style="font-size: 14px; line-height: 140%;"><span style="font-size: 16px; line-height: 22.4px;">That’s right! It’s not just your birthday,</span></p>
<p style="font-size: 14px; line-height: 140%;"><span style="font-size: 16px; line-height: 22.4px;">but we’re celebrating you for a whole month!</span></p>
  </div>

      </td>
    </tr>
  </tbody>
</table>

  <!--[if (!mso)&(!IE)]><!--></div><!--<![endif]-->
  </div>
</div>
<!--[if (mso)|(IE)]></td><![endif]-->
      <!--[if (mso)|(IE)]></tr></table></td></tr></table><![endif]-->
    </div>
  </div>
  </div>
  


  
  
<div class="u-row-container" style="padding: 0px 10px;background-color: rgba(255,255,255,0)">
  <div class="u-row" style="margin: 0 auto;min-width: 320px;max-width: 600px;overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;background-color: transparent;">
    <div style="border-collapse: collapse;display: table;width: 100%;height: 100%;background-color: transparent;">
      <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding: 0px 10px;background-color: rgba(255,255,255,0);" align="center"><table cellpadding="0" cellspacing="0" border="0" style="width:600px;"><tr style="background-color: transparent;"><![endif]-->
      
<!--[if (mso)|(IE)]><td align="center" width="600" style="width: 600px;padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;" valign="top"><![endif]-->
<div class="u-col u-col-100" style="max-width: 320px;min-width: 600px;display: table-cell;vertical-align: top;">
  <div style="height: 100%;width: 100% !important;">
  <!--[if (!mso)&(!IE)]><!--><div style="box-sizing: border-box; height: 100%; padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;"><!--<![endif]-->
  
  <!--[if (!mso)&(!IE)]><!--></div><!--<![endif]-->
  </div>
</div>
<!--[if (mso)|(IE)]></td><![endif]-->
      <!--[if (mso)|(IE)]></tr></table></td></tr></table><![endif]-->
    </div>
  </div>
  </div>
  


  
  
<div class="u-row-container" style="padding: 30px;background-color: #f0f0f0">
  <div class="u-row" style="margin: 0 auto;min-width: 320px;max-width: 600px;overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;background-color: transparent;">
    <div style="border-collapse: collapse;display: table;width: 100%;height: 100%;background-color: transparent;">
      <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding: 30px;background-color: #f0f0f0;" align="center"><table cellpadding="0" cellspacing="0" border="0" style="width:600px;"><tr style="background-color: transparent;"><![endif]-->
      
<!--[if (mso)|(IE)]><td align="center" width="600" style="width: 600px;padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;" valign="top"><![endif]-->
<div class="u-col u-col-100" style="max-width: 320px;min-width: 600px;display: table-cell;vertical-align: top;">
  <div style="height: 100%;width: 100% !important;">
  <!--[if (!mso)&(!IE)]><!--><div style="box-sizing: border-box; height: 100%; padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;"><!--<![endif]-->
  
<table style="font-family:'Raleway',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
  <tbody>
    <tr>
      <td style="overflow-wrap:break-word;word-break:break-word;padding:20px;font-family:'Raleway',sans-serif;" align="left">
        
  <div style="font-size: 14px; line-height: 120%; text-align: left; word-wrap: break-word;">
    <div style="font-family: arial, helvetica, sans-serif;"><span style="font-size: 12px; color: #999999; line-height: 14.4px;">You received this email because you signed up for My Company Inc.</span></div>
<div style="font-family: arial, helvetica, sans-serif;">&nbsp;</div>
<div style="font-family: arial, helvetica, sans-serif;"><span style="font-size: 12px; color: #999999; line-height: 14.4px;">Unsubscribe</span></div>
  </div>

      </td>
    </tr>
  </tbody>
</table>

  <!--[if (!mso)&(!IE)]><!--></div><!--<![endif]-->
  </div>
</div>
<!--[if (mso)|(IE)]></td><![endif]-->
      <!--[if (mso)|(IE)]></tr></table></td></tr></table><![endif]-->
    </div>
  </div>
  </div>
  


    <!--[if (mso)|(IE)]></td></tr></table><![endif]-->
    </td>
  </tr>
  </tbody>
  </table>
  <!--[if mso]></div><![endif]-->
  <!--[if IE]></div><![endif]-->
</body>

</html>`

func SendEmailReg(item Data) {
	subject := "C др!"
	first_name := item.First_name
	last_name := item.Last_name

	replacer := strings.NewReplacer("${first_name}", first_name, "${last_name}", last_name)

	htmlBytes, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Println("Ошибка при чтении файла index.html:", err)
		return
	}
	html := string(htmlBytes)
	html = replacer.Replace(html)
	log.Println("=fba203=", html)

	m := gomail.NewMessage()
	m.SetHeader("From", "support@crypto-emergency.com")
	log.Println("=30747d=", item.Email)
	m.SetHeader("To", item.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	d := gomail.NewDialer("smtp.yandex.ru", 465, "support@crypto-emergency.com", os.Getenv("EMAIL_PASS"))
	if err := d.DialAndSend(m); err != nil {
		// log.Println("Error SendEmailReg", err)
		time.Sleep(10 * time.Second)
		log.Fatal(err)
	}

}

// log.Println("=fba203=", reflect.TypeOf(today))
// date := "11/06/1969"
// t, err := time.Parse("02/01/2006", date)
// fmt.Println("=008c37=", "День", t.Day(), "Месяц", t.Month())
// fmt.Println("=008c37=", "День", today.Day(), "Месяц", today.Month())

// reflect.TypeOf(tst)
// \n
// today := time.Now()
// today.Day() == birthdate.Day() && today.Month()
