package elements

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/qvistgaard/openrms/internal/tui/style"
)

func Button(caption string, focused bool) string {
	var (
		prefixFocusedBaseStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Foreground(lipgloss.Color("205"))
		prefixBlurredBaseStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Foreground(lipgloss.Color("7"))

		suffixFocusedBaseStyle = prefixFocusedBaseStyle.Copy()
		suffixBlurredBaseStyle = prefixBlurredBaseStyle.Copy()

		captionFocusedBaseStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Foreground(lipgloss.Color("15"))
		captionBlurredBaseStyle = captionFocusedBaseStyle.Copy().Foreground(style.InactiveForegroundColor)

		prefix string
		suffix string
	)

	if focused {
		prefix = prefixFocusedBaseStyle.Render("[")
		suffix = suffixFocusedBaseStyle.Render("]")
		caption = captionFocusedBaseStyle.Render(caption)
	} else {
		prefix = prefixBlurredBaseStyle.Render("[")
		suffix = suffixBlurredBaseStyle.Render("]")
		caption = captionBlurredBaseStyle.Render(caption)
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, prefix, caption, suffix)

}
