<?php
$i = 0; while (false) { echo 'bad'; $i++; }
echo $i, ':';
while ($i < 1) { echo $i++; } echo ':';
while ($i < 4) echo $i++;
echo ':', $i, "\n";
