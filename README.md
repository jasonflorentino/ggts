# 🚆 GGTS 

Go GO Train Schedule: A web app for GO Train schedules written in Go lang.

_"Go, go gadget GO train schedule!"_

It's some time after 5pm and I've just finished firing off a final Slack message or a quick _"LGTM!"_ before stuffing my laptop in my bag, putting my mug in the dishwasher, and hitting the elevator button. I know the final express train is coming soon. But when, exactly? Would I make it in time? Or would I be stuck waiting another thirty minutes for the next train, which will also take twenty minutes longer.

The GO Transit website is loading up on my phone. Which tab do I want? It wants me to _search_ for the stations? Why don't they just lead with a list of "Upcoming Departures" that clearly show the terminus? It's not like this is the [second-busiest railway station in North America](https://en.wikipedia.org/wiki/List_of_busiest_railway_stations_in_North_America). I'm on the Schedules page. It's theoretically straightforward. _Theoretically_. The first item on the list is the trip from 6am this morning. Others on the list include multiple transfers and take 1.5x longer than a direct trip. Of course there's filter, but I've learned that messing with the wrong input on this site might crater the responsiveness of this brutal UI. _Ding! Ding!_ A cyclist is coming as I'm looking to cross. I hit the bottom of the list. Ah, a 12pm bus+train excusion and a button that says _Load More_. 

I'll need to hit that button twice before seeing the part of the list I'm looking for. It's either this or the _Plan Your Trip_ page that'll only show me 3 options. Don't bet on them being the 3 that I want.

All this just to find an excuse to try out [HTMX](https://htmx.org/).

— Jason (Nov 15, 2024)

# Development
You will need:
- The Go programming language: https://go.dev/doc/install
- Air for live-reloading during dev: https://github.com/air-verse/air?tab=readme-ov-file#installation
- The Tailwind CLI for building the stylesheet: https://tailwindcss.com/blog/standalone-cli

To run the app locally:
- Create .env file
```
cp .env.template .env
```
- Run tailwind in watch mode:
```
./tailwindcss -i ./css/style.css -o ./static/style.css --watch
```
- Run air
```
air
```

- The server will be running on the port defined in your .env:
```
open http://locahost:5400
```

# Deployment
I'm running this out of a cheap linux box that already has deps installed.
To deploy the latest changes, ssh in there, pull main, and run `./update.sh`
