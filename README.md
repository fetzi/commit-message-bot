# commit-message-bot

[![Build Status](https://img.shields.io/travis/karriereat/commit-message-bot.svg?style=flat-square)](https://travis-ci.org/karriereat/commit-message-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/karriereat/commit-message-bot?style=flat-square)](https://goreportcard.com/report/github.com/karriereat/commit-message-bot)
[![license](https://img.shields.io/badge/license-Apache%202.0-brightgreen.svg?style=flat-square)](https://github.com/karriereat/commit-message-bot/blob/master/LICENSE)

A bot that validates gitlab commit messages and informs the commit author about the broken commit message rule(s).

## Installation

- copy the `sample.toml` from `/conf` and edit it to your needs
- run `commit-message-bot`
- Point the gitlab webhook to `http://your.domain/hooks/gitlab`
