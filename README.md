# Go GoTrain Schedule
A web app for Go Train schedules written in Go lang

TODO: 
- Add Date selector for dates other than today. maybe just tomorrow selector?
- Choose to allow non-direct trips?
  - but how to layout transfers?

# Development
You need tailwind CLI:
https://tailwindcss.com/blog/standalone-cli

- Create .env file
```
cp .env.template .env
```
- Run air
```
air
```
- Run tailwind watcher:
```
./tailwindcss -i ./css/style.css -o ./static/style.css --watch
```