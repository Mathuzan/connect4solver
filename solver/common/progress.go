package common

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

const ProgressBarResolution = 1_000_000_000

func NewProgressBar() *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		ProgressBarResolution,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
	bar.RenderBlank()
	return bar
}
