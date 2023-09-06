# Use an official Golang runtime as a parent image
FROM golang:1.21

# Set the working directory in the container
WORKDIR /app

# Install necessary packages for Python virtual environment and cron
RUN apt-get update && apt-get install -y python3-venv python3-pip cron

# Create a virtual environment for Python
RUN python3 -m venv /venv

# Activate the virtual environment
ENV PATH="/venv/bin:$PATH"

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Install pipx inside the virtual environment
RUN pip install pipx

# Use pipx to install yt-dlp
RUN pipx install yt-dlp

# Create and configure the script to update yt-dlp
RUN echo -e '#!/bin/sh\npipx upgrade yt-dlp' > /app/update-ytdlp.sh
RUN chmod +x /app/update-ytdlp.sh

# Add a cron job to run your script daily
RUN (crontab -l ; echo "0 0 * * * /app/update-ytdlp.sh") | crontab -

# Expose port 8080
EXPOSE 8080

# Run the application and start cron
CMD ["/bin/sh", "-c", "(./main &) && cron -f"]


# docker build -t music-api .
# docker run -p 8080:8080 --name music-api-container music-api
