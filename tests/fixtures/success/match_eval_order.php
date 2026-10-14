<?php
function show($x) { echo $x; return $x; } echo "|", match (show(2)) { show(1) => 9, show(2) => 8, default => 7 }, "\n";
