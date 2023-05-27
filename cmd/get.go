/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type TrackReponse struct {
	Query  string `json:"query"`
	Tracks struct {
		TotalCount int `json:"totalCount"`
		Items      []struct {
			Data struct {
				URI          string `json:"uri"`
				ID           string `json:"id"`
				Name         string `json:"name"`
				AlbumOfTrack struct {
					URI      string `json:"uri"`
					Name     string `json:"name"`
					CoverArt struct {
						Sources []struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"sources"`
					} `json:"coverArt"`
					ID          string `json:"id"`
					SharingInfo struct {
						ShareURL string `json:"shareUrl"`
					} `json:"sharingInfo"`
				} `json:"albumOfTrack"`
				Artists struct {
					Items []struct {
						URI     string `json:"uri"`
						Profile struct {
							Name string `json:"name"`
						} `json:"profile"`
					} `json:"items"`
				} `json:"artists"`
				ContentRating struct {
					Label string `json:"label"`
				} `json:"contentRating"`
				Duration struct {
					TotalMilliseconds int `json:"totalMilliseconds"`
				} `json:"duration"`
				Playability struct {
					Playable bool `json:"playable"`
				} `json:"playability"`
			} `json:"data"`
		} `json:"items"`
		PagingInfo struct {
			NextOffset int `json:"nextOffset"`
			Limit      int `json:"limit"`
		} `json:"pagingInfo"`
	} `json:"tracks"`
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all paths",
	Long: `This command returns a list of features provided by this cli app
	All paths of the spotify api
`,
	Run: func(cmd *cobra.Command, args []string) {

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		var response TrackReponse

		query, err := cmd.Flags().GetString("track")

		url := "https://spotify23.p.rapidapi.com/search/?q=" + query + "&type=tracks&offset=0&limit=5&numberOfTopResults=5"
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("X-RapidAPI-Key", os.Getenv("APIKEY"))
		req.Header.Add("X-RapidAPI-Host", os.Getenv("APIHOST"))

		fmt.Println("Fetching...")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("An error occurred: ", err.Error())
		}

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		json.Unmarshal(body, &response)

		for index := range response.Tracks.Items {
			artists := response.Tracks.Items[index].Data.Artists.Items
			song := response.Tracks.Items[index].Data.Name
			fmt.Printf("Song: %s |", song)
			fmt.Printf("Artists: ")
			for ind := range artists {
				fmt.Println(artists[ind].Profile.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().String("track", "", "Enter the track name you want to search")
}
