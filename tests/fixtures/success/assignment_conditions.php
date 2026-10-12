<?php
$x = 0;
if ($x = 0) echo 'bad';
elseif (($x = 2) === 2) echo $x;
elseif ($x = 3) echo 'bad';
else $x = 4;
echo $x;
if (true) $x = 5; else $x = 6;
echo $x, "\n";
