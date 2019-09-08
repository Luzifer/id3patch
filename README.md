[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/id3patch)](https://goreportcard.com/report/github.com/Luzifer/id3patch)
![](https://badges.fyi/github/license/Luzifer/id3patch)
![](https://badges.fyi/github/downloads/Luzifer/id3patch)
![](https://badges.fyi/github/latest-release/Luzifer/id3patch)
![](https://knut.in/project-status/id3patch)

# Luzifer / id3patch

`id3patch` is a small CLI wrapper around [bogem/id3v2](https://github.com/bogem/id3v2) to manipulate ID3v2 tag values on the CLI without having to install dependencies.

To install just `go get -u github.com/Luzifer/id3patch` it or download a [release binary](https://github.com/Luzifer/id3patch/releases).

It can be used to quickly check or adjust tags:

```console
# id3patch -f Seelennacht\ -\ \ Die\ Nächtliche\ Stadt.mp3
WARN[0000] No supported ID3v2 tags found                 file="Seelennacht -  Die Nächtliche Stadt.mp3"
INFO[0000] File opened successfully                      album= artist= file="Seelennacht -  Die Nächtliche Stadt.mp3" tag_version=4 title= year=
INFO[0000] No tags changed, no write needed              file="Seelennacht -  Die Nächtliche Stadt.mp3"

# id3patch -f Seelennacht\ -\ \ Die\ Nächtliche\ Stadt.mp3 --artist Seelennacht --album Gaslichtromantik --title 'Die Nächtliche Stadt' --year 2014
WARN[0000] No supported ID3v2 tags found                 file="Seelennacht -  Die Nächtliche Stadt.mp3"
INFO[0000] File opened successfully                      album= artist= file="Seelennacht -  Die Nächtliche Stadt.mp3" tag_version=4 title= year=
INFO[0000] Tags written successfully                     file="Seelennacht -  Die Nächtliche Stadt.mp3"

# id3patch -f Seelennacht\ -\ \ Die\ Nächtliche\ Stadt.mp3
INFO[0000] File opened successfully                      album=Gaslichtromantik artist=Seelennacht file="Seelennacht -  Die Nächtliche Stadt.mp3" tag_version=4 title="Die Nächtliche Stadt" year=2014
INFO[0000] No tags changed, no write needed              file="Seelennacht -  Die Nächtliche Stadt.mp3"
```
