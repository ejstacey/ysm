
# Joyrex YSM

## Standard Info

![Janky YSM Logo](https://repo.joyrex.net/ejstacey/ysm/raw/branch/main/assets/120px-logo.png "YSM")

Joyrex YSM is a TUI-based application (using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for UI elements) used for keeping a local index of your Youtube Subscriptions and letting you attach notes and tags to them. The idea is this is a backend application that prepares the database so a frontend application can handle displaying the subscription list for a specific tag/whatever, however there is a basic html page generator included as well.

This was created as a way for me to learn/improve writing stuff in Go, but quickly turned into me learning how to use Bubble Tea, Lip Gloss and [The Elm Architecture](https://guide.elm-lang.org/architecture/). This is all to say that while the code is functional, there's very likely cleaner ways to write it. It's a great learning experience and I'm proud of it, though.

You can see what the built-in generator output looks like [here](https://joyrex.net/ysm/).

### Downloading

Official releases are on the [releases page](https://repo.joyrex.net/ejstacey/ysm/releases).

### Compiling

You can build this yourself by doing the following:

    git clone https://repo.joyrex.net/ejstacey/ysm.git
    cd ysm
    go build .

### Running

You need to create a settings.json before you start. You can modify it if you want, but the defaults are fine.

    1. cp settings.sample.json settings.json
    2. ./ysm
    3. log into google and let it sync your subscriptions from youtube
    4. hit 't' to go to the tags list menu, then 'n' to create a new tag.
    5. fill in the tag info and submit it.
    6. when you're done adding tags, hit 'c' to go to the channel list menu
    7. select a channel and hit '<enter>' to go to the channel modify screen
    8. add notes if required. on the tag select, use the arrow keys to move around the tags, then hit '<space>' to select the tag
    9. submit when done/ready
    10. when you've done all the channel noting/tagging you want, from the channel or tag list screen, hit 'g' to go to the generate pagepage
    11. set/change info as wanted, the defaults are probably fine, but you can customise it as you get more used to it.
    12. press 'q' on the channel/tag list screen to quit.
    13. look at the html/ directory to see the page.
    14. upload the contents of the html/ directory onto a web host somewhere.
    15. send the link to your friends and enemies so they can see what you like on youtube.

### Help

Help is at the bottom of most screens. If you hit a bug, email me at [ejstacey@joyrex.net](mailto:ejstacey@joyrex.net). I'm going to get this mirrored to github for bug lodging purposes soon.

## FAQ

1. This is too much work, why didn't you do a webpage to do all this so I can edit/display in my browser?  
I really, really, really don't want your data (and you shouldn't want to hand it over to me). I think there's value in people learning to create/edit webpages, and doing something basic like having the majority filled out for you so you can learn by doing tweaks to something that is there for the most part. There's movements out there to embrace/return to what is called the [small web](https://benhoyt.com/writings/the-small-web-is-beautiful/). I don't follow this perfectly (using javascript), but self-hosting is a bit part of it, and this does that.
2. So where can I host this?
ISPs used to give webspace, yours may still do it. Otherwise I've heard good things about [neocities](https://neocities.org/).
3. Why is this so ugly?
I am really not a UX guy, which means this was a ridiculous task for me to pick up, since there's both a web frontend and the program frontend. I have plans to let people define colours in settings.json so they can make it look nicer.
4. I was looking at your code and there's some sloppiness in there
This project was a year of on and off work, and I was learning as I went, so while the logic works (I think), the code might not be as clean or efficient as it could be. I am going back and still cleaning stuff up as I add more polish to the program.

## Tech stuff

### Todo

In no particular order:

- ~~Add help entries on screens that don't have them~~
- Create better handling of non-fatal errors
- Create confirmation screens
- Abstract the program themeing out to let people have custom themes
- ~~Mirror to github for ticket lodging~~
- Create screenshots for the README
- Create a youtube video showing a typical process flow.
- Squash bugs

### Database format

The database is a simple sqlite database with three tables (diagram made with [draw.io](https://www.drawio.com/)):

![DB ERD for YSM](https://repo.joyrex.net/ejstacey/ysm/raw/branch/main/assets/ysm-erd.png "YSM")

### Youtube Access

This program uses Google's OAuth to retrieve your subscription list. The only permission it uses is "youtube.YoutubeReadonlyScope" so it can grab your subscriptions. All data it grabs is stored locally on your machine. Nothing is sent to me/stored on my side/etc. The code for authenticating to youtube is [here](https://repo.joyrex.net/ejstacey/ysm/src/branch/main/utils/youtube.go) and the code for grabbing the subscriptions is [here](https://repo.joyrex.net/ejstacey/ysm/src/branch/main/channel/channel.go#L29).

Youtube subscription info is retrieved in two situations:

- The channel database is empty
- You set "Refresh" to "true" in your settings.json file.

In these situations the program will give you a browser link to auth with google and try to open that link in your normal browser using xdg-open. Once you auth with google, the time-limited authentication credential is stored in ~/.credentials/ysm.json.

## Privacy Policy and Terms of Service

Please see our [privacy policy](https://repo.joyrex.net/ejstacey/ysm/raw/branch/main/PRIVACY-POLICY) and [terms of service](https://repo.joyrex.net/ejstacey/ysm/raw/branch/main/TERMS-OF-SERVICE) in the root of this repository for the respective information.
