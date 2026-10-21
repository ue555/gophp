<?php
$count = 0;
function nextKey() { global $count; $count++; return 'k'; }
$a = [];
$v = 3;
$a[nextKey()] =& $v;
$v = 5;
echo $count, '|', $a['k'];
