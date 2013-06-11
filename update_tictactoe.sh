#!/bin/bash

BASE_URL=https://raw.github.com/crhym3/go-tictactoe/master

curl "$BASE_URL/static/images/caution.gif" > static/images/caution.gif
curl "$BASE_URL/static/css/base.css" > static/css/tictactoe.css
curl "$BASE_URL/static/js/base.js" > static/js/tictactoe.js
curl "$BASE_URL/templates/tictactoe.html" > templates/tictactoe.html

sed -i .bak -e 's/\/js\/base\.js/\/js\/tictactoe.js/g' templates/tictactoe.html
sed -i .bak -e 's/\/css\/base\.css/\/css\/tictactoe.css/g' templates/tictactoe.html
rm -f templates/tictactoe.html.bak
