# Tsniper
Multi-campus Rutgers course sniper implemented as a Go Discord bot. This project has the ability to monitor multiple Rutgers campuses and seasons simultaneously. This project also contains the directions required for deployment as a daemon service but can be run with nohup, tmux, etc.

## Installation
- Clone this Repository
```
git clone git@github.com:brandonliau/Tsniper.git
```

## Create and Configure Discord Bot
- Developer Portal: https://discord.com/developers/applications
- Create application with `New Application`
- Enable all settings under `Privileged Gateway Intents` in the `Bot` section
    
## Inviting the Bot
- Under `OAuth2`, select the `URL Generator`
- Select `bot` in the `scopes` section
- Select `Administrator` in the `Bot Permissions` section
- Go to the `Generated URL` link and select a server to add the bot

## Configure Project
- Create config file
```
touch config.yml
```
- Copy the following into the config file
```
token:
guild:
boarding:
image:
current_campuses:
    - foo
    - bar
    - baz
current_seasons:
    - foo
    - bar
seasons:
    foo:
        name:
        term:
        year:

```
- Configure the config file
    - `token` - Discord API token
    - `guild` - Guild ID of the server
    - `boarding` - Channel ID to send welcome messages
    - `image` - Bot logo (image file path or an imgur link)
    - `current_campuses` - List of Rutgers campuses to monitor
    - `current_seasons` - List of registration seasons to monitor
    - `seasons` - Season data for each of the provided seasons in `current_seasons`

## Create Daemon Service
- Create a daemon service
```
sudo nano /etc/systemd/system/Tsniper.service
```
- Copy the following into the service file
```
[Unit]
Description=Run Tsniper
After=multi-user.target
[Service]
Type=simple
Restart=always
ExecStart=/root/Tsniper/Tsniper
WorkingDirectory=/root/Tsniper
StandardOutput=append:/root/Tsniper/log.log
StandardError=append:/root/Tsniper/log.log
[Install]
WantedBy=multi-user.target
```
- Paths for `WorkingDirectory`, `StandardOutput`, and `StandardError` may differ depending on your system configuration

## Starting and Stopping the Service
- Reload service configuration
```
sudo systemctl daemon-reload
```
- Enable the service
```
sudo systemctl enable Tsniper.service
```
- Disable the service
```
sudo systemctl disable Tsniper.service
```
- Start the service
```
sudo systemctl start Tsniper.service
```
- Stop the service
```
sudo systemctl stop Tsniper.service
```
- Restart the service
```
sudo systemctl restart Tsniper.service
```
- View service status
```
sudo systemctl status Tsniper.service
```
- For more information on daemons, visit the following
    - https://medium.com/p/f0cc55a42267
    - https://github.com/torfsen/python-systemd-tutorial
