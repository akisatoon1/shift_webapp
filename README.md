# process
`sudo docker build -t shift_webapp:dev .`
`docker run -it -p 8080:80 -v ./code:/go/src/code -v ./data:/go/src/data --name shift_webapp shift_webapp:dev`

