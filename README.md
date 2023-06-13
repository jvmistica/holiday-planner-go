# holiday-planner-go
![Build status](https://github.com/jvmistica/holiday-planner-go/workflows/holiday-planner-go/badge.svg)
![Coverage](https://img.shields.io/sonar/coverage/jvmistica_holiday-planner-go/main?server=https%3A%2F%2Fsonarcloud.io)

A tool that fetches holidays from Google Calendar and adds them together with weekends to provide vacation leave suggestions.


## Requirements
#### Google API key

To get a Google API key:
1. Go to https://console.cloud.google.com
2. Navigate to "APIs & Services" -> "Credentials"
3. Click "Create Credentials" -> "API key"


#### Trello API key and token

To get Trello API key and token:
1. Go to https://trello.com/power-ups/admin
2. Select an existing power-up and integration or create a new one
3. Select "API key" -> "Generate a new API key"
4. Follow the link to generate a token

### Environment Variables
```
export GCP_API_KEY=<gcp-api-key>
export TRELLO_API_KEY=<trello-api-key>
export TRELLO_API_TOKEN=<trello-api-token>
```

## Usage
`go run main.go -start=2023-06-01 -end=2024-01-31`
