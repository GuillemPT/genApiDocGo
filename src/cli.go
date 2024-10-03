package main

import (
	"fmt"
	"genApiDocGo/src/fileslogic"
	"log"

	"github.com/manifoldco/promptui"
)

func getUniqueSelect(label string, items []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, fileTypeResult, err := prompt.Run()

	if err != nil {
		log.Fatal("Error in selection: ", err)
	}
	return fileTypeResult
}

func getExcludedDirectories(selectedPosition int,
	directories []*fileslogic.Item) ([]string, error) {
	// Always prepend a "Done" item to the slice if it doesn't
	// already exist.
	const doneID = "Done"
	const size = 5
	if len(directories) > 0 && directories[0].ID != doneID {
		var items = []*fileslogic.Item{
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
		Size:      size,
		// Start the cursor at the currently selected index
		CursorPos:    selectedPosition,
		HideSelected: true,
	}

	selectionIndex, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenItem := directories[selectionIndex]

	if chosenItem.ID != doneID {
		// If the user selected something other than "Done",
		// toggle selection on this item and run the function again.
		chosenItem.IsSelected = !chosenItem.IsSelected
		return getExcludedDirectories(selectionIndex, directories)
	}

	// If the user selected the "Done" item, return
	// all selected items.
	var selectedItems []string
	for _, directory := range directories {
		if directory.IsSelected {
			selectedItems = append(selectedItems, "/"+directory.ID+"/")
		}
	}

	return selectedItems, nil
}
