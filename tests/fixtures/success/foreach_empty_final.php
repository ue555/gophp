<?php
$k = 'keep'; $v = 9;
foreach ([] as $k => $v) { echo 1 / 0; }
echo $k, ':', $v, "\n";
foreach ([7] as $v) echo $v;
echo ':', $v, "\n";
