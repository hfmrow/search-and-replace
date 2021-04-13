# Search and Replace ©2018-21 H.F.M

## information

At the bottom you can find a compiled standalone ".deb" version with its checksum. The ".tar.gz" sources contain a "vendor" directory ensuring you can always compile it even if the official libraries have been changed.

## Changelog

All notable changes to this project will be documented in this file.

#### [1.10] 2021-04-14

#### Added

- When the files are dropped, the parent directory is now considered as root, it is stored in the file selector, if the button "Scan files" is pressed , this directory will be fully analyzed.
- Adding expand all in found files window.
- Adding check boxes to un-select lines to be preserved during replace actions.
- Adding un-select all and invert selection in found window.

#### Fixed

- The text preview window does not correctly highlight the searched pattern when escaped characters [Esc search] are used. This is fixed and no longer occurs.

#### Changed

- The application can be closed during the search process, of course the operation will be interrupted.
- By clicking on the search button, if there is no selected file, it will be considered that all files are selected (previously, file selection was explicitly required). Of course, the explicit selection of file (s) remains possible.
- Some incompatible search combinations are now disabled depending on the options checked.

---

#### [1.9] 2021-04-02

Since this version: [Golang GtkSourceView binding for use with gotk3](https://github.com/hfmrow/gotk3_gtksource). `$ sudo apt install libgtksourceview-4-dev` is required to be compiled.

If your OS does not implement GtkSourceView library natively, you can install it using:

`$ sudo apt install libgtksourceview-4-0`

If your Debian version is < **focal** (*20.04*), this means like **xenial** (*16.04*), **bionic** (*18.04*), this version does not work, it must re-compiled with appropriate `libgtksourceview-3.0-dev` library and version restrictions commands:

`$ go build -tags "gtksourceview_3_24" source.go`

#### Added

- Improving Source-code displaying [500ms for source > 350Kb]
- Allow the use of 'escape' characters with replace pattern. An option was added to enable or not the use of it.
- Source map is now included in preview window.

#### Fixed

- Application close while adding numbers to text > 360Kb, fixed using [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource)

#### Changed

- Repository name was changed to [https://github.com/hfmrow/search-and-replace](https://github.com/hfmrow/search-and-replace) instead of `https://github.com/hfmrow/sAndReplace`
- A lot of source code has been modified to be able to use the features described below.
- Text numbering provided by [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource) now.
- Preview search is now provided using [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource) library
- Syntax highlighting is now provided by [GtkSourceView](https://github.com/hfmrow/gotk3_gtksource) too
- Code updated to fit actual version of [GitHub - Go bindings for GTK3](https://github.com/gotk3/gotk3)

---

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

- Code updated to fit actual version of [GitHub - Go bindings for GTK3](https://github.com/gotk3/gotk3)

---

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

---

#### [1.7.3] 2019-07-19

#### Added

- Switch-control to enabling/disabling file-chooser-button on command line argument passing.

#### Fixed

- Files display on change extension(s) masks

- Directory handling on command line arguments
