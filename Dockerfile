# Use an official Golang runtime as a parent image
FROM golang:1.21

# Set the working directory in the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Install pipx
RUN apt-get update && apt-get install -y python3-pip && \
    python3 -m pip install --user pipx && \
    python3 -m pipx ensurepath

# Install yt-dlp using pipx
RUN pipx install yt-dlp

# Create a directory for yt-dlp data and set permissions
RUN mkdir /yt-dlp-data
RUN chmod 777 /yt-dlp-data

# Create a shell script to update yt-dlp daily
RUN echo -e '#!/bin/sh\npipx run yt-dlp --update --update -o "/yt-dlp-data"' > /app/update-ytdlp.sh
RUN chmod +x /app/update-ytdlp.sh

# Set up a cron job to run the update script daily
RUN echo "0 0 * * * /bin/sh /app/update-ytdlp.sh >> /var/log/cron.log 2>&1" | crontab -

# Expose port 8080
EXPOSE 8080

# Run the application and start cron
CMD ["/bin/sh", "-c", "(./main &) && (crond -f -L /dev/stdout & cron -f)"]

# docker build -t music-api .
# docker run -p 8080:8080 --name music-api-container music-api
