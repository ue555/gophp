<?php
foreach ([null, true, false, 3, 2.5] as $v) echo isset($v[0]) + 0, empty($v[0]) + 0;
echo "\n";
