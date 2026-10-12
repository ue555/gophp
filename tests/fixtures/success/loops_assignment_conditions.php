<?php
$i = 3; while ($v = $i--) { echo $v; }
echo ':', $v, ':', $i, "\n";
$a = [1, 2, 3]; for ($i = 0; $i < 3; ++$i) $a[$i] *= 2;
echo ($a === [2, 4, 6]) + 0, "\n";
