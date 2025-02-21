# Auth API

## Instructions to run Auth application

- If you have go installed, just clone the repo, do `go mod tidy`,and then run using `go run main.go`.
- If you are using Linux OS and not having Go in your system, just run build. Download the build file named `LinuxAuth`, and run using `./LinuxAuth`.
- If you are using MacOs and not having Go installed, run the build file `MacAuth` using `./MacAuth`.
- If you are using Windows and not having Go installed, run the exe file named `WindowsAuth`.

## cURL commands for each testcase

- **SignUp**: `curl --location 'http://localhost:8000/signUp' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAwNzk5NTMsInVzZXJuYW1lIjoibW1nIn0.0Puc0Q9v8dkxajY4_Sning2OMNlRlwUQCDQhT7NYjW0' \
  --header 'Content-Type: application/json' \
  --data '{
  "userName":"mukund",
  "password":"Ariqt"
  }'`



- **SignIn**: `curl --location 'http://localhost:8000/signIn' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAwNzk5NTMsInVzZXJuYW1lIjoibW1nIn0.0Puc0Q9v8dkxajY4_Sning2OMNlRlwUQCDQhT7NYjW0' \
  --header 'Content-Type: application/json' \
  --data '{
  "userName":"mukund",
  "password":"Ariqt"
  }'`


- **Authentication using JWT**: `curl --location 'http://localhost:8000/get' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAxMjk5ODgsImp0aSI6MC40Njg5NjM0MzQ3ODA1NTI1LCJ1c2VybmFtZSI6Im11a3VuZCJ9.-XNEBpN6MNmzpPu-h0KHZ4b3Vx9xPb7hfqGkF-xHICI'`
  
    **NOTE:** Please update the token in above cURL whatever got as a response of signin, because above token would have expired.


- **Revoking a token**: `curl --location 'http://localhost:8000/revokeToken' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAxMjk5ODgsImp0aSI6MC40Njg5NjM0MzQ3ODA1NTI1LCJ1c2VybmFtZSI6Im11a3VuZCJ9.-XNEBpN6MNmzpPu-h0KHZ4b3Vx9xPb7hfqGkF-xHICI'`


- **Getting new token using Refresh Token**: `curl --location 'http://localhost:8000/refreshToken' \
  --header 'Refresh-Token: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDA3MzM0MDgsInVzZXJuYW1lIjoibXVrdW5kIn0.K3b3xExrDfEuzpU4hBdWO628dp-MtCtAt2OfrD9OR0w'`