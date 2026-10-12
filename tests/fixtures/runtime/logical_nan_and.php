<?php
$n = (1e308 * 10) - (1e308 * 10);
echo ($n && true) + 0;
echo (true && $n) + 0;
echo ($n && $n) + 0;
if ($n && true) { echo '|then'; }
echo "\n";
