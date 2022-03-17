# scutil
screen shot utils

## Install

go install github.com/digitalcircle-com-br/scutil@latest

## How it works:

scutil has the following op modes:

### sc

In this mode, a screenshot will be taken. Optional timeout may be provided throught -to parameter, being to the number of seconds to wait

- find
    This mode will try to find the image named by fname in the screen
- daemon
    In this mode sc will run as a daemon, in the background, allowing the following:
    - ctrl+shift+w: will take a screenshot of the screen at each press
    - ctrl+shift+f: In this mode, every 1 sec scutil will try to find every .png file on its directory in the screen. If found, a notification will be displayed.

