<?php
// For such a stupid task, I felt like writing some PHP.

$search_replace = [
	'fill="#050505"' => 'fill="#FFFFFF"',
	'fill="#030303"' => 'fill="#F5F5F5"',
	'd="M371.611,149.506H353.4v33.801h41' => 'fill="#FFFFFF" d="M371.611,149.506H353.4v33.801h41',
];

$files = glob("static/logos/logo-*.svg");
foreach ($files as $fname) {
	$file = file_get_contents($fname);
	foreach ($search_replace as $s => $r) {
		$file = str_replace($s, $r, $file);
	}
	file_put_contents(substr($fname, 0, -4) . "-dark.svg", $file);
}
