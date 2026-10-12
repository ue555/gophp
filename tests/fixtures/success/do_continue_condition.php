<?php
$i = 0; $checks = 0;
do { ++$i; if ($i < 3) continue; echo $i; } while (++$checks < 3);
echo ':', $i, ':', $checks, "\n";
do { break; echo 'bad'; } while (1 / 0);
echo 'done', "\n";
