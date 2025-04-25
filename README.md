# Jawbreaker Game 

## Jawbreaker 2005

I wrote Jawbreaker back in 2005.  The web hosting service went out of business somewhere along the way
(was it Textdrive??) and I lost the code.  However, some gaming website had copied it, and are still hosting. 
The Javascript still works, ancient as it is.  The backend has been lost-in-time, probably was PHP with 
MySQL as the datastore.  

Remember adding some advertising to the page and making $75-$300 per month for close to a year. 
Unfortunately, Google changed advertising algorithms and making significantly less thereafter.
The Javascript is ancient, using the `prototype.js` library (don't remember that), that provides a 
**Ajax** library.

If I were to rewrite this, I'd remove the reliance on Prototype.js, write a new backend (to keep track
of scores), and replace the GIFs with Emojis.  For example, the red gif 
![Jawbreaker 2005 Screenshot](docs/p_red.gif "Jawbreaker") could be replaced with the
red emoji: 🔴(or stylized CSS).

Steve Spencer  
August, 2024

🟣🔵🟢🔴🟡


![Jawbreaker Screenshot](docs/jawbreaker-2005.png "Jawbreaker")

## Jawbreaker 2024

Rewrote Javascript with modern javascript.  No external dependencies required,
first or third party.  No images used, just CSS.  Single file.  The code for
this edition lives in the [2024/](2024/) directory.  Game looks same as 
screenshot below.

## Jawbreaker 2025

Recently (circa April 2025), [Datastar](https://data-star.dev/) has caught my
attention. Reimplemented this new version (basically same html/css) with the 
game logic residing in a Go server.  The Javscript interaction comes from
triggering the Datastar library with `data-*` tags. The code for
this edition lives in the [2025/](2025/) directory.

![Jawbreaker 2025 Screenshot](docs/jawbreaker.png "Jawbreaker")
