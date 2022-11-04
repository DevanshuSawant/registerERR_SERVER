package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func firebaseChecker() {

        ctxS := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	clientS := getClient(config)

	srvS, err := sheets.NewService(ctxS, option.WithHTTPClient(clientS))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "1tzOP-w3eW9qa6DWPD1D9kntFmRf4mmy9W5kKw0Jek8s"


	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
        

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	it := client.Collection("EVENT1").Snapshots(ctx)
	defer it.Stop()

	for {
		snap, err := it.Next()
		// DeadlineExceeded will be returned when ctx is cancelled.
		if status.Code(err) == codes.DeadlineExceeded {
			return
		}
		if err != nil {
			fmt.Printf("Snapshots.Next: %v", err)
			return
		}
		if snap != nil {
			for _, change := range snap.Changes {
				switch change.Kind {
				case firestore.DocumentAdded:
                                        writeRange := "EVENT1!A:A" 
                                        values := []interface{}{change.Doc.Data()["reg_no"], change.Doc.Data()["name"], change.Doc.Data()["branch"], change.Doc.Data()["batch"], change.Doc.Data()["roll_no"], change.Doc.Data()["batch_roll_no5"], change.Doc.Data()["smart_card_no"], change.Doc.Data()["gmail"], change.Doc.Data()["msteams_mail"], change.Doc.Data()["phone_no"]}
                                        var vr sheets.ValueRange
                                        vr.Values = append(vr.Values, values)
                                        _, err = srvS.Spreadsheets.Values.Append(spreadsheetId, writeRange, &vr).ValueInputOption("USER_ENTERED").Do()
                                        if err != nil {
                                                fmt.Println(err)
                                        }
                                        fmt.Printf("Added: %v\n", change.Doc.Data()["reg_no"])
				case firestore.DocumentRemoved:
					fmt.Printf("Removed: %v\n", change.Doc.Data()["reg_no"])
				}
			}
		}
	}
}

//ValueRange Struct to hold the values to be written to the spreadsheet.

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func spreadsheets() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "1tzOP-w3eW9qa6DWPD1D9kntFmRf4mmy9W5kKw0Jek8s"

	writeRange := "LOOKUP!A:A" 
	values := []interface{}{"16105", "hi", "hello"}
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)
	_, err = srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	firebaseChecker()
	spreadsheets()
}
