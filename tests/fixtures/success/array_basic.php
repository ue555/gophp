<?php
$a = [1, 'two', true, null, [5]];
$b = array('name' => 'PHP', 'nested' => array(8, 9),);
echo $a[0], ':', $a[1], ':', $a[2], ':', $a[3], ':', $a[4][0], "\n";
echo $b['name'], $b['nested'][1], [3, 4][1], "\n";
echo ([] === array()) + 0, "\n";
