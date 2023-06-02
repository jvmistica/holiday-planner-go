# gcal
![Build status](https://github.com/jvmistica/gcal/workflows/gcal/badge.svg)

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

### Environment Variables
```
export GCP_API_KEY=<gcp-api-key>
export TRELLO_API_KEY=<trello-api-key>
export TRELLO_SECRET=<trello-secret>
```

## Usage
`go run main.go -start=2023-05-01T00:00:00Z -end=2023-05-31T00:00:00Z`
