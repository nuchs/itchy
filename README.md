# itchy
Tool for handling multiple windows in the sway scratchpad

You call itchy with the and identifier of an application you want to manage and a
command specifying how to start the application. The behaviour is then as
follows

* If the application is not running, it will be started and its window displayed
in the active workspace
* If the application is running and its window is focused the window will be sent
to the scratchpad.
* If the application is running but not focused it will be focused
* If the application is running and it is in the scratchpad it will be moved to
the active workspace and focused.

## Credit

Idea shamelessly nicked from this reddit post

https://www.reddit.com/r/swaywm/comments/qop9ys/more_than_one_scratchpad_window/

