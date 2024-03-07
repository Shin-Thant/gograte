FROM golang:bullseye

WORKDIR /app

RUN go install github.com/Shin-Thant/gograte/cmd/gograte@latest

# 1. /bin/bash - start the interactive shell
# 2. -c - pass the following command to the shell
# 3. cd /app - change the working directory to /app
# 4. /bin/bash - start the interactive shell in the new working directory
CMD ["/bin/bash", "-c", "cd /app && /bin/bash"]

# docker run --name gograte_db -e POSTGRES_PASSWORD=pwd -e POSTGRES_USER=postgres -p 5000:5432 postgres:bullseye