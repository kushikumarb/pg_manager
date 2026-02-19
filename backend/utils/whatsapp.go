package utils

import (
	"fmt"
	"pg-manager-backend/config"// Import your new config package
	"strings"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendWhatsAppMessage(to string, messageBody string) error {
	// 1. Use the typed config instead of os.Getenv
	accountSid := config.App.TwilioSID
	authToken := config.App.TwilioAuthToken
	fromPhone := config.App.TwilioFromNumber

	// Validation check before attempting to send
	if accountSid == "" || authToken == "" {
		return fmt.Errorf("twilio credentials are not configured in config.go")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// 2. Format the phone number
	formattedTo := to
	if !strings.HasPrefix(to, "+") {
		// Professional tip: Use +91 for your specific region as default
		formattedTo = "+91" + to
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo("whatsapp:" + formattedTo)
	params.SetFrom(fromPhone)
	params.SetBody(messageBody)

	// 3. Execute the request
	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("twilio api error: %v", err)
	}

	if resp.Sid != nil {
		fmt.Printf("✅ WhatsApp sent to %s! SID: %s\n", formattedTo, *resp.Sid)
	}
	return nil
}
