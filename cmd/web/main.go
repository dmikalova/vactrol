// Command web serves and runs the Vactrol browser client. The same binary is
// built two ways: compiled to WebAssembly it runs the interactive UI in the
// browser (app.RunWhenOnBrowser blocks there); built natively it serves that
// wasm bundle and the required go-app resources over HTTP.
package main

import (
	"log"
	"net/http"

	"github.com/dmikalova/vactrol/internal/web"
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

func main() {
	app.Route("/", func() app.Composer { return web.NewGame() })
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:            "Vactrol",
		ShortName:       "Vactrol",
		Description:     "Vactrol — a KeyForge-style card game, playable in the browser.",
		BackgroundColor: "#1c1c1b",
		ThemeColor:      "#1c1c1b",
		// Plain CSS served as a static file from web/ — no CDN, no build step. The dev
		// server serves it from disk, so editing web/app.css and refreshing the browser
		// applies changes with the server left running (no restart).
		Styles:     []string{"/web/app.css"},
		RawHeaders: []string{boardScript},
	})

	const addr = ":8000"
	log.Printf("Vactrol web client on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

// boardScript turns a vertical scroll wheel into horizontal scrolling over the
// card strips (convenient when a battleline runs off-screen) and keeps the game
// log pinned to its newest entry.
const boardScript = `<script>
(function () {
  document.addEventListener('wheel', function (e) {
    var strip = e.target && e.target.closest ? e.target.closest('.card-strip') : null;
    if (strip && strip.scrollWidth > strip.clientWidth && e.deltaY !== 0) {
      strip.scrollLeft += e.deltaY;
      e.preventDefault();
    }
  }, { passive: false });
  new MutationObserver(function () {
    var log = document.getElementById('gamelog');
    if (log) { log.scrollTop = log.scrollHeight; }
  }).observe(document.documentElement, { childList: true, subtree: true });
})();
</script>`
