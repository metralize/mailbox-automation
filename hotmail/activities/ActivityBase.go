package activities

import (
	"fmt"

	"github.com/mvinturis/mailbox-automation/activity"

	"github.com/chromedp/chromedp"
)

// ActivityBase class extends the Activity struct with common reusable methods
type ActivityBase struct {
	activity.Activity
}

// GetSearchKeyword tests the current search string
func (self *ActivityBase) GetSearchKeyword() (keyword, filter string) {

	chromedp.Run(self.Context,
		chromedp.EvaluateAsDevTools(`$x('//*[@id="filtersButtonId"]/span')[0].innerText`, &filter))
	chromedp.Run(self.Context,
		chromedp.EvaluateAsDevTools(`$x('//*[@id="searchBoxId"]/div/div/input')[0].value`, &keyword),
	)

	fmt.Println("[DEBUG] GetSearchKeyword(): filter = '%s', keyword = '%s'", filter, keyword)
	return
}

// SetSearchKeyword sets the specified keyword to the search box
func (self *ActivityBase) SetSearchKeyword(keyword string, filter string) {
	fmt.Println("[DEBUG] SetSearchKeyword(): keyword = '%s'", keyword)

	localKeyword, localFilter := self.GetSearchKeyword()

	if filter != localFilter {
		chromedp.Run(self.Context,
			// Click Search box
			chromedp.DoubleClick(`#searchBoxId`),
			// chromedp.Click(`#searchBoxId`),
			self.RandomSleep(),
			self.RandomSleep(),

			// Click Filters
			chromedp.Click(`#filtersButtonId > span`), self.RandomSleep(),
			// // Select search Inbox folder
			chromedp.Click(`//div[contains(@id, "Dropdown")][@role="listbox"]`), self.RandomSleep(),
			chromedp.Click(`//button[@title="`+filter+`"][contains(@id, "Dropdown")][@role="option"]`), self.RandomSleep(),
			// // Click Search button
			chromedp.Click(`//div[.="Search"]/ancestor::button`), self.RandomSleep(),
		)
	}
	if keyword != localKeyword {
		chromedp.Run(self.Context,
			// Click Search box
			chromedp.DoubleClick(`#searchBoxId`), self.RandomSleep(),
			chromedp.KeyEvent("\b\b", chromedp.KeyModifiers(0)), self.RandomSleep(),
			// Input search keyword
			chromedp.SendKeys(`//input[@aria-label="Search"][1]`, keyword+"\n"), self.RandomSleep(),
		)
	}
}

func (self *ActivityBase) IsAvailableMailActionByName(name, dual string) bool {
	var value string
	errName := chromedp.Run(self.Context, chromedp.EvaluateAsDevTools(`$x('//button[@name="`+name+`"]')[0].type`, &value))
	errDual := chromedp.Run(self.Context, chromedp.EvaluateAsDevTools(`$x('//button[@name="`+dual+`"]')[0].type`, &value))

	if errName != nil && errDual != nil {
		// Open More mail actions menu
		chromedp.Run(self.Context,
			chromedp.EvaluateAsDevTools(`$x('//*[@aria-label="More mail actions"]')[1].type`, &value),
			chromedp.Click(`(//*[@aria-label="More mail actions"])[1]`, chromedp.NodeVisible), self.RandomSleep(),
		)
	}

	errName = chromedp.Run(self.Context, chromedp.EvaluateAsDevTools(`$x('//*[@name="`+name+`"]')[0].type`, &value))

	if errName != nil {
		fmt.Println("[WARN] IsAvailableMailActionByName(%s, %v): %v", name, dual, false)
		return false
	}
	fmt.Println("[DEBUG] IsAvailableMailActionByName(%s, %v): %v", name, dual, true)
	return true
}

func (self *ActivityBase) SetMailActionByName(name, dual string) {
	if self.IsAvailableMailActionByName(name, dual) {
		chromedp.Run(self.Context,
			chromedp.Click(`(//*[@name="`+name+`"])[1]`, chromedp.NodeVisible), self.RandomSleep(),
		)
	}
}
