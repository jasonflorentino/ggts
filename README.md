# Go GoTrain Schedule
A web app for Go Train schedules written in Go lang

TODO: 
- Selecting element from the FROM drop down should add all query params and refetch the possible destinations, and clean the timetable.
- Existing FROM destination should remain, if it's available.
- Selecting an element from the TO drop down should update the Timetable.

- Add Date selector for dates other than today. maybe just tomorrow selector?

- Add SQLite for caching destination data to avoid hitting gotransit on every request
- 5 or 10 min TTL is probs ok.


# Development
You need tailwind CLI:
https://tailwindcss.com/blog/standalone-cli

- Run air
```
air
```
- Run tailwind watcher:
```
./tailwindcss -i ./css/style.css -o ./static/style.css --watch
```