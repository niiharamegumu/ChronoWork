package widgets

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/db"
	"github.com/niiharamegumu/ChronoWork/models"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/niiharamegumu/ChronoWork/util/strutil"
	"github.com/rivo/tview"
)

var (
	projectHeader = []string{
		"ID",
		"Name",
		"Tags",
	}
)

type Project struct {
	Layout       *tview.Grid
	Form         *tview.Form
	ReadOnlyForm *tview.Form
	Table        *tview.Table
}

func NewProject() *Project {
	return &Project{
		Layout: tview.NewGrid().
			SetRows(0, 0).
			SetColumns(0, 0).
			SetBorders(true),
		Form: tview.NewForm().
			SetButtonBackgroundColor(tcell.ColorPurple).
			SetLabelColor(tcell.ColorPurple).
			SetFieldTextColor(tcell.ColorGray).
			SetFieldBackgroundColor(tcell.ColorWhite),
		ReadOnlyForm: tview.NewForm().
			SetButtonBackgroundColor(tcell.ColorPurple).
			SetLabelColor(tcell.ColorPurple).
			SetFieldTextColor(tcell.ColorGray).
			SetFieldBackgroundColor(tcell.ColorWhite),
		Table: tview.NewTable().
			SetSelectable(true, false).
			SetFixed(1, 1),
	}
}

func (p *Project) GenerateInitProject(tui *service.TUI) *Project {
	p.setStoreProjectForm(tui)
	p.RestoreTable()

	p.Layout.AddItem(p.Form, 0, 0, 1, 1, 0, 0, false)
	p.Layout.AddItem(p.ReadOnlyForm, 0, 1, 1, 1, 0, 0, false)
	p.Layout.AddItem(p.Table, 1, 0, 1, 2, 0, 0, true)

	p.tableCapture(tui)
	p.formCapture(tui)
	return p
}

func (p *Project) RestoreTable() {
	p.Table.Clear()
	p.setTable()
}

func (p *Project) setStoreProjectForm(tui *service.TUI) {
	p.Form.Clear(true)
	p.ReadOnlyForm.Clear(true)

	p.ReadOnlyForm.AddTextArea("Selected Tags", "", 50, 5, 0, nil)
	tags := models.AllTagNames(db.DB)
	tags = append([]string{notSelectText}, tags...)
	p.Form.AddInputField("Project Name : ", "", 50, nil, nil).
		AddDropDown("Tags : ", tags, 0, func(option string, optionIndex int) {
			link := p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea)
			if option == notSelectText {
				link.SetText("", false)
				return
			}
			if link.GetText() == "" {
				link.SetText(option, false)
				return
			}
			linkTagNames := strings.Split(link.GetText(), ",")
			linkTagNames = append(linkTagNames, option)
			linkTagNames = strutil.RemoveDuplicates(linkTagNames)
			link.SetText(strings.Join(linkTagNames, ","), false)
		}).
		AddButton("Save", func() {
			p.storeProject()
			tui.SetFocus("projectTable")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("projectTable")
		})
}

func (p *Project) setUpdateProjectForm(tui *service.TUI, project *models.ProjectType) {
	p.Form.Clear(true)
	p.ReadOnlyForm.Clear(true)

	p.ReadOnlyForm.AddTextArea("Selected Tags", "", 50, 5, 0, nil)
	tags := models.AllTagNames(db.DB)
	tags = append([]string{notSelectText}, tags...)
	p.Form.AddInputField("Project Name : ", project.Name, 50, nil, nil).
		AddDropDown("Tags : ", tags, 0, func(option string, optionIndex int) {
			link := p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea)
			if option == notSelectText {
				link.SetText("", false)
				return
			}
			if link.GetText() == "" {
				link.SetText(option, false)
				return
			}
			linkTagNames := strings.Split(link.GetText(), ",")
			linkTagNames = append(linkTagNames, option)
			linkTagNames = strutil.RemoveDuplicates(linkTagNames)
			link.SetText(strings.Join(linkTagNames, ","), false)
		}).
		AddButton("Update", func() {
			p.updateProject(project)
			tui.SetFocus("projectTable")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("projectTable")
		})
	p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea).SetText(strings.Join(project.GetTagNames(), ","), false)
}

func (p *Project) formCapture(tui *service.TUI) {
	p.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlB:
			tui.SetFocus("projectTable")
		}
		return event
	})
}

func (p *Project) storeProject() error {
	projectName := p.Form.GetFormItemByLabel("Project Name : ").(*tview.InputField).GetText()
	projectTags := p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea).GetText()
	if projectName == "" {
		return nil
	}
	project := models.ProjectType{
		Name: projectName,
	}
	if projectTags != "" {
		tags := models.TagsByNames(db.DB, strings.Split(projectTags, ","))
		project.Tags = tags
	}
	result := db.DB.Create(&project)
	if result.Error != nil {
		return result.Error
	}
	p.RestoreTable()

	return nil
}

func (p *Project) updateProject(project *models.ProjectType) {
	// TODO: It is necessary to determine the conditions under which it can be updated because it will affect the tasks connected to it.
	projectName := p.Form.GetFormItemByLabel("Project Name : ").(*tview.InputField).GetText()
	projectTags := p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea).GetText()
	if projectName == "" {
		return
	}
	db.DB.Model(&project).Association("Tags").Clear()
	project.Name = projectName
	if projectTags != "" {
		tags := models.TagsByNames(db.DB, strings.Split(projectTags, ","))
		project.Tags = tags
	}
	result := db.DB.Save(&project)
	if result.Error != nil {
		return
	}
	p.RestoreTable()
}

func (p *Project) setTable() {
	p.setTableHeader()
	p.setTableBody()
}

func (p *Project) setTableHeader() {
	for i, header := range projectHeader {
		tableCell := tview.NewTableCell(header).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorPurple).
			SetSelectable(false).
			SetExpansion(1)
		p.Table.SetCell(0, i, tableCell)
	}
}

func (p *Project) setTableBody() {
	projects := models.AllProjectTypeWithTags(db.DB)
	for i, project := range projects {
		p.Table.SetCell(i+1, 0,
			tview.NewTableCell(fmt.Sprint(project.ID)).
				SetAlign(tview.AlignCenter),
		)
		p.Table.SetCell(i+1, 1,
			tview.NewTableCell(project.Name).
				SetAlign(tview.AlignCenter),
		)
		tags := strings.Join(project.GetTagNames(), ",")
		p.Table.SetCell(i+1, 2,
			tview.NewTableCell(fmt.Sprint(tags)).
				SetAlign(tview.AlignCenter),
		)
	}
}

func (p *Project) tableCapture(tui *service.TUI) {
	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'a':
				// store project
				p.setStoreProjectForm(tui)
				tui.SetFocus("projectForm")
			case 'u':
				// update project
				row, _ := p.Table.GetSelection()
				cell := p.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				id := cell.Text
				if intId, err := strconv.ParseUint(id, 10, 0); err == nil {
					unitId := uint(intId)
					project, err := models.FindProjectTypeByID(db.DB, unitId)
					if err != nil {
						log.Println(err)
						break
					}
					p.setUpdateProjectForm(tui, project)
					tui.SetFocus("projectForm")
				}
			case 'd':
				// delete project
				row, _ := p.Table.GetSelection()
				cell := p.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				id := cell.Text

				var projectType *models.ProjectType
				var chronoWorks []models.ChronoWork
				var err error
				var intId uint64
				isExist := false

				intId, _ = strconv.ParseUint(id, 10, 0)
				uintId := uint(intId)
				projectType, _ = models.FindProjectTypeByID(db.DB, uintId)
				chronoWorks, err = models.FindChronoWorksByProjectTypeID(db.DB, projectType.ID)
				if err != nil {
					break
				}
				if len(chronoWorks) > 0 {
					isExist = true
				}
				var modal *tview.Modal
				if isExist {
					modal = tview.NewModal().
						SetText("Can't delete this project. Exist work that use this project.").
						AddButtons([]string{"Close"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							tui.DeleteModal()
							tui.SetFocus("projectTable")
							p.Table.ScrollToBeginning().Select(row, 0)
						})
				} else {
					modal = tview.NewModal().
						SetText("Are you sure you want to delete this project?").
						AddButtons([]string{"Yes", "No"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							if buttonLabel == "Yes" {
								if err := projectType.DeleteProjectType(db.DB); err != nil {
									log.Println(err)
								}
								p.RestoreTable()
							}
							tui.DeleteModal()
							tui.SetFocus("projectTable")
							p.Table.ScrollToBeginning().Select(1, 0)
						})
				}
				tui.SetModal(modal)
				tui.SetFocus("modal")
			}
		}
		return event
	})
}
