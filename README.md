# SMS-Gateway microservice app

SMS gateway microservice application is a gateway for send and manage SMS.
This microservice will help you to send your SMS as quick as possible with a single API without getting involved in different providers API methods.

## Components

1. [Frontend](/front-end) management panel to display and filter sent text messages and status, send quick text messages, view logs and events recorded in other services. (**It's incomplete and under construction**)
2. [authentication](/authentication-service) managing users (crud), change passwords, authenticating and generate jwt tokens
3. [broker](/broker-service) connect services to each other behind the scenes.
4. [logger](/logger-service) recording events and errors related to services.
5. [sms](/sms-service) It connect to different provider and send SMS based on priority.

Take a look at the components diagram that describes them and their interactions.

![sms-gateway-diagram](https://user-images.githubusercontent.com/20874565/202277834-462358da-2143-47ea-8be5-789157d87886.png)

## Use cases

- if you have one or more systems that require SMS services and you want their connect to one API and stored sms data in single database.
- if you use several different SMS service providers.
- if you send a lot of mass SMS.
- if you don't want to change the code of your software when transferring the SMS number to another service provider.
- if sending speed is important to you

## Which providers are available

| provider      | address                                |
| ------------- | -------------------------------------- |
| **rahyab**    | [sms.rahyab.ir](http://sms.rahyab.ir)  |
| **rahyabPG**  | [rahyabcp.ir](https://rahyabcp.ir)     |
| **kavenegar** | [kavenegar.com](http://kavenegar.com)  |
| **hamyarsms** | [hamyarsms.com](https://hamyarsms.com) |
| **More...**   | **Coming soon**                        |

## How to start

1. Clone the project

```
git clone https://github.com/mohaphez/sms-gateway-Golang.git
```

2. Change sms-gateway-Golang/project/.env content for more security

3. Fill sms-gateway-Golang/sms-service/config.js parameters for which provider you want to use.

4. Go to the project folder

```
cd sms-gateway-Golang/project
```
5. Enter below command

```
make up_build
```

## Usage

With the following code samples, you can connect to microservice API and use them in your software.

- How to get authentication token

```
URL : IP:8000/getToken  <!-- replace IP with your local or server ip/domain -->

METHOD: POST

BODY : JSON/raw

PARAMETERS :
  {
    "username":"demo", <!-- string -->
    "password":"demo"  <!-- string -->
  }

RESPONSE :
   {
    "error": false,
    "message": "You Successfully Authenticate",
    "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njg2Nzk1MjAsImkiOjF9.kToEp7c3Sceva7EG"
   }

```

- How to send single sms

```
URL : IP:8000/send-sms  <!-- replace IP with your local or server ip/domain  -->

METHOD: POST

BODY : JSON/raw

HEADER : Authorization: Bearer Token

PARAMETERS :
  {
    "message":"test",           <!-- string -->
    "receptor":["09910000000"], <!-- array[string] -->
    "sender":"rahyabPG",        <!-- string --> <!-- / According to the above table / rahyab|rahyabPG|kavenegar|hamyarsms | .. -->
    "sender_number":"10001222", <!-- string -->
    "localid":""                <!-- string --> <!-- optional -->
  }

RESPONSE :
   {
    "error": false,
    "message": "All Messages sent successfully !",
    "data": [
        {
            "batchid": "b5a98eb2-a8c5-4877-899f-993b9fe5efbf",
            "date": 1668677882,
            "lang": "en",
            "message": "test",
            "receptor": "09910000000",
            "sender": "rahyabPG",
            "sender_number": "10001222",
            "sms_count": 1,
            "status": 1,
            "status_text": "در صف ارسال"
        }
    ]
  }

```

- How to send p2p sms

```
URL : IP:8000/send-sms  <!-- replace IP with your local or server ip/domain  -->

METHOD: POST

BODY : JSON/raw

HEADER : Authorization: Bearer Token

PARAMETERS :
  {
    "message":["num1","num2"], <!-- array[string] -->
    "receptor":["09910000000","09910000001"], <!-- array[string] -->
    "sender":"rahyabPG",      <!-- string --> <!-- / According to the above table / rahyab|rahyabPG|kavenegar|hamyarsms | .. -->
    "sender_number":"10001222", <!-- string -->
    "localid":[]                <!-- array[string] --> <!-- optional -->
  }

RESPONSE :
   {
    "error": false,
    "message": "All Messages sent successfully !",
    "data": [
            {
                "batchid": "cadd891a-2f20-4cf9-ad6f-58a601bfbdbb",
                "date": 1668677982,
                "lang": "en",
                "message": "num1",
                "receptor": "09910000000",
                "sender": "rahyabPG",
                "sender_number": "10001222",
                "sms_count": 1,
                "status": 1,
                "status_text": "در صف ارسال"
            },
            {
                "batchid": "c25fddd4-e529-4cbe-a527-51955ebeae7b",
                "date": 1668677982,
                "lang": "en",
                "message": "num2",
                "receptor": "09910000001",
                "sender": "rahyabPG",
                "sender_number": "10001222",
                "sms_count": 1,
                "status": 1,
                "status_text": "در صف ارسال"
            }
        ]
    }

```

- How to get sms status

```
URL : IP:8000/get-sms-status  <!-- replace IP with your local or server ip/domain  -->

METHOD: POST

BODY : JSON/raw

HEADER : Authorization: Bearer Token

PARAMETERS :
  {
    "batchid":["b5a98eb2-a8c5-4877-899f-993b9fe5efbf"] <!-- array[string] -->
  }

RESPONSE :
   {
    "error": false,
    "message": "Messages fetched successfully !",
    "data": [
        {
            "batchid": "b5a98eb2-a8c5-4877-899f-993b9fe5efbf",
            "created_at": "2022-11-17T09:38:02.325Z",
            "date": 1668677882,
            "lang": "en",
            "message": "test",
            "receive_time": "",
            "receptor": "9910000000",
            "send_time": "",
            "sender": "rahyabPG",
            "sender_number": "10001222",
            "sms_count": 1,
            "status": 1,
            "status_text": "در صف ارسال",
            "updated_at": "2022-11-17T09:38:02.325Z"
        }
    ]
  }

```

- How to get sms list with filter

```
URL : IP:8000/get-sms-list  <!-- replace IP with your local or server ip/domain  -->

METHOD: POST

BODY : JSON/raw

HEADER : Authorization: Bearer Token

PARAMETERS :
  {
	"message":"",      <!-- string -->
	"receptor":[],     <!-- array[string] -->
	"sender":[],       <!-- array[string] -->
	"senderNumber":[], <!-- array[string] -->
	"offset":0,        <!-- number -->
	"limit":10,        <!-- number -->
	"sort":""          <!-- string -->  <!-- "asc" | "desc" -->
   }


RESPONSE :
   {
    "error": false,
    "message": "Messages fetched successfully !",
    "data": [
        {
            "batchid": "b5a98eb2-a8c5-4877-899f-993b9fe5efbf",
            "created_at": "2022-11-17T09:38:02.325Z",
            "date": 1668677882,
            "error": "0",     <!-- This field appears when the provider returns an error code when sending a message, you should find its meaning from the provider own document.-->
            "lang": "en",
            "message": "test",
            "receive_time": "0001-01-01T00:00:00Z",
            "receptor": "9910000000",
            "send_time": "0001-01-01T00:00:00Z",
            "sender": "rahyabPG",
            "sender_number": "10001222",
            "sms_count": 1,
            "status": 1,
            "status_text": "در صف ارسال",
            "updated_at": "2022-11-17T09:38:02.325Z"
        }
    ]
}
```

## How to stop

1. Go to the project folder

```
cd sms-gateway-Golang/project
```

2. Enter below command

```
make down
```

## Roadmap

- [ ] Add listener service with rabbitmq
- [ ] Add automatic check sms status job
- [ ] Create front-end panel ui
  - [ ] Create dashboard with various chart.
  - [ ] Create page for SMS Send list with filter.
  - [ ] Create page for send sms.
  - [ ] Create page for log events list.
- [ ] Add more providers such as sms.ir , melipayamak.com , farapayamak.ir,farazsms.com ,etc...

## Contribution

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement". Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

MIT
