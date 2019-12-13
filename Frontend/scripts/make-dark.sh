#!/usr/bin/env bash
echo 'preparing...'
sed -i 's/static\//static\/dark\//g' semantic.json
mv semantic/src/theme.config{,.default}
mv semantic/src/theme.config{.dark,}
echo 'building...'
gulp build-semantic
echo 'rolling back...'
sed -i 's/static\/dark\//static\//g' semantic.json
mv semantic/src/theme.config{,.dark}
mv semantic/src/theme.config{.default,}
