<h1 align="center">
  <br>
  <a href="https://www.morphysm.com/"><img src="./assets/morph_logo_rgb.png" alt="Morphysm" ></a>
  <br>
  <h5 align="center"> Morphysm is a community of engineers, designers and researchers
contributing to security, cryptography, cryptocurrency and AI.</h5>
  <br>
</h1>

<h1 align="center">
  <img src="https://img.shields.io/badge/Go-^1.18.0-green" alt="python badge">

  <img src="https://img.shields.io/badge/version-1.1-green" alt="version badge">
  <img src="https://img.shields.io/gitlab/pipeline-status/dicu.chat/server?branch=master" alt="docker build">
  <a href="https://codecov.io/gh/morphysm/famed-github-backend">
    <img src="https://codecov.io/gh/morphysm/famed-github-backend/branch/master/graph/badge.svg?token=P5ZUKZF9XN"/>
  </a>
  <a href="https://lgtm.com/projects/g/morphysm/famed-github-backend/alerts/">
    <img alt="Total alerts" src="https://img.shields.io/lgtm/alerts/g/morphysm/famed-github-backend.svg?logo=lgtm&logoWidth=18"/>
  </a>
  <a href="https://lgtm.com/projects/g/morphysm/famed-github-backend/context:go">
    <img alt="Language grade: Go" src="https://img.shields.io/lgtm/grade/go/g/morphysm/famed-github-backend.svg?logo=lgtm&logoWidth=18"/>
  </a>
</h1>


# Famed-Backend

This repository contains the code of the Famed-Backend.

# Table of Contents

<!--ts-->

- [How to Famed](#how-to-famed)
- [Security Considerations](#security-considerations)
- [Develop](#develop)
  - [Prerequisites](#prerequisites)
  - [Run](#run)
    - [Env Variables](#env-variables)
- [Troubleshooting](#troubleshooting)
- [Code Owners](#code-owners)
- [Contribute](#contribute)
- [License](#license)
- [Contact](#contact)
<!--te-->


# How to Famed

🚧 [New guide in construction](https://github.com/morphysm/famed-github-backend/wiki/Installation-guide-&-first-start-%F0%9F%90%A7) 🚧
1. Install the Famed GitHub App (https://github.com/apps/get-famed) and allow the app to access to your repository.</br>
   ***Note:** We populate the issue labels when you allow the app to access your repository: "famed", "none", "low", "medium", "high", "critical". We do not overwrite your labels if labels with the same name are present.*
2. Setup frontend:
   1. You can find your public board at `https://www.famed.morphysm.com/teams/<owner>/<repoName>`
   2. Use our famed-board react component (work in progress)
   3. Use our famed-board js script (work in progress)
3. Label your repository issues:
   1. Assign a “famed” label to the issues you want to track with Famed
   2. Assign a severity label to each issue tracked by Famed. We follow the Common Vulnerability Scoring System (CVSS). (Low, Medium, High, Critical)
   3. Make sure the issue has an assignee when closing the issue<br><br>
      
   You will see comments by the Famed bot on your issues labeled with "famed" - the frontend is updated once the first issues are closed.


4. Join Famed on Telegram: https://t.me/+iQPfZQNshl04YmIy

# Security Considerations
We memmemory encrypted the GitHub keywith https://github.com/awnumar/memguard to mitigate memmory dump readout attacks.

We use -buildmode=pie resulting in all addresses except the stack being randomized. (https://rain-1.github.io/golang-aslr.html)


# Self Host
Coming Soon

### GitHub App
Coming Soon


# Develop

## Prerequisites

Please make sure that your system has the following programs:

- [go min. v1.17](https://go.dev/doc/install)

1. Create your own GitHub app.
2. Add a webhook secret to your GitHub app.
3. Use a reverse proxy method of your choice to forward requests from github to your localhost port. (e.g. https://ngrok.com/)
4. Add the reverse proxy endpoint for callbacks (famed/webhooks/event) at the GitHub app.
5. Set up the Env variables.

## Run

### Env Variables
🚧 [New env variables list in construction](https://github.com/morphysm/famed-github-backend/wiki/Configuration-and-environment-variables
) 🚧

- GITHUB_API_KEY: Secret key of the Famed GitHub app (GoLand might format your API key wrongly - Go to .idea/workspace.xml with a alternative editor and set  <env name="GITHUB_API_KEY" value=<Key>/> where you replace newlines with &#10;).
- GITHUB_APP_ID: ID of the Famed GitHub app
- GITHUB_BOT_LOGIN: Login Name of the Famed GitHub app bot (GitHub App name - spaces replaced by "-" + [bot] e.g. : get-famed[bot] )
- GITHUB_WEBHOOK_SECRET: Webhook secret key of the Famed GitHub app
- GITHUB_FAMED_LABEL: Label used to assign issues to the Famed Process
- ADMIN_USERNAME: Username for simple auth admin calls
- ADMIN_PASSWORD: Password for simple auth admin calls
- NEWRELIC_ENABLED: Enable New Relic tracing (feature still experimental / in development)
- NEWRELIC_KEY: New Relic authentication key (leave empty if NEWRELIC_ENABLED=false)
- NEWRELIC_NAME: New Relic service name (leave empty if NEWRELIC_ENABLED=false)

# Troubleshooting

If you have encountered any problems while running the code, please open a new issue in this repo and label it bug, and we will assist you in resolving it.

# Code Owners

@morphysm/team 😎

# Contribute

Developers interested in contributing should read the [Contribution Guide](https://github.com/morphysm/famed-github-backend/wiki/Contribution-Guide).

# License

Our repository is licensed under the terms of the [GNU Affero General Public License v3.0](https://github.com/morphysm/famed-github-backend/blob/master/LICENSE).

# Contact

If you'd like to know more about us visit [https://www.morphysm.com/](https://www.morphysm.com/), or contact us at [contact@morphysm.com](mailto:contact@morphysm.com).
