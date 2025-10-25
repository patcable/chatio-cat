# chatio-cat

chatio-cat will expire messages in a discord channel you specify. It's a golang-based AWS lambda function.

## Why?

Discord doesn't seem to have settings to expire messages in channels, and I wanted to learn the discord API. 

## Prerequesites

chatio-cat is built to take in an event using AWS eventbridge and run a lambda function. So you'll need an AWS account to run this in. You could probably make this work with other function as a service providers if you fork this and rip the AWS stuff out for something else.

## Setup
Create a discord bot, give it guild install permissions. You'll want to give "View Channels," "Manage Messages," and "Read Message History" permissions. Your install URL should look something like this:

```
https://discord.com/oauth2/authorize?client_id=REMOVED&permissions=74752&integration_type=0&scope=bot
```

Generate a bot token and hang onto it, you'll need it later. Then go to the URL and install the bot onto your server (which are called guilds in the API i guss)

Create a new AWS Lambda function. Set the runtime to Amazon Linux 2. Create an environment variable for the function named `CHATIO_CAT_TOKEN` and use the bot token above.

Create a trigger for the function, use "EventBridge" and set some kind of schedule expresison for how often you want it to run. You'll need to flip over to the EventBridge UI and configure the thing to send a constant JSON blob that looks like:

```
{
    "channel_id": "ID of channel you wish to clean up",
    "message_duration": "how long messages should last in the format of 0d0h0m"
}
```

For example, to delete messages after 24h in channel 1431739945235779604 you'd put:
```
{
    "channel_id": "1431739945235779604",
    "message_duration": "1d0h0m"
}
```

You can make an eventbridge trigger for multiple channels, if you want to point it at multiple channels.

## Help

This one's provided as-is, just like the MIT license states. It will delete messages, so be careful with it.