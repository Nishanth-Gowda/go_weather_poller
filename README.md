# Weather Polling App

The Weather Polling App is a simple Go application that retrieves weather data from an external API and sends it via SMS using Twilio. It provides an example of how to structure a Go application, use interfaces, and integrate with external services.

## Features

- Periodically fetches weather data (temperature, rain, showers, etc.) from an external API.
- Sends the weather data via SMS using Twilio.
- Supports extensibility for adding more notification methods.

## Getting Started

### Prerequisites

- Go installed on your machine
- Twilio account with credentials

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/weather-polling.git
cd weather-polling
```

2. Create a .env file with your Twilio credentials:

```
ACCOUNT_SID=your_account_sid
AUTH_TOKEN=your_auth_token
FROM_PHONE=your_twilio_phone_number
TO_PHONE=your_destination_phone_number
```

3. Build and run the application:

```
go build -o weather-polling
./weather-polling
```

4. Configuration
You can configure the polling interval and other parameters in the main.go file.

5. Extending Notification Methods
To extend the application to support additional notification methods, you can implement the Sender interface. The SMSSender is an example implementation.

6. Contributing
Feel free to contribute by opening issues or submitting pull requests. Your feedback and contributions are welcome!


Make sure to replace placeholders like `your_account_sid`, `your_auth_token`, etc., wit
