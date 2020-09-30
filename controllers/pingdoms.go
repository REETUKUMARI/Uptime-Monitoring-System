package controllers

import (
	"fmt"
	"net/http"
	"time"
	"website_sc/models"

	"github.com/gin-gonic/gin"
)

func GetUrls(c *gin.Context) {
	var url []models.Pingdom
	if err := repo.DatabaseGets(&url); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		models.DB.Find(&url)
		c.JSON(http.StatusOK, url)
	}
}
func GetUrl(c *gin.Context) {
	id := c.Params.ByName("id")
	var url models.Pingdom
	if err := repo.DatabaseGet(&url); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		models.DB.Where("id = ?", id).First(&url)
		c.JSON(http.StatusOK, url)
	}
}
func CreateUrl(c *gin.Context) {
	var url models.Pingdom
	c.BindJSON(&url)
	url.Status = checkurl(url.URLLink, url.CrawlTimeout)
	if err := repo.DatabaseCreate(&url); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		models.DB.Create(&url)
		c.JSON(http.StatusOK, url)
	}
}
func Updateurl(c *gin.Context) {
	var url models.Pingdom
	id := c.Params.ByName("id")

	if err := models.DB.Where("id = ?", id).First(&url).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.BindJSON(&url)
	err := repo.DatabaseSave(&url)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		models.DB.Save(&url)

		c.JSON(http.StatusOK, url)
	}
}
func Deleteurl(c *gin.Context) {
	id := c.Params.ByName("id")
	var url models.Pingdom
	if err := repo.DatabaseDelete(&url); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		//d := models.DB.Where("id = ?", id).Delete(&url)
		d := models.DB.Delete(&url)
		fmt.Println(d)
		c.JSON(200, gin.H{"id #" + id: "deleted"})
	}

}

func Checklink() {
	rows, err := models.DB.Raw("select * from pingdoms WHERE status != ?", "inactive").Rows()
	if err != nil {
		fmt.Println("connection failed")
	}

	defer rows.Close()

	var (
		id              string
		urllink         string
		crawltimeout    time.Duration
		frequency       int
		failurethresold int
		status          string
		failurecount    int
	)

	c := make(chan string)
	for rows.Next() {
		rows.Scan(&id, &urllink, &crawltimeout, &frequency, &failurethresold, &status, &failurecount)
		//fmt.Println(id, urllink, crawltimeout, frequency, failurethresold, status, failurecount)
		go checkLink(urllink, c)
	}

	for l := range c {
		go func(link string) {
			time.Sleep(1 * time.Second)
			if link != "" {
				checkLink(link, c)
			}
		}(l)
	}
}

func checkurl(url string, crawltimeout time.Duration) string {
	client := http.Client{
		Timeout: crawltimeout * time.Second,
	}
	_, err := client.Get(url)
	if err != nil {
		return "inactive"
	} else {
		return "active"
	}
}

func checkLink(urllink string, c chan string) {
	var p models.Pingdom
	models.DB.First(&p, "url_link  = ?", urllink)
	client := http.Client{
		Timeout: p.CrawlTimeout * time.Second,
	}
	_, err := client.Get(urllink)
	if err != nil {
		failurecount := p.FailureCount + 1
		//fmt.Println(p.FailureCount)
		models.DB.Model(&p).Update("FailureCount", failurecount)
		models.DB.Model(&p).Update("Status", "inactive")
		if failurecount >= p.FailureThreshold {
			models.DB.Model(&p).Update("FailureCount", 0)
			c <- ""
		}
		//fmt.Println(urllink, "is down!")
	} else {
		if p.Status != "active" {
			models.DB.Model(&p).Update("Status", "active")
			models.DB.Model(&p).Update("FailureCount", 0)
		}
		//fmt.Println(urllink, "is up!")
	}
	c <- urllink
}
