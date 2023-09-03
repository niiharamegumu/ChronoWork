package widgets

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/db"
	"github.com/niiharamegumu/ChronoWork/models"
	"github.com/niiharamegumu/ChronoWork/pkg"
	"github.com/niiharamegumu/ChronoWork/service"
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
	p.SetStoreProjectForm(tui)
	p.RestoreTable()

	p.Layout.AddItem(p.Form, 0, 0, 1, 1, 0, 0, false)
	p.Layout.AddItem(p.ReadOnlyForm, 0, 1, 1, 1, 0, 0, false)
	p.Layout.AddItem(p.Table, 1, 0, 1, 2, 0, 0, true)
	return p
}

func (p *Project) SetStoreProjectForm(tui *service.TUI) {
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
			linkTagNames = pkg.RemoveDuplicates(linkTagNames)
			link.SetText(strings.Join(linkTagNames, ","), false)
		}).
		AddButton("Save", func() {
			p.StoreProject()
			tui.SetFocus("projectTable")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("projectTable")
		})
}

func (p *Project) SetUpdateProjectForm(tui *service.TUI, project *models.ProjectType) {
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
			linkTagNames = pkg.RemoveDuplicates(linkTagNames)
			link.SetText(strings.Join(linkTagNames, ","), false)
		}).
		AddButton("Update", func() {
			p.UpdateProject(project)
			tui.SetFocus("projectTable")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("projectTable")
		})
	p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea).SetText(strings.Join(project.GetTagNames(), ","), false)
}

func (p *Project) FormCapture(tui *service.TUI) {
	p.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlB:
			tui.SetFocus("projectTable")
		}
		return event
	})
}

func (p *Project) StoreProject() error {
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

func (p *Project) UpdateProject(project *models.ProjectType) {
	// TODO: exist work that use this project tags. if exist, can't delete, update.
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

func (p *Project) SetTable() {
	p.SetTableHeader()
	p.SetTableBody()
}

func (p *Project) SetTableHeader() {
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

func (p *Project) SetTableBody() {
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

func (p *Project) RestoreTable() {
	p.Table.Clear()
	p.SetTable()
}

func (p *Project) TableCapture(tui *service.TUI) {
	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'a':
				// store project
				p.SetStoreProjectForm(tui)
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
					p.SetUpdateProjectForm(tui, project)
					tui.SetFocus("projectForm")
				}
			case 'd':
				// delete project
				// TODO: exist work that use this project. if exist, can't delete.
				row, _ := p.Table.GetSelection()
				cell := p.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				modal := tview.NewModal().
					SetText("Are you sure you want to delete this project?").
					AddButtons([]string{"Yes", "No"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "Yes" {
							id := cell.Text
							if intId, err := strconv.ParseUint(id, 10, 0); err == nil {
								uintId := uint(intId)
								projectType, err := models.FindProjectTypeByID(db.DB, uintId)
								if err != nil {
									log.Println(err)
								}
								if err := projectType.DeleteProjectType(db.DB); err != nil {
									log.Println(err)
								}
								p.RestoreTable()
							}
						}
						tui.DeleteModal()
						tui.SetFocus("projectTable")
						p.Table.ScrollToBeginning().Select(1, 0)
					})
				tui.SetModal(modal)
				tui.SetFocus("modal")
			}
		}
		return event
	})
}
