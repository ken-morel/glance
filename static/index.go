package handler

import (
  // "context"
  "fmt"
  "github.com/glanceapp/glance/intern/assets"
  "github.com/glanceapp/glance/intern/glance"
  "gopkg.in/yaml.v3"
  "log/slog"
  "net/http"
  // "os"
  "path/filepath"
  "time"
)

var yamlConfig = `server:
  port: 80
theme:
  background-color: 232 23 18
  contrast-multiplier: 1.5
  primary-color: 220 83 75
  positive-color: 105 48 72
  negative-color: 351 74 73
pages:
  - name: Home
    columns:
      - size: small
        widgets:
          - type: clock
            hour-format: 12h
            timezones:
              - timezone: Europe/Paris
                label: Paris
              - timezone: America/New_York
                label: New York
          - type: html
            source: |

              <img src="https://memer-new.vercel.app" id="dev-memes" onclick="alert('dkd')" />
              sdsd
              <script>
                  alert("load");
                  let img = document.querySelector('#dev-memes');
                  let i = 0;
                  function switchMeme() {
                      i++;
                      alert("upload");
                      img.src = "https://memer-new.vercel.app?i=" + i.toString();
                  }
                  setInterval(switchMeme, 1000);
                  setTimeout(switchMeme, 1000);
                  window.switchMeme = switchMeme;
              </script>
          - type: calendar
          - type: monitor
            cache: 1m
            title: Services
            sites:
              - title: API
                url: https://ken-morel.vercel.app
                icon: /static/ama.png
              - title: Sbook
                url: https://sbook.up.railway.app
                icon: /static/antimony.png
          - type: repository
            repository: ken-morel/pyoload
            pull-requests-limit: 5
            issues-limit: 5
      - size: full
        widgets:
          - type: search
            search-engine: google
            bangs:
              - title: Github
                shortcut: "!gh"
                url: https://www.github.com/search?q={QUERY}
              - title: Google
                shortcut: "!go"
                url: https://www.google.com/search?q={QUERY}
          - type: videos
            limit: 15
            channels:
              - UCbKWv2x9t6u8yZoB3KcPtnw
          - type: videos
            limit: 15
            channels:
              - UCS_7tplUgzJG4DhA16re5Yg
          - type: videos
            limit: 15
            channels:
              - UCG9G2dyRv04FDSH1FSYuLBg
          - type: hacker-news
            limit: 15
            collapse-after: 5
          - type: rss
            title: Pypi packages
            style: horizontal-cards
            feeds:
              - url: https://pypi.org/rss/packages.xml
                title: New packages
              - url: https://pypi.org/rss/updates.xml
                title: Updates
              - url: https://pypi.org/rss/project/pyoload/releases.xml
                title: pyoload
      - size: small
        widgets:
          - type: html
            source: |
              <img src="https://ken-morel-stats.vercel.app/api?username=ken-morel&count_private=true&show_icons=true&include_all_commits=true&show=reviews,discussions_started,discussions_answered,prs_merged,prs_merged_percentage&theme=nord&hide_border=true&border_radius=0" />
              <img src="https://ken-morel-api.up.railway.app/counters/github-profile/remove.svg" />
              <img src="https://stackoverflow.com/users/flair/22719308.png?theme=dark&cache=300" />
          - type: html
            source: |
              <img src="https://github-readme-stats.vercel.app/api/pin/?username=ken-morel&repo=pyoload&theme=nord&bg_color=55114455&hide_border=true&border_radius=20" />
          - type: html
            source: |
              <img src="https://wakatime.com/share/embeddable/kenmorel/8508f973-313d-4618-99a9-a023d1e576f2.svg" />
          - type: weather
            units: metric
            hour-format: 12h
            location: Yaounde, Cameroon
          - type: releases
            repositories:
              - glanceapp/glance
              - ken-morel/pyoload
          - type: bookmarks
            groups:
              - links:
                  - title: Gmail
                    url: https://mail.google.com/mail/u/0/
                  - title: stack
                    url: https://www.amazon.com/
                  - title: Github
                    url: https://github.com/
                  - title: Wikipedia
                    url: https://en.wikipedia.org/
              - title: Entertainment
                color: 10 70 50
                links:
                  - title: YouTube
                    url: https://www.youtube.com/
              - title: Social
                color: 200 50 50
                links:
                  - title: Twitter
                    url: https://twitter.com/
`

var serverCreated bool
var mux *http.ServeMux

func GetConfig() (*glance.Config, error) {
  config := glance.NewConfig()

  err := yaml.Unmarshal([]byte(yamlConfig), config)

  if err != nil {
    return nil, err
  }

  // if err = glance.configIsValid(config); err != nil {
  //  return nil, err
  // }

  return config, nil
}

func GetServer() int {
  config, err := GetConfig()

  if err != nil {
    fmt.Printf("failed parsing config file: %v\n", err)
    return 1
  }
  app, err := glance.NewApplication(config)

  if err != nil {
    fmt.Printf("failed creating application: %v\n", err)
    return 1
  }
  mux, err = MuxServer(app)
  if err != nil {
    fmt.Printf("http server error: %v\n", err)
    return 1
  }
  return 0
}
func Handler(w http.ResponseWriter, r *http.Request) {
  if !serverCreated {
    GetServer()
    serverCreated = true
  }
  mux.ServeHTTP(w, r)
}

func MuxServer(a *glance.Application) (*http.ServeMux, error) {
  // TODO: add gzip support, static files must have their gzipped contents cached
  // TODO: add HTTPS support
  mux := http.NewServeMux()

  mux.Handle("GET /static/{path...}", http.StripPrefix("/static/", glance.FileServerWithCache(http.FS(assets.PublicFS), 2*time.Hour)))

  if a.Config.Server.AssetsPath != "" {
    absAssetsPath, err := filepath.Abs(a.Config.Server.AssetsPath)

    if err != nil {
      return mux, fmt.Errorf("invalid assets path: %s", a.Config.Server.AssetsPath)
    }

    slog.Info("Serving assets", "path", absAssetsPath)
    assetsFS := glance.FileServerWithCache(http.Dir(a.Config.Server.AssetsPath), 2*time.Hour)
    mux.Handle("/assets/{path...}", http.StripPrefix("/assets/", assetsFS))
  }
  a.Config.Server.StartedAt = time.Now()

  slog.Info("Created mux server server", "host", a.Config.Server.Host, "port", a.Config.Server.Port)
  return mux, nil
}
