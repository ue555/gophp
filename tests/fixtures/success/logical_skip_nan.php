<?php
$n = (1e308 * 10) - (1e308 * 10);
echo (false && $n) === false, (true || $n) === true;
echo (false && !$n) === false, (true || !!$n) === true;
if (false && $n) { echo !$n, 'must not print'; } else { echo '|else'; }
if (true || $n) { echo '|then'; } else { echo !$n, 'must not print'; }
echo "\n";
