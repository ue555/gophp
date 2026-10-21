<?php
[$a, $b] = [10, 20];
['left' => $c, 'nested' => [$d, $e]] = ['left' => 30, 'nested' => [40, 50]];
echo $a, '|', $b, '|', $c, '|', $d, '|', $e;
