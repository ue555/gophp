<?php
$log = ''; for ($i = 0; $i < 4; $log .= 'U', $i++) {
    if ($i === 1) continue;
    if ($i === 3) break;
    $log .= $i;
}
echo $log, ':', $i, "\n";
