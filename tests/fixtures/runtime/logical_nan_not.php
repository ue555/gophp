<?php
$n = (1e308 * 10) - (1e308 * 10);
echo (!$n) === false;
echo (!!$n) === true;
if (!$n) { echo 'bad'; } else { echo '|else'; }
echo "\n";
