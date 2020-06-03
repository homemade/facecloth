package flannel

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
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

	var status int
	var result map[string]interface{}
	status, result, err = c.CreateFundraiser(params)
	checkResult := func(testcase string, status int, result map[string]interface{}, err error) {
		if err != nil {
			t.Errorf("failed to create %s %v", testcase, err)
		}
		if _, exists := result["id"]; exists {
			t.Logf("succesfully created %s %d %v", testcase, status, result["id"])
		} else {
			t.Errorf("invalid result returned when creating %s %d %v", testcase, status, result)
		}
	}
	checkResult("fundraiser with required fields", status, result, err)

	externalID = fmt.Sprintf("%s_2", ts)
	params.Title = fmt.Sprintf("Test Fundraiser %s", externalID)
	params.Description = fmt.Sprintf("The description for Test Fundraiser %s", externalID)

	// test an image within size limit (< 4MB)
	httpClient := &http.Client{Timeout: time.Second * 20}
	res, err := httpClient.Get("https://images.unsplash.com/photo-1576086265779-619d2f54d96b")
	if err != nil {
		t.Fatalf("failed to download a test cover photo image %v", err)
	}
	defer res.Body.Close()
	status, result, err = c.CreateFundraiser(params, WithFundraiserCoverPhotoImage("photo-1576086265779-619d2f54d96b", res.Body))
	checkResult("fundraiser with cover photo image", status, result, err)

	externalID = fmt.Sprintf("%s_3", ts)
	params.Title = fmt.Sprintf("Test Fundraiser %s", externalID)
	params.Description = fmt.Sprintf("The description for Test Fundraiser %s", externalID)
	status, result, err = c.CreateFundraiser(params,
		WithFundraiserField("external_fundraiser_uri", "http://www.example.com/"),                                             // URI of the fundraiser on the external site.
		WithFundraiserField("external_event_name", fmt.Sprintf("The external event name for Test Fundraiser %s", externalID)), // Name of the event this fundraiser belongs to.
		WithFundraiserField("external_event_uri", "http://www.example.org/"),                                                  // URI of the event this fundraiser belongs to.
		WithFundraiserField("external_event_start_time", fmt.Sprintf("%d", time.Now().AddDate(0, 0, 1).Unix())),               // Unix timestamp of the day when the event takes place.
	)
	checkResult("fundraiser with optional fields", status, result, err)

	externalID = fmt.Sprintf("%s_4", ts)
	params.Title = fmt.Sprintf("Test Fundraiser %s", externalID)
	params.Description = fmt.Sprintf("The description for Test Fundraiser %s", externalID)
	// test an image url within size limit (< 4MB)
	var coverPhotoURL *url.URL
	coverPhotoURL, err = url.Parse("https://images.unsplash.com/photo-1576086265779-619d2f54d96b")
	if err != nil {
		t.Fatalf("failed to parse test cover photo url %v", err)
	}
	status, result, err = c.CreateFundraiser(params, WithFundraiserCoverPhotoURL(path.Base(coverPhotoURL.Path), *coverPhotoURL))
	checkResult("fundraiser with cover photo url", status, result, err)

	externalID = fmt.Sprintf("%s_5", ts)
	params.Title = fmt.Sprintf("Test Fundraiser %s", externalID)
	params.Description = fmt.Sprintf("The description for Test Fundraiser %s", externalID)
	// test an image url over the size limit (> 4MB)
	coverPhotoURL, err = url.Parse("https://images.unsplash.com/photo-1576086476234-1103be98f096")
	if err != nil {
		t.Fatalf("failed to parse test cover photo url %v", err)
	}
	status, result, err = c.CreateFundraiser(params, WithFundraiserCoverPhotoURL(path.Base(coverPhotoURL.Path), *coverPhotoURL))
	if IsErrorWithFundraiserCoverPhoto(err) {
		t.Logf("fundraiser with cover photo url over the size limit returned an error as expected: %v", err)
	}

	externalID = fmt.Sprintf("%s_6", ts)
	params.Title = fmt.Sprintf("Test Fundraiser %s", externalID)
	params.Description = fmt.Sprintf("The description for Test Fundraiser %s", externalID)
	// test an image url over the dimensions limit
	coverPhotoURL, err = url.Parse("https://via.placeholder.com/30001x1.jpg")
	if err != nil {
		t.Fatalf("failed to parse cover photo url %v", err)
	}
	status, result, err = c.CreateFundraiser(params, WithFundraiserCoverPhotoURL(path.Base(coverPhotoURL.Path), *coverPhotoURL))
	if IsErrorWithFundraiserCoverPhoto(err) {
		t.Logf("fundraiser with cover photo url containing too many pixels returned an error as expected: %v", err)
	}

}
