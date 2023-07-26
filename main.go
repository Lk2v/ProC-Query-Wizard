package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"gioui.org/font/gofont"
)

const commun_h = "commun.h"

func verif_err(e error) {
	if e != nil {
		panic(e)
	}
}

func nom_fichier(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func main() {

	ui := UI{
		Window: app.NewWindow(app.Title("Query-Wizard")),
		Shaper: text.NewShaper(gofont.Collection()),
		Theme:  NewTheme(gofont.Collection()),

		ResultatReq: "",

		Resize: component.Resize{Ratio: 0.5},
	}
	go func() {
		if err := ui.Loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func traitement(requete string) string {
	var err error

	t1 := time.Now()
	// On récupère le nom du fichier requête

	// RegEx pour trouver les variables
	// On cherche le premier terme qui commence par ":" et qui se termine par un espace
	fmt.Println("Recherche des variables...")

	reg := regexp.MustCompile(":[a-zA-Z0-9_]+")
	ls_var := reg.FindAllString(requete, -1)

	val, val_undef := chercher_val(ls_var)

	fmt.Println("Remplacement des variables par leur valeur...")
	// On remplace les variables par leurs valeurs
	for key, element := range val {
		requete = strings.Replace(requete, key, element, -1)
	}

	// Recuperer des valeurs non définies
	fmt.Println("\nVariables non trouvée(s) :")
	fmt.Println("• :NL -> ' ' (remplacement par un espace)")
	for key, occurences := range val_undef {
		fmt.Println("• " + key + " -> " + strconv.Itoa(occurences) + " occurence(s)")
	}

	/*
		############### STRUCT ###############
		On enleve le prefix de la struct
		Exemple : :struct->nom -> :nom
	*/

	for key, occurences := range val_undef {
		fmt.Println("• " + key + " -> " + strconv.Itoa(occurences) + " occurence(s)")
	}

	reg = regexp.MustCompile(":[\\w]+->")
	requete = reg.ReplaceAllString(requete, ":")

	/*
		###### SELECT ANULLATION DU INTO #####
		On met en commentaire le INTO
	*/

	reg = regexp.MustCompile(`(?Ui)(SELECT(.|\n)*FROM)`)
	for _, match := range reg.FindAllString(requete, -1) {
		intoIndex := strings.Index(match, "INTO")
		intoFrom := strings.Index(match, "FROM")
		if intoIndex != -1 {
			substring := match[intoIndex:intoFrom]

			lignes := strings.Split(substring, "\n")
			for j := range lignes {
				if strings.TrimSpace(lignes[j]) == "" {
					continue
				}
				lignes[j] = "--" + lignes[j]
			}
			requete = strings.Replace(requete, substring, strings.Join(lignes, "\n"), -1)
		}
	}

	/*
		############### ECRITURE ###############
	*/

	fmt.Println("\nEcriture de la requête...")
	now := time.Now()
	// Formater la date
	date := now.Format("02012006")
	// Formater l'heure
	heure := now.Format("150405")
	// Concaténer la date et l'heure avec un tiret "-"
	result := date + "-" + heure

	// On écrit la requête dans un fichier
	err = os.WriteFile(result+".sql", []byte(requete), 0644)
	verif_err(err)

	t2 := time.Now()
	diff := t2.Sub(t1)

	fmt.Println("\n-> Automatisation executée en " + diff.String() + " !")

	return requete
}

func chercher_val(vars []string) (map[string]string, map[string]int) {
	fmt.Println("Recherche des valeurs dans le commun.h...")

	valeurs := make(map[string]string)
	val_undef := make(map[string]int)

	dat, err := os.ReadFile(commun_h)
	verif_err(err)

	contenu := string(dat)
	for i := range vars {
		if vars[i] == ":NL" {
			valeurs[vars[i]] = "' '"
			continue
		}
		// Format : [a-zA-Z0-9_]+[ ]+=+.+; OU .*; pour recuperer toute la ligne;
		reg := regexp.MustCompile(vars[i][1:] + "+[ \\[\\]]+=+.+;")
		ligne := reg.FindString(contenu)

		// Recuperer la valeur
		reg = regexp.MustCompile("\".*\"")
		val := reg.FindString(ligne)

		if val != "" {
			// replace " par '
			reg := regexp.MustCompile("^\\s*\"|\"\\s*$")
			valeurs[vars[i]] = reg.ReplaceAllString(val, "'")
		} else {
			_, ok := val_undef[vars[i]]
			if !ok {
				val_undef[vars[i]] = 1
			} else {
				val_undef[vars[i]] += 1
			}
		}

	}
	return valeurs, val_undef
}

type (
	C = layout.Context
	D = layout.Dimensions
)

// UI specifies the user interface.
type UI struct {
	// External systems.
	// Window provides access to the OS window.
	startButton widget.Clickable
	Window      *app.Window
	// Theme contains semantic style data. Extends `material.Theme`.
	Theme *Theme
	// Shaper cache of registered fonts.
	Shaper *text.Shaper
	// Renderer tranforms raw text containing markdown into richtext.

	Editor         widget.Editor
	ResultatEditor widget.Editor

	component.Resize

	ResultatReq string
}

// Theme contains semantic style data.
type Theme struct {
	// Base theme to extend.
	Base *material.Theme
}

// NewTheme instantiates a theme, extending material theme.
func NewTheme(font []text.FontFace) *Theme {
	return &Theme{
		Base: material.NewTheme(font),
	}
}

// Loop drives the UI until the window is destroyed.
func (ui UI) Loop() error {
	var ops op.Ops

	for {
		e := <-ui.Window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			if ui.startButton.Clicked() {
				ui.ResultatReq = traitement(ui.Editor.Text())
				ui.ResultatEditor.SetText(ui.ResultatReq)
			}

			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		}
	}
}

// Update processes events from the previous frame, updating state accordingly.
func (ui *UI) Update(gtx C) {
	for _, event := range ui.Editor.Events() {
		if _, ok := event.(widget.ChangeEvent); ok {
			// valeur editor change
		}
	}
}

// Layout renders the current frame.
func (ui *UI) Layout(gtx C) D {
	ui.Update(gtx)
	return layout.Flex{
		// Vertical alignment, from top to bottom
		Axis: layout.Vertical,
		// Empty space is left at the start, i.e. at the top
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return ui.Resize.Layout(gtx,
				func(gtx C) D {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
						return material.Editor(ui.Theme.Base, &ui.Editor, "Saisir la requête..").Layout(gtx)
					})
				},
				func(gtx C) D {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
						return material.Editor(ui.Theme.Base, &ui.ResultatEditor, "Resultat").Layout(gtx)
					})
				},
				func(gtx C) D {
					rect := image.Rectangle{
						Max: image.Point{
							X: (gtx.Dp(unit.Dp(4))),
							Y: (gtx.Constraints.Max.Y),
						},
					}
					paint.FillShape(gtx.Ops, color.NRGBA{A: 200}, clip.Rect(rect).Op())
					return D{Size: rect.Max}
				},
			)
		}),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(ui.Theme.Base, &ui.startButton, "Convertir")
				return btn.Layout(gtx)
			},
		),
	)
}
