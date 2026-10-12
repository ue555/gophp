<?php
$a = [[1, 2], [3]]; $b = $a;
foreach ($a as $k => $row) { $row[0] += 10; foreach ($row as $v) echo $v, ','; }
echo ':', ($a === $b) + 0, ':', $row[0], ':', $v, "\n";
