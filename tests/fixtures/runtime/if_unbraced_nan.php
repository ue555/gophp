<?php
$inf = 1e308 * 10; $nan = $inf - $inf;
if (false) echo 'bad'; elseif ($nan) echo 'nan'; else echo 'bad';
echo "\n";
