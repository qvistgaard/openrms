package telemetry

import (
	"embed"
	"html/template"
	"os"
	"time"
)

//go:embed report.gohtml
var reportTemplate embed.FS

func report(race Race) error {
	if len(race) > 0 {
		filename := time.Now().Format("2006-01-02T15-04-05.html")

		t, err := template.ParseFS(reportTemplate, "report.gohtml")
		if err != nil {
			return err
		}

		file, err := os.Create(filename)

		if err != nil {
			return err
		}
		defer file.Close()

		err = t.Execute(file, &race)
		if err != nil {
			return err
		}
	}
	return nil
}
