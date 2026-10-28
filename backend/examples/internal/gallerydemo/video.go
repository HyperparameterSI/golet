// Author: HyperparameterSI

package gallerydemo

import (
	"embed"

	"github.com/HyperparameterSI/golet/backend/ui"
)

//go:embed waves.gif
var videoAssets embed.FS

// videoCard creates the decoder only when the drawing page is visited.
// The gallery owns the player and releases it when its window closes.
func (g *gallery) videoCard(c *ui.Context) {
	card(c, "Video · 按帧播放", func() {
		if g.video == nil {
			data, err := videoAssets.ReadFile("waves.gif")
			if err != nil {
				ui.Text(c, err.Error())
				return
			}
			decoder, err := ui.NewGIFDecoder(data)
			if err != nil {
				ui.Text(c, err.Error())
				return
			}
			player, err := ui.NewVideoPlayer(decoder)
			if err != nil {
				ui.Text(c, err.Error())
				return
			}
			g.video = player.SetLoop(true)
			g.win.OnClosed(player.Close)
		}
		ui.Text(c, "后台解码、循环播放；不可见时停止推进，遵循系统减少动态效果设置。").TextColor(c.Theme().TextMuted)
		ui.Video(c, g.video).Size(480, 270).Fit(ui.Cover).Radius(12)
		label := "暂停"
		if !g.video.Playing() {
			label = "播放"
		}
		if ui.Button(c.Key("video-play"), label).Clicked() {
			if g.video.Playing() {
				g.video.Pause()
			} else {
				g.video.Play()
			}
		}
		if err := g.video.Err(); err != nil {
			ui.Text(c, err.Error()).TextColor(c.Theme().Danger)
		}
	})
}
