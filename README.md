# SMS-Gateway microservice app

SMS gateway microservice application is a gateway for send and manag SMS.
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

| provider | address |
| --- | --- |
| **rahyab** | [sms.rahyab.ir](http://sms.rahyab.ir) |
| **rahyabPG** | [rahyabcp.ir](https://rahyabcp.ir)|
| **kavenegar** | [kavenegar.com](http://kavenegar.com)|
| **hamyarsms** | [hamyarsms.com](https://hamyarsms.com)|
| **More...** | **Coming soon**|

## How to start

**build and startup project .**

```
git clone https://github.com/mohaphez/sms-gateway-Golang.git
```

1. Go to the project folder

```
cd sms-gateway-Golang/project
```
2. Change .env Content for more security

3. Fill sms-service/config.js parameters for which provider you want use. 

4. Enter below command

```
make up_build
```

**stop and remove containers**

```
make down
```

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
