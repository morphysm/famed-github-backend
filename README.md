<h1 align="center">
  <br>
  <a href="https://www.morphysm.com/"><img src="./assets/morph_logo_rgb.png" alt="Morphysm" ></a>
  <br>
  <h5 align="center"> Morphysm is a community of engineers, designers and researchers
contributing to security, cryptography, cryptocurrency and AI.</h5>
  <br>
</h1>

<h1 align="center">
  <img src="https://img.shields.io/badge/Go-^1.17.0-red" alt="python badge">

 <img src="https://img.shields.io/badge/version-1.1-orange" alt="version badge">
 <img src="https://img.shields.io/gitlab/pipeline-status/dicu.chat/server?branch=master" alt="docker build">
</h1>

# Table of Contents

<!--ts-->

- [Famed-Backend](#famed-backend)
- [How to Famed](#use-to-famed)
- [Develop](#installation)
- [Prerequisites](#prerequisites)
- [Troubleshooting](#troubleshooting)
- [Code Owners](#code-owners)
- [License](#license)
- [Contact](#contact)
<!--te-->

# Famed-Backend

This repository contains the code of the Famed-Backend.

# How to Famed

1. Install the Famed GitHub App (https://github.com/apps/get-famed) and allow it access to your repository.
2. Join Famed on Telegram: https://t.me/+iQPfZQNshl04YmIy
3. Setup frontend:
   1. Let us know  through telegram that you installed the app, and we will set up a frontend for you. All we need is the app owner and the repository name. You can copy them from your repository URL: https://github.com/<owner>/<repoName>
   2. Use our famed-board react component (work in progress)
   3. Use our famed-board js script (work in progress)
4. Setup your repository issues:
   1. Assign a “famed” label to the issues you want to track with Famed 
   2. Assign a severity label to each issue tracked by Famed. We follow the Common Vulnerability Scoring System (CVSS). (Low, Medium, High, Critical	)
   3. Make sure the issue has an assignee when closing the issue

If all is set up correctly, you will see comments by the Famed bot on your closed issues, and the frontend should be updated accordingly.

# Develop

## Prerequisites

Please make sure that your system has the following programs:
- [go min. v1.17](https://go.dev/doc/install)


## Run
### Env Variables
- GITHUB_API_KEY: Secret key of the Famed GitHub app
- GITHUB_APP_ID: ID of the Famed GitHub app
- GITHUB_BOT_ID: ID of the Famed GitHub app bot
- GITHUB_WEBHOOK_SECRET: Webhook secret key of the Famed GitHub app
- GITHUB_FAMED_LABEL: Label used to assign issues to the Famed Process


# Troubleshooting

If you have encountered any problems while running the code, please open a new issue in this repo and label it bug, and we will assist you in resolving it.

# Code Owners

@morphysm/team :sunglasses:

# License:

Our repository is licensed under the terms of the GNU Affero General Public License v3.0 see the license [here!](https://gitlab.com/dicu.chat/rasa/-/blob/master/LICENSE)

# Contact

If you'd like to know more about us visit [https://www.morphysm.com/](https://www.morphysm.com/), or contact us at [contact@morphysm.com](mailto:contact@morphysm.com).
