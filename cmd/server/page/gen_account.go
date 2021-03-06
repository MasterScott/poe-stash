package page

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/cptpingu/poe-stash/generate"
	"github.com/cptpingu/poe-stash/scraper"
)

// GenAccountHandler handles refresh of an account.
func GenAccountHandler(c *gin.Context) {
	verbosity := c.MustGet("verbosity").(int)

	account := c.Params.ByName("account")
	poeSessID := ""
	realm := "pc"
	league := "standard"

	query := c.Request.URL.Query()
	if param, ok := query["poesessid"]; ok {
		poeSessID = param[0]
	}
	if param, ok := query["realm"]; ok {
		realm = param[0]
	}
	if param, ok := query["league"]; ok {
		league = param[0]
	}

	// It could be a refresh, try to deduce it from previous run.
	if poeSessID == "" {
		fmt.Println("no poesessid for account ", account, ", let's deduce it!")
		if sess, errFile := ioutil.ReadFile(scraper.DataCacheDir + account + ".poesessid"); errFile != nil {
			fmt.Println("can't find any cached poesessid for account", account)
		} else {
			poeSessID = string(sess)
		}
	}

	scrap := scraper.NewScraper(account, poeSessID, realm, league)
	scrap.SetVerbosity(verbosity)
	data, errScrap := scrap.ScrapEverything()
	if errScrap != nil {
		c.HTML(http.StatusOK, "error", errScrap)
		return
	}

	name := fmt.Sprintf("%s-%s-%s", account, league, realm)
	output := fmt.Sprintf("%s%s.html", scraper.DataDir, name)
	file, err := os.Create(output)
	if err != nil {
		c.HTML(http.StatusOK, "error", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(file)
	gen := generate.NewGenerator(w)
	if errGen := gen.GenerateHTML(data); errGen != nil {
		c.HTML(http.StatusOK, "error", errGen)
		return
	}
	w.Flush()

	// Everything is fine, let's store the poeSessID.
	if errFile := ioutil.WriteFile(scraper.DataCacheDir+account+".poesessid", []byte(poeSessID), 0644); errFile != nil {
		// Non fatal error, storing the session is not mandatory.
		fmt.Println("error occurred:", errFile)
	}
	c.HTML(http.StatusOK, "redirect", "/view/"+name)
}
