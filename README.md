# LENIC
![image](demo.gif)

Lenic is a bonafide social network where users are able to:
- follow each other
- make posts (private and public) to their feed
- comment on each others posts
- add pictures to their user profile and posts for further personalization
- search for users
- receive notifications related to activity relevant to them
- send and receive DMs

Want to give it a spin inside your work network? Thought wou might... ;P

## Requirements:

You'll need:
- Docker
- The [mailer_sender](https://github.com/Anacardo89/mailer_sender) service

## Setup:



The whole thing is developed in Go with Gorilla and HTML templates, with minimal use of JS.
DB is MySQL and it interacts with [mailer_sender](https://github.com/Anacardo89/mailer_sender) through RabbitMQ to send registration emails
It's to show i can do backend :D
