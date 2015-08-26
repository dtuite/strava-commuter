# Strava Commuter

I got sick of manually creating my commute on Strava twice a day, day after day
so I wrote this small CLI program to do it.

Originally I was hoping I could take an existing GPX file, change the dates on
the waypoints of it and upload the file to Strava as if I had just completed the
route described in the file. The application does have that feature but it's
unusable because Strava's activity duplication detection is better than I
expected.

Eventually I realized I'd have to settle for simply creating a manual activity
via the CLI.

It's written in Go and the code is probably horrible because I'm just learning
the language.

## Usage

Create a `~/.strava-commuter/config.yml` config file which contains your Strava
access token (more on that below) and some other details.

```
access_token: hiuhf98hfchu893u89j8dj8jd832
bike_gear_id:
default_activity_description: "Created with Strava Commuer."
default_activity_duration: 600
default_activity_distance: 2900
default_activity_is_private: false
default_activity_is_commute: true
```

You can copy the config.yml.dist in this repository for concenience.

Create commutes like this:

    commuter --finish-time="10:20"

## Strava Setup

Create a developer application at [Strava Developers](https://www.strava.com/developers).
Give it any name, website and description you like and set the Authorization Callback Domain
to `http://localhost:3000`.

You should now have a Client ID, Client Secret and Access Token from Strava.

Open a browser and visit the following URL. Replace `CLIENT_ID` with your Strava
application's Client ID.

```
https://www.strava.com/oauth/authorize?client_id=[CLIENT_ID]&response_type=code&redirect_uri=http://localhost:3000&scope=write
```

Authorize the application on the page that load. Eventually, the browser should
end up on a blank page. Look at the URL in the browser. You should see that the
end of the URL looks something like this: `code=[RANDOM_LETTERS_AND_NUMBERS]`.

Copy the random letters and numbers to the clipboard. Open a terminal and paste
in the following code, once again replacing the values in the square brackets.

```
curl -X POST https://www.strava.com/oauth/token \
     -F client_id=[CLIENT_ID] \
     -F client_secret=[CLIENT_SECRET] \
     -F code=[RANDOM_LETTERS_AND_NUMBERS]
```

Strava should return a blob of stuff. One of the things it returns is an `access_token`
which we will refer to as the `[ACCESS_TOKEN]` from now on. Copy this to the
clipboard.

Open a file called `config.yml` and put the following into it.

```
access_token: [ACCESS_TOKEN]
```

## Activity Setup

TODO: Talk through creating the activity YAML file.
