<?php
for ($i = 0; $i < 3; ++$i) {
    foreach ([1, 2] as $v) { echo $i, $v; continue 2; }
    echo 'bad';
}
echo ':', $i, "\n";
$i = 0;
do { ++$i; foreach ([1] as $v) continue 2; echo 'bad'; } while ($i < 2);
echo $i, ':';
while (true) { foreach ([7] as $v) { while (true) { echo $v; break 3; } } echo 'bad'; }
echo "\n";
