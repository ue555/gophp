<?php
function source() { echo 'S'; return [3, 4, 5]; }
foreach (source() as $v) echo $v;
echo ':', $v, "\n";
