<?php
foreach ([1, 2, 3] as $v) {
    for ($i = 0; $i < 2; ++$i) { echo $v; continue 2; }
    echo 'bad';
}
echo ':', $v, ':', $i, "\n";
foreach ([4, 5] as $v) { do { echo $v; break 2; } while (1 / 0); }
echo ':', $v, "\n";
