<?php
foreach ([1, 2, 3, 4] as $k => $v) {
    if ($v === 2) continue;
    if ($v === 4) break;
    echo $v;
}
echo ':', $k, ':', $v, "\n";
