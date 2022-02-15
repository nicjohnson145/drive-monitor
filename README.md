# drive-monitor
Crappy binary for monitoring MD1000 health and pushing to pushbullet

## Usage

Designed to be ran through cron:

```
5 * * * * bash -c '/opt/dell/srvadmin/bin/omreport storage pdsik controller=0 | APP_TOKEN=<Pushbullet App Token> USER_TOKEN=<Pushbullet User Token> /root/drive-monitor
```
