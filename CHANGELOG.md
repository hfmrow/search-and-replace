# Changelog

## Search and Replace ©2018-19 H.F.M

All notable changes to this project will be documented in this file.

## [1.7.6] 2019-10-02

### Added

- Adding lines numbers in preview window.
- Highlight found patterns in preview window.
- Option to show/hide forbidden files in result window, (binary and those < 8 bytes).
- Adding filenames, found count in title of the preview window and in found result window.
- The found result windows, now allow to double-click on sub nodes to show it in preview window and positioning preview to specified line.

### Fixed

- Some issues concerning directories scanning, Dnd and/or command line handling. Now it  preserve the chosen source files (filechooser act independently) a switch prevent unwanted changes.

### Changed

- Scanning sub directory now handle depth of directory scanning.
- Lot of parts of the code have been rewritten for more stability and a really must faster processing.
- The interface has been reorganized.