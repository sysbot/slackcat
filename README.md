# slackcat
Slackcat is a simple commandline utility to post snippets to Slack.


  <img width="500px" src="https://raw.githubusercontent.com/vektorlab/slackcat/master/demo.gif" alt="slackcat"/>

## Made some changes

## Usage
Pipe command output:
```bash
$ echo -e "hi\nthere" | slackcat --channel general --filename hello
file hello uploaded to general
```

Post an existing file:
```bash
$ slackcat -c general /home/user/bot.png
file bot.png uploaded to general
```

## Installing

Download the latest release for your platform:

#### OS X

```bash
curl -sLo slackcat https://github.com/vektorlab/slackcat/releases/download/v0.7/slackcat-0.7-darwin-amd64
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
```

#### Linux

```bash
wget https://github.com/vektorlab/slackcat/releases/download/v0.7/slackcat-0.7-linux-amd64 -O slackcat
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
```

## Configuration

Generate a new Slack token with:
```bash
slackcat --configure
```
A new browser window will be opened for you to confirm the request via Slack, and you'll be returned a token.

Create a Slackcat config file and you're ready to go!
```bash
echo '<your-slack-token>' > ~/.slackcat
```

## Options

Option | Description
--- | ---
--tee, -t | Print stdin to screen before posting
--noop | Skip posting file to Slack. Useful for testing
--channel, -c | Slack channel or group to post to
--filename, -n | Filename for upload. Defaults to given filename or current timestamp if reading from stdin.
