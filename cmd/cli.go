package main

import (
	"fmt"
	"log"

	"genApiDocGo/internal"

	"github.com/manifoldco/promptui"
)

func getSelectLanguage() string {

	prompt := promptui.Select{
		Label: "Enter type of files to generate the documentation:",
		Items: internal.FILE_TYPE_OPTIONS,
	}
	_, fileTypeResult, err := prompt.Run()

	if err != nil {
		log.Fatal(err)
	}
	return fileTypeResult
}

func getExcludedDirectories(selectedPosition int,
	directories []*item) ([]string, error) {
	// Always prepend a "Done" item to the slice if it doesn't
	// already exist.
	const doneID = "Done"
	if len(directories) > 0 && directories[0].ID != doneID {
		var items = []*item{
			{
				ID: doneID,
			},
		}
		directories = append(items, directories...)
	}

	// Define promptui template
	templates := &promptui.SelectTemplates{
		Label: `{{if .IsSelected}}
                    ✔
                {{end}} {{ .ID | green }} - label`,
		Active:   "→ {{if .IsSelected}}✔ {{end}}{{ .ID | white }}",
		Inactive: "{{if .IsSelected}}✔ {{end}}{{ .ID | white }}",
	}

	prompt := promptui.Select{
		Label:     "Select the directories that want exclude",
		Items:     directories,
		Templates: templates,
		Size:      5,
		// Start the cursor at the currently selected index
		CursorPos:    selectedPosition,
		HideSelected: true,
	}

	selectionIdx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenItem := directories[selectionIdx]

	if chosenItem.ID != doneID {
		// If the user selected something other than "Done",
		// toggle selection on this item and run the function again.
		chosenItem.IsSelected = !chosenItem.IsSelected
		return getExcludedDirectories(selectionIdx, directories)
	}

	// If the user selected the "Done" item, return
	// all selected items.
	var selectedItems []string
	for _, directory := range directories {
		if directory.IsSelected {
			selectedItems = append(selectedItems, directory.ID+"/")
		}
	}

	return selectedItems, nil
}
