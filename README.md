
## Joyrex YSM

![Janky YSM Logo](https://gitea.joyrex.net/ejstacey/ysm/raw/branch/main/assets/120px-logo.png "YSM")

Joyrex YSM is a TUI-based application (using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for UI elements) used for keeping a local index of your Youtube Subscriptions and letting you attach notes and tags to them. The idea is this is a backend application that prepares the database so a frontend application can handle displaying the subscription list for a specific tag/whatever.

This was created as a way for me to learn/improve writing stuff in Go, but quickly turned into me learning how to use Bubble Tea, Lip Gloss and [The Elm Architecture](https://guide.elm-lang.org/architecture/). This is all to say that while the code is functional, there's very likely cleaner ways to write it. It's a great learning experience and I'm proud of it, though.

### Compiling

You can build this yourself by doing the following:

    git clone https://gitea.joyrex.net/ejstacey/ysm.git
    cd ysm
    go build .

### Running

You need to create a settings.json before you 

    cp settings.sample.json settings.json
    ./ysm

### Database format

The database is a simple sqlite database with three tables (diagram made with [draw.io](https://www.drawio.com/)):

![DB ERD for YSM](https://gitea.joyrex.net/ejstacey/ysm/raw/branch/main/assets/ysm-erd.png "YSM")

### Youtube Access

This program uses Google's OAuth to retrieve your subscription list. The only permission it uses is "youtube.YoutubeReadonlyScope" so it can grab your subscriptions. All data it grabs is stored locally on your machine. Nothing is sent to me/stored on my side/etc. The code for authenticating to youtube is [here](https://gitea.joyrex.net/ejstacey/ysm/src/branch/main/utils/youtube.go) and the code for grabbing the subscriptions is [here](https://gitea.joyrex.net/ejstacey/ysm/src/branch/main/channel/channel.go#L29).

Youtube subscription info is retrieved in two situations:
- The channel database is empty
- You set "Refresh" to "true" in your settings.json file.

In these situations the program will give you a browser link to auth with google and try to open that link in your normal browser using xdg-open. Once you auth with google, the time-limited authentication credential is stored in ~/.credentials/ysm.json.

### Privacy Policy and Terms of Service

Please see PRIVACY-POLICY and TERMS-OF-SERVICE in the root of this repository for the respective information.
