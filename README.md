# 🚆 GGTS 

Go GO Train Schedule: A web app for GO Train schedules written in Go lang.

> _"Go, go gadget GO train schedule!"_

It's some time after 5pm and I've just finished firing off a quick _"LGTM!"_ before stuffing my laptop in my bag, putting my mug in the dishwasher, and jabbing the elevator button. I know the final express train is coming soon. But when, exactly? Would I make it in time? Or would I be stuck waiting another thirty minutes for the next train, which also travels twenty minutes longer.

The GO Transit website is loading on my phone.  Which tab do I want? It wants me to _search_ for the stations? Why don't they lead with "Upcoming Union Station Departures" and clearly state the terminus?<sup>1</sup> It's not like this is the [second-busiest railway station in North America](https://en.wikipedia.org/wiki/List_of_busiest_railway_stations_in_North_America). I don't think a jog is necessary – yet. I'm on the Schedules page. It's theoretically straightforward. Theoretically. The first line is the trip from 6am this morning. Other items on the list feature multiple transfers and take 1.5x longer than a direct trip. Of course there's a filter, but I've learned that touching the wrong input on this site might crater the responsiveness of an already brutal UI. _Ding! Ding!_ A cyclist is coming as I'm looking to cross the street. I hit the bottom of the list. It's a 12pm bus+train excursion and a button that reads _Load More_. 

I'll need to hit that button twice before seeing the part of the list I'm looking for. It's this or the _Plan Your Trip_ page that'll only show me 3 options. Don't bet on them being the 3 that I want. I guess it's that time again – time to spend my weekend making a thing!

All this stress just to have an excuse to try out [HTMX](https://htmx.org/).

— Jason (Nov 15, 2024)

<sup>1</sup> Not every train goes on a line goes to the same stops. This is true for the terminal stop too. So to find out if a train is going _that far_ you have to wait around watching the a list of stops cycle through to the end, 3 stops at a time.

# Development
You will need:
- The Go programming language: https://go.dev/doc/install
- Air for live-reloading during dev: https://github.com/air-verse/air?tab=readme-ov-file#installation
- The Tailwind CLI for building the stylesheet: https://tailwindcss.com/blog/standalone-cli

To run the app locally:
- Create a `.env` file and fill in your values.
```
cp .env.template .env
```
- Run `tailwindcss` in watch mode.
```
./tailwindcss -i ./css/style.css -o ./static/style.css --watch
```
- Run `air`.
```
air
```

- The server will be running on the port defined in your `.env`.
```
open http://locahost:5400
```

# Deployment
I'm running this out of a cheap linux box that already has deps installed.
To deploy the latest changes, ssh in there, pull main, and run `./update.sh`
