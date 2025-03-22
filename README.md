# process
`docker build -t shift_webapp:dev .`
`docker run -it -p 8080:80 -v ./:/go/project --name shift_webapp shift_webapp:dev`
