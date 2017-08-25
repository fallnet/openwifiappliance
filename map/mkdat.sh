#!/bin/bash

# expects input from parser mac-gps

echo "var mapFeatures = ["
while read a; do
mac=`echo $a | awk ' { print $3 } '`
lon=`echo $a | awk ' { print $2 } '`
lat=`echo $a | awk ' { print $1 } '`
echo "new ol.Feature({name: '$mac', geometry: new ol.geom.Point(ol.proj.fromLonLat([$lon, $lat]))}),"
done
echo "];"