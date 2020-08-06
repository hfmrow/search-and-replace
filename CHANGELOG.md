# Search and Replace ©2018-20 H.F.M

## Informations

At the bottom you can find a compiled standalone ".7z" version with its checksum. The ".tar.gz" & ".zip" sources contain a "vendor" directory ensuring you can always compile it even if the official libraries have been changed.

## Changelog

All notable changes to this project will be documented in this file.

#### [1.8] 2020-08-06

#### Added

- Text preview window: context menu allow you to open the directory that contain the current viewed file or opening it using your default file viewer. It's configured to work with most of desktop environement "*xdg-open*", you may change that inside ".opt" file to fit your's if that doesn't work.

- Files list main window: context menu allow you to open selected file or directory that contain the current selection. It's configured to work with most of desktop environement "*xdg-open*", you may change that inside ".opt" file to fit your's if that doesn't work.

- Syntax highlighting cover multiple languages and styles.

- Options window added (icon at top right).

- Option to ignore files where size does not match lower and higher values defined as byte.

#### Fixed

- Some little issues using preview window while lines number are visible.

- Search using escape characters match correctly when pattern contain regex' reserved characters [*.].

#### Changed

- Switch-control to enabling/disabling file-chooser-button has been changed to control the auto refresh on directory change.

- Some options was moved to new options window.

- Global artwork changed, icons, top image.

- The search notifications (information at the top right) using button, have been replaced by a progress bar.

- Found files window and preview window now keep there place on screen when reopened in the same session.

- Preview window: the function that allow line number display was rewritten. Better handling large text. More precision when placing cursor at row number containing the found pattern, and better global stability.

- Cleaning some useless code.

- Rewriting some others part of code to improve speed while searching in files.

#### [1.7.6] 2019-10-02

#### Added

- Adding lines numbers in preview window.

- Highlight found patterns in preview window.

- Option to show/hide forbidden files in result window, (binary and those < 8 bytes).

- Adding filenames, found count in title of the preview window and in found result window.

- The found result windows, now allow to double-click on sub nodes to show it in preview window and positioning preview to specified line.

#### Fixed

- Some issues concerning directories scanning, Dnd and/or command line handling. Now it  preserve the chosen source files (filechooser act independently) a switch prevent unwanted changes.

#### Changed

- Scanning sub directory now handle depth of directory scanning.

- Lot of parts of the code have been rewritten for more stability and a really must faster processing.

- The interface has been reorganized.

#### [1.7.3] 2019-07-19

#### Added

- Switch-control to enabling/disabling file-chooser-button on command line argument passing.

#### Fixed

- Files display on change extension(s) masks

- Directory handling on command line arguments
