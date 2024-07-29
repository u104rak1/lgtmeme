# LGTMeme
LGTMeme is a simple LGTM (Looks Good To Me) image generator. You can quickly insert LGTM characters into your favorite image and copy markdown to the clipboard. Brighten up the monotonous code review and approval process with humorous LGTM image.

## URL
https://lgtmeme.onrender.com

## Development Background
Initially, I was playing around with creating an Auth server while learning about Open ID Connect, with the goal of deepening my understanding of it. I thought it would be a waste to just run it in a local environment after making it, so I also created a front end to show it to people. I chose the LGTM image generation service as it was something that I had a minimum of need for, was not complicated, and could use the authorization function I created. Read, Write, and Update are authorized with Client Credentials Grant, and Delete is authorized with Authorization Code Grant. Also, since I wanted to use the deployment service for free, the front, Auth, and Resource servers are not independent but are combined into one, and the authorization sequence is reproduced in the form of an internal server request.

## Documents
- [Folder configuration diagram](https://github.com/ucho456job/lgtmeme/blob/develop/document/folder_configuration.md)
- [Endpoint list](https://github.com/ucho456job/lgtmeme/blob/develop/document/endpoint.md)
- [Architecture diagram](https://github.com/ucho456job/lgtmeme/blob/develop/document/architecture.md)
- [ER diagram](https://github.com/ucho456job/lgtmeme/blob/develop/document/er.md)
- [Authorization Code Grant sequence diagram](https://github.com/ucho456job/lgtmeme/blob/develop/document/authorization_code_grant_flow.md)
- [Refresh Token Grant sequence diagram](https://github.com/ucho456job/lgtmeme/blob/develop/document/refresh_token_grant_flow.md)
- [Client Credentials Garnt sequence diagram](https://github.com/ucho456job/lgtmeme/blob/develop/document/client_credentials_grant_flow.md)
- [Technology selection](https://github.com/ucho456job/lgtmeme/blob/develop/document/Technology_selection.md)
- [Developer knowledge](https://github.com/ucho456job/lgtmeme/blob/develop/document/developer_knowledge.md)
