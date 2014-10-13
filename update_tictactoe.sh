#!/bin/bash

BASE_URL=https://raw.githubusercontent.com/crhym3/go-tictactoe/master

curl -sS "$BASE_URL/app/static/images/caution.gif" > static/images/caution.gif
curl -sS "$BASE_URL/app/static/css/base.css" > static/css/tictactoe.css
curl -sS "$BASE_URL/app/static/js/base.js" > static/js/tictactoe_base.js
curl -sS "$BASE_URL/app/static/js/render.js" > static/js/tictactoe_render.js
curl -sS "$BASE_URL/app/static/tictactoe.html" > templates/tictactoe.html

sed -i .bak -e 's/\/js\/base\.js/\/js\/tictactoe_base.js/g' templates/tictactoe.html
sed -i .bak -e 's/\/js\/render\.js/\/js\/tictactoe_render.js/g' templates/tictactoe.html
sed -i .bak -e 's/\/css\/base\.css/\/css\/tictactoe.css/g' templates/tictactoe.html
rm -f templates/tictactoe.html.bak
