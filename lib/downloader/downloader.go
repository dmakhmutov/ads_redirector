package downloader

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/telegram-bot-api.v4"
	"io"
	"net/http"
	"os"
)

const RESULT_FILE = "result.txt"

func GetFile(fileID string, bot *tgbotapi.BotAPI) {
	var (
		f   tgbotapi.File
		err error
	)

	fc := tgbotapi.FileConfig{FileID: fileID}

	if f, err = bot.GetFile(fc); err != nil {
		log.Errorf("Unable to get file FileID [%s]: %s", fileID, err)
		return
	}

	if err = download(f.Link(viper.GetString("TELEGRAM_KEY")), RESULT_FILE); err != nil {
		log.Errorf("Unable to download file for FileID [%s]: %s", fileID, err)
		return
	}

	log.Debugf("File downloaded for FileID [%s] in %s", fileID, RESULT_FILE)
}

func download(url string, filename string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}
	return
}
