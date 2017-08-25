<!doctype html>
<html lang="en">
<head>
	<meta http-equiv="Content-Type" content="text/html;charset=ISO-8859-1">
	<link rel="stylesheet" href="ol/css/ol.css" type="text/css">
	<style>
	.map {
	height: 80%;
	width: 90%;
	border-width: 2px;
	border-color: black;
	border-style: solid;
	}
	</style>
	<script src="ol/build/ol.js" type="text/javascript"></script>
	<script src="map.js" type="text/javascript"></script>
<?php
	$mapf = "map.dat.js";
	if (isset($_GET['f'])) {
		$mapf = "map.dat.".$_GET['f'].".js";
	}
	echo "<script src=\"$mapf\" type=\"text/javascript\"></script>";
?>

	<title>WIFI</title>
</head>
<body>
<h2>Wifi Map</h2>
<center><div id="map" class="map"><script type="text/javascript">run();</script></div></center>
</body>
</html>