var fill = new ol.style.Fill({ color: 'rgba(255,255,255,0.4)' });
var stroke = new ol.style.Stroke({color: '#5555ff', width: 1.25});
var circleStyles = [];

function getCircle(r) {
	if (circleStyles.length > r) {
		return circleStyles[r];
	}
	circleStyles.length = r + 1;
	for (var i = 1; i <= r; i++) {
		if (circleStyles[i])  {
			continue;
		}
		circleStyles[i] = new ol.style.Style({
			image: new ol.style.Circle({fill: fill, stroke: stroke, radius: i}),
			fill: fill,
			stroke: stroke
		});
	}
	return circleStyles[r]
}

var map
function run() {
	var llsum = [0.0, 0.0]
	for (var i = 0; i < mapFeatures.length; i++) {
		var r = 1
		var n = mapFeatures[i].get("name")
		for (j = 0; j < n.length; j++) {
			if (n[j] == ',') {
				r++
			}
		}
		mapFeatures[i].setStyle(getCircle(r))

		p = mapFeatures[i].getGeometry().getCoordinates()
		lonlat = ol.proj.toLonLat(p)
		llsum[0] += lonlat[0]
		llsum[1] += lonlat[1]
	}

	llsum[0] /= (mapFeatures.length + 0.0)
	llsum[1] /= (mapFeatures.length + 0.0)

	var vectorSource = new ol.source.Vector({ features: mapFeatures })
	var vectorLayer = new ol.layer.Vector({ source: vectorSource });
	map = new ol.Map({
		target: document.getElementById('map'),
		layers: [
			new ol.layer.Tile({source: new ol.source.OSM()}),
			vectorLayer,
		],
		view: new ol.View({
			center: ol.proj.fromLonLat(llsum),
			zoom: 14,
		})
	})
}
