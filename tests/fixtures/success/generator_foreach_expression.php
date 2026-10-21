<?php
$arrays = [[1], [2]];
foreach ($arrays ?? [] as &$item) { $item[0] += 10; }
echo $arrays[0][0], '|', $arrays[1][0], '|', $item[0];
