package ggui

import (
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/leerumler/gengar/ggdb"
)

// readSel reads the currently selected line and returns a string
// containing its contents, without trailing spaces.
func readSel(curView *gocui.View) string {

	_, posY := curView.Cursor()
	selection, _ := curView.Line(posY)
	selection = strings.TrimSpace(selection)

	return selection
}

// readCat reads the currently selected category name and matches it to a ggdb.Category.
func readCat(catView *gocui.View, cats []ggdb.Category) *ggdb.Category {

	// Read the name of the currently selected category.
	curCatName := readSel(catView)

	// Search for a category that matches that name.
	var curCat ggdb.Category
	for _, cat := range cats {
		if curCatName == cat.Name {
			curCat = cat
		}
	}

	// And return it.
	return &curCat
}

// readExp reads the currently selected expansion name and matches it to a ggdb.Expansion.
func readExp(expView *gocui.View, exps []ggdb.Expansion) *ggdb.Expansion {

	// Read the name of the currently selected expansion.
	curExpName := readSel(expView)

	// Search for an expansion that matches that name.
	var curExp ggdb.Expansion
	for _, exp := range exps {
		if curExpName == exp.Name {
			curExp = exp
		}
	}

	// And return it.
	return &curExp
}

// selUp moves the cursor/selection up one line.
func selUp(gooey *gocui.Gui, view *gocui.View) error {

	if view != nil {
		view.MoveCursor(0, -1, false)
	}

	return nil
}

// selDown moves the selected menu item down one line, without moving past the last line.
func selDown(gooey *gocui.Gui, view *gocui.View) error {

	if view != nil {
		view.MoveCursor(0, 1, false)

		// If the cursor moves to an empty line, move it back. :P
		if readSel(view) == "" {
			view.MoveCursor(0, -1, false)
		}
	}

	return nil
}

// resetHighlights changes the SelBgColor of the categories, expansions,
// and phrases views back to their "default" blue.
func resetHighlights(gooey *gocui.Gui) error {

	if catView, err := gooey.View("categories"); err == nil {
		catView.SelBgColor = gocui.ColorBlue
	} else {
		return err
	}

	if expView, err := gooey.View("expansions"); err == nil {
		expView.SelBgColor = gocui.ColorBlue
	} else {
		return err
	}

	if phraseView, err := gooey.View("phrases"); err == nil {
		phraseView.SelBgColor = gocui.ColorBlue
	} else {
		return err
	}

	return nil
}

// focusCat changes the focus to the categories view.
func focusCat(gooey *gocui.Gui, view *gocui.View) error {

	// Focus on the categories view.
	if err := gooey.SetCurrentView("categories"); err != nil {
		return err
	}

	// Reset every view's highlight colors back to blue, then set the
	// set the categories view's highlight color to cyan.
	// So everyone's blue but categories.
	resetHighlights(gooey)
	if catView, err := gooey.View("categories"); err == nil {
		catView.SelBgColor = gocui.ColorCyan
	} else {
		return err
	}

	return nil
}

// focusExp changes the focus to the expansions view.
func focusExp(gooey *gocui.Gui, view *gocui.View) error {

	// Focus on the epxansions view.
	if err := gooey.SetCurrentView("expansions"); err != nil {
		return err
	}

	// Reset every view's highlight colors back to blue, then set the
	// expansions view's highlight color to cyan.
	// So everyone's blue but expansions.
	resetHighlights(gooey)
	if expView, err := gooey.View("expansions"); err == nil {
		expView.SelBgColor = gocui.ColorCyan
	} else {
		return err
	}

	return nil
}

// focusPhrase changes the focus to the phrases view.  It also sets
// the highlight color to Cyan for clarity.
func focusPhrase(gooey *gocui.Gui, view *gocui.View) error {

	// Focus on the phrases view.
	if err := gooey.SetCurrentView("phrases"); err != nil {
		return err
	}

	// Reset every view's highlight colors back to blue, then set the
	// set the phrases view's highlight color to cyan.
	// So everyone's blue but phrases.
	resetHighlights(gooey)
	if phraseView, err := gooey.View("phrases"); err == nil {
		phraseView.SelBgColor = gocui.ColorCyan
	} else {
		return err
	}

	return nil
}

// setKeyBinds is a necessary evil.
func setKeyBinds(gooey *gocui.Gui) error {

	// If the categories view is focused and ↑ is pressed, move the selected menu item up one.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowUp, gocui.ModNone, selUp); err != nil {
		return err
	}

	// If the categories view is focused and ↓ is pressed, move the selected menu item down one.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowDown, gocui.ModNone, selDown); err != nil {
		return err
	}

	// If the categories view is focused and → is pressed, move the focus to the expansions view.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowRight, gocui.ModNone, focusExp); err != nil {
		return err
	}

	// If the categories view is focused and Enter is pressed, move the focus to the expansions view.
	if err := gooey.SetKeybinding("categories", gocui.KeyEnter, gocui.ModNone, focusExp); err != nil {
		return err
	}

	// If the expansions view is focused and ↑ is pressed, move the selected menu item up one.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowUp, gocui.ModNone, selUp); err != nil {
		return err
	}

	// If the expansions view is focused and ↓ is pressed, move the selected menu item down one.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowDown, gocui.ModNone, selDown); err != nil {
		return err
	}

	// If the expansions view is focused and ← is pressed, move the focus to the categories menu.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowLeft, gocui.ModNone, focusCat); err != nil {
		return err
	}

	// If the expansions view is focused and → is pressed, move the focus to the phrases view.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowRight, gocui.ModNone, focusPhrase); err != nil {
		return err
	}

	// If the phrases view is focused and ↑ is pressed, move the selected menu item up one.
	if err := gooey.SetKeybinding("phrases", gocui.KeyArrowUp, gocui.ModNone, selUp); err != nil {
		return err
	}

	// If the phrases view is focused and ↓ is pressed, move the selected menu item down one.
	if err := gooey.SetKeybinding("phrases", gocui.KeyArrowDown, gocui.ModNone, selDown); err != nil {
		return err
	}

	// If the phrases view is focused and ← is pressed, move the focus to the expansions menu.
	if err := gooey.SetKeybinding("phrases", gocui.KeyArrowLeft, gocui.ModNone, focusExp); err != nil {
		return err
	}

	return nil
}
