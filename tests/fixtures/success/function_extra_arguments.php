<?php
function valueOf($x) { return $x; }
function side() { echo 'side'; return 99; }
echo valueOf(7, side()), "\n";
