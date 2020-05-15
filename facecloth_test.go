package facecloth

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCreateFundraiser(t *testing.T) {

	accessToken := os.Getenv("ACCESS_TOKEN")
	charityID := os.Getenv("CHARITY_ID")

	externalID := time.Now().Format("20060102150405")

	c, err := CreateAPIClient()
	if err != nil {
		t.Fatalf("failed to create api client %v", err)
	}

	var id string
	var status int
	id, status, err = c.CreateFundraiser(CreateFundraiserParams{
		AccessToken: accessToken,
		CharityID:   charityID,
		Title:       fmt.Sprintf("Test Fundraiser %s", externalID),
		Description: fmt.Sprintf("The description for Test Fundraiser %s", externalID),
		Goal:        100000,
		Currency:    "GBP",
		EndTime:     time.Now().AddDate(1, 0, 0),
		ExternalID:  externalID,
	})
	if err != nil {
		t.Fatalf("failed to create fundraiser %v", err)
	}

	t.Logf("result %s %d", id, status)

}
