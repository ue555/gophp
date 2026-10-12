<?php
$a = [1, 2, 3]; foreach ($a as $k => $v) { unset($a[$k]); echo $v; }
echo ':', count($a), ':', $k, ':', $v, "\n";
