<?php
$x = 'kept';
$nan = (1e308 * 10) - (1e308 * 10);
if (true) { echo "first|"; }
elseif ($undefined) { $x = 'bad'; echo "bad"; }
else if (1 / 0) { $x = 'bad'; echo "bad"; }
elseif ($nan) { $x = 'bad'; echo "bad"; }
else { echo $missing; $x = 1 / 0; }
if (false) { echo $missing; $x = 1 / 0; }
elseif (1 == '1') { $x = $x . '|middle'; }
else if ($missing) { $x = 'bad'; }
else { if ($nan) { echo "bad"; } $x = 1 / 0; }
echo $x, "\n";
