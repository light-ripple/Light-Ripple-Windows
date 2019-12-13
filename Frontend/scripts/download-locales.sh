#!/bin/bash
rm data/{js-,}locales/templates-*.po
cd data/locales
for i in it es de pl ru fr nl sv "fi" ro ko vi; do
	echo "$i"
	wget -O templates-$i.po --quiet "https://cutebirbs.ripple.moe/export/?path=/$i/Hanayo/"
	cd ../js-locales
	wget -O templates-$i.po --quiet "https://cutebirbs.ripple.moe/export/?path=/$i/HanayoJS/"
	i18next-conv -l $i -s templates-$i.po -t ../../static/locale/$i.json
	cd ../locales
done
