# Strava Commuter

I got sick of uploading my commute to Strava day after day so I wrote this
small CLI program to read activity data from a YAML file (my commute is the
same route and roughly the same time each day) and send to to Strava via
their API.

It's written in Go and the code is probably horrible because I'm just learning it.

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
