<?php
$i = 0; do { echo $i++; } while (false); echo ':';
do echo $i++; while ($i < 4);
echo ':', $i, "\n";
