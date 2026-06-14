package main

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func showAboutDialog(w fyne.Window) {
	projectURL, _ := url.Parse(appProjectURL)
	projectLink := widget.NewHyperlink(appProjectURL, projectURL)

	checkUpdatesButton := widget.NewButton(tr(msgCheckUpdates), func() {
		checkUpdates(w)
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle(tr(msgWindowTitle), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("%s: %s", tr(msgVersion), appVersion)),
		widget.NewLabel(appCopyright),
		widget.NewLabel(tr(msgProjectLink)+":"),
		projectLink,
		widget.NewSeparator(),
		checkUpdatesButton,
	)

	dialog.ShowCustom(tr(msgAbout), tr(msgClose), content, w)
}

func checkUpdates(w fyne.Window) {
	info, err := CheckForUpdates()
	if err != nil {
		dialog.ShowError(fmt.Errorf("%s: %w", tr(msgUpdateCheckFailed), err), w)
		return
	}

	if !info.IsNewer {
		dialog.ShowInformation(
			tr(msgNoUpdatesTitle),
			fmt.Sprintf("%s\n%s: %s", tr(msgNoUpdatesBody), tr(msgVersion), info.CurrentVersion),
			w,
		)
		return
	}

	message := fmt.Sprintf(
		"%s\n\n%s: %s\n%s: %s",
		tr(msgUpdateAvailableBody),
		tr(msgCurrentVersion),
		info.CurrentVersion,
		tr(msgLatestVersion),
		info.LatestVersion,
	)
	if info.ReleaseName != "" {
		message += "\n" + info.ReleaseName
	}

	confirm := dialog.NewConfirm(tr(msgUpdateAvailableTitle), message, func(download bool) {
		if !download {
			return
		}
		releaseURL, err := url.Parse(info.ReleaseURL)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err := fyne.CurrentApp().OpenURL(releaseURL); err != nil {
			dialog.ShowError(err, w)
		}
	}, w)
	confirm.SetConfirmText(tr(msgDownload))
	confirm.SetDismissText(tr(msgClose))
	confirm.Show()
}
