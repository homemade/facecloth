package flannel

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCreateFundraiser(t *testing.T) {

	accessToken := os.Getenv("ACCESS_TOKEN")
	charityID := os.Getenv("CHARITY_ID")
	ts := time.Now().Format("20060102150405")

	c, err := CreateAPIClient(WithLogger(t, true))
	if err != nil {
		t.Fatalf("failed to create api client with test logger %v", err)
	}

	externalID := fmt.Sprintf("%s_1", ts)
	params := CreateFundraiserParams{
		AccessToken: accessToken,
		CharityID:   charityID,
		Title:       fmt.Sprintf("Test Fundraiser %s", externalID),
		Description: fmt.Sprintf("The description for Test Fundraiser %s", externalID),
		Goal:        100000,
		Currency:    "GBP",
		EndTime:     time.Now().AddDate(1, 0, 0),
		ExternalID:  externalID,
	}

	var id string
	var status int
	id, status, err = c.CreateFundraiser(params)
	if err != nil {
		t.Fatalf("failed to create fundraiser with required fields %v", err)
	}
	t.Logf("succesfully created fundraiser with required fields %s %d", id, status)

	externalID = fmt.Sprintf("%s_2", ts)
	params.Title = fmt.Sprintf("Test Fundraiser %s", externalID)
	params.Description = fmt.Sprintf("The description for Test Fundraiser %s", externalID)

	coverPhoto, err := os.Open("./FundraiserCoverPhoto.jpg")
	if err != nil {
		t.Fatalf("failed to read cover photo %v", err)
	}
	defer coverPhoto.Close()
	id, status, err = c.CreateFundraiser(params, WithFundraiserCoverPhotoImage("FundraiserCoverPhoto.jpg", coverPhoto))
	if err != nil {
		t.Fatalf("failed to create fundraiser with cover photo %v", err)
	}
	t.Logf("succesfully created fundraiser with cover photo %s %d", id, status)

	externalID = fmt.Sprintf("%s_3", ts)
	params.Title = fmt.Sprintf("Test Fundraiser %s", externalID)
	params.Description = fmt.Sprintf("The description for Test Fundraiser %s", externalID)

	id, status, err = c.CreateFundraiser(params,
		WithFundraiserField("external_fundraiser_uri", "http://www.example.com/"),                                             // URI of the fundraiser on the external site.
		WithFundraiserField("external_event_name", fmt.Sprintf("The external event name for Test Fundraiser %s", externalID)), // Name of the event this fundraiser belongs to.
		WithFundraiserField("external_event_uri", "http://www.example.org/"),                                                  // URI of the event this fundraiser belongs to.
		WithFundraiserField("external_event_start_time", fmt.Sprintf("%d", time.Now().AddDate(0, 0, 1).Unix())),               // Unix timestamp of the day when the event takes place.
	)
	if err != nil {
		t.Fatalf("failed to create fundraiser with optional fields %v", err)
	}
	t.Logf("succesfully created fundraiser with optional fields %s %d", id, status)

}
