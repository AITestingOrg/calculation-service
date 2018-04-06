package tests

func TestLogin(t *testing.T) {

	// Create Pact, connecting to local Daemon
	// Ensure the port matches the daemon port!
	pact := Pact{
		Port:     8080,
		Consumer: "calculation-service",
		Provider: "gmaps-adapter",
	}
	// Shuts down Mock Service when done
	defer pact.Teardown()

	// Pass in your test case as a function to Verify()
	var test = func() error {
		_, err := http.POST("http://localhost:8080/")
		return err
	}

	// Set up our interactions. Note we have multiple in this test case!
	pact.
		AddInteraction().
		Given("a post request from weston, FL to Miami Lakes, FL"). // Provider State
		UponReceiving("A request to login"). // Test Case Name
		WithRequest(Request{
			Method: "POST",
			Path:   "/api/v1/directions",
			Body: {
				origin: "weston, fl",
				destination: "Miami lakes, fl",
				departure_time: 15220998650000000
			},
			Headers: {"Accept" => "application/json"}
		}).
		WillRespondWith(Response{
			Status: 200,
			Headers: {"Content-Type" => "application/json"},
			Body: {
				duration: Pact.like(1666),
				distance: Pact.like(34101)
			}
		})

	// Run the test and verify the interactions.
	err := pact.Verify(test)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}

	// Write pact to file
	pact.WritePact()
}
