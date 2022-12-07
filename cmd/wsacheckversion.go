/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	. "labs/utils"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type WSAVersionCheck struct {
	currentVersion string
	cancel         context.CancelFunc
	ctx            context.Context
}

const (
	RemoteURL = "https://itunes.apple.com/cn/lookup?bundleId=net.whatsapp.WhatsApp"
)

func PushZoomWebHookMessage(message string) error {
	type pushReq struct {
		Version string `json:"version"`
	}

	url := "https://inbots.zoom.us/incoming/hook/Rx7c-Ks7QBDzlRAsKKlJsd0-?format=fields"
	opts := []RequestOptions{
		WithHeaders(map[string]string{
			"Authorization": "rk6uerDtimnjPt0un65Xr60g",
			"Content-Type":  "application/json",
		}),
	}

	req, _ := json.Marshal(&pushReq{
		Version: message,
	})
	body := bytes.NewBuffer(req)
	buffer, err := HttpRequest(url, "POST", body, opts...)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("[time]:%v push [message]:%v result:%v.", time.Now().Unix(), message, string(buffer)))
	return nil
}

func CreateUpgradeChecker() *WSAVersionCheck {
	return &WSAVersionCheck{
		currentVersion: "0.0.0",
	}
}

func (wsa *WSAVersionCheck) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wsa.cancel, wsa.ctx = cancel, ctx

	var (
		version     string
		err         error
		needUpgrade = false
	)
	defer func() {
		if err != nil {
			PushZoomWebHookMessage(fmt.Sprintf("WSAVersionCheck ERROR:%v", err.Error()))
			return
		}

		PushZoomWebHookMessage("WSAVersionCheck OVER")
	}()

	ticker := time.NewTicker(60 * time.Minute)

	version, err = wsa.checkVersionUpgrade()
	if err == nil {
		_, err = wsa.checkVersionNotify(version)
	}
	if err != nil {
		return
	}
	PushZoomWebHookMessage(fmt.Sprintf("WSAVersionCheck START, version:%v", version))

	for true {
		select {
		case <-ticker.C:
			version, err = wsa.checkVersionUpgrade()
			if err != nil {
				fmt.Println("check upgrade", err.Error())
				break
			}

			needUpgrade, err = wsa.checkVersionNotify(version)
			if err != nil {
				fmt.Println("check notify", err.Error())
				break
			}

			if needUpgrade {
				PushZoomWebHookMessage(fmt.Sprintf("WhatsApp Upgraded, Version:%v", version))
			}
		case <-ctx.Done():
			return
		}
	}
}

func (wsa *WSAVersionCheck) Cancel() {
	if wsa.cancel != nil {
		wsa.cancel()
	}
}

func (wsa *WSAVersionCheck) checkVersionUpgrade() (max string, err error) {
	type VersionUpgradeResult struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			ArtistViewUrl                      string        `json:"artistViewUrl"`
			ArtworkUrl60                       string        `json:"artworkUrl60"`
			ArtworkUrl100                      string        `json:"artworkUrl100"`
			IsGameCenterEnabled                bool          `json:"isGameCenterEnabled"`
			Advisories                         []string      `json:"advisories"`
			Features                           []interface{} `json:"features"`
			SupportedDevices                   []string      `json:"supportedDevices"`
			IpadScreenshotUrls                 []interface{} `json:"ipadScreenshotUrls"`
			AppletvScreenshotUrls              []interface{} `json:"appletvScreenshotUrls"`
			ArtworkUrl512                      string        `json:"artworkUrl512"`
			ScreenshotUrls                     []string      `json:"screenshotUrls"`
			Kind                               string        `json:"kind"`
			TrackCensoredName                  string        `json:"trackCensoredName"`
			TrackViewUrl                       string        `json:"trackViewUrl"`
			ContentAdvisoryRating              string        `json:"contentAdvisoryRating"`
			AverageUserRating                  float64       `json:"averageUserRating"`
			ReleaseNotes                       string        `json:"releaseNotes"`
			IsVppDeviceBasedLicensingEnabled   bool          `json:"isVppDeviceBasedLicensingEnabled"`
			ReleaseDate                        time.Time     `json:"releaseDate"`
			Description                        string        `json:"description"`
			SellerName                         string        `json:"sellerName"`
			GenreIds                           []string      `json:"genreIds"`
			TrackId                            int           `json:"trackId"`
			TrackName                          string        `json:"trackName"`
			BundleId                           string        `json:"bundleId"`
			PrimaryGenreName                   string        `json:"primaryGenreName"`
			PrimaryGenreId                     int           `json:"primaryGenreId"`
			CurrentVersionReleaseDate          time.Time     `json:"currentVersionReleaseDate"`
			Currency                           string        `json:"currency"`
			MinimumOsVersion                   string        `json:"minimumOsVersion"`
			LanguageCodesISO2A                 []string      `json:"languageCodesISO2A"`
			FileSizeBytes                      string        `json:"fileSizeBytes"`
			SellerUrl                          string        `json:"sellerUrl"`
			FormattedPrice                     string        `json:"formattedPrice"`
			AverageUserRatingForCurrentVersion float64       `json:"averageUserRatingForCurrentVersion"`
			UserRatingCountForCurrentVersion   int           `json:"userRatingCountForCurrentVersion"`
			TrackContentRating                 string        `json:"trackContentRating"`
			Price                              float64       `json:"price"`
			ArtistId                           int           `json:"artistId"`
			ArtistName                         string        `json:"artistName"`
			Genres                             []string      `json:"genres"`
			Version                            string        `json:"version"`
			WrapperType                        string        `json:"wrapperType"`
			UserRatingCount                    int           `json:"userRatingCount"`
		} `json:"results"`
	}

	var buff []byte
	buff, err = HttpRequest(RemoteURL, "GET", nil)
	if err != nil {
		return
	}

	var resp VersionUpgradeResult
	err = json.Unmarshal(buff, &resp)
	if err != nil {
		return
	}

	for counter := 0; counter != resp.ResultCount; counter++ {
		result := resp.Results[counter]

		if max == "" || max < result.Version {
			max = result.Version
		}
	}
	return
}

func (wsa *WSAVersionCheck) checkVersionNotify(newVersion string) (needUpgrade bool, err error) {
	needUpgrade = false

	v0 := strings.Split(wsa.currentVersion, ".")
	v1 := strings.Split(newVersion, ".")

	if len(v0) != 3 || len(v1) != 3 {
		err = errors.New("error version format")
		return
	}

	current, err0 := strconv.Atoi(v0[1])
	if err0 != nil {
		err = errors.New("error version format")
		return
	}

	upgrade, err1 := strconv.Atoi(v1[1])
	if err1 != nil {
		err = errors.New("error version format")
		return
	}

	if upgrade > current {
		wsa.currentVersion = newVersion
		needUpgrade = true
	}
	return
}

// wsacheckversionCmd represents the wsacheckversion command
var wsacheckversionCmd = &cobra.Command{
	Use:   "wsacheckversion",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		checker := CreateUpgradeChecker()

		go func() {
			signalChannel := make(chan os.Signal, 1)
			signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
			<-signalChannel

			checker.Cancel()
		}()

		checker.Run()
		//fmt.Println((&WSAVersionCheck{}).checkVersionUpgrade())
	},
}

func init() {
	rootCmd.AddCommand(wsacheckversionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wsacheckversionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wsacheckversionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
