<?php
$n = (1e308 * 10) - (1e308 * 10);
echo 'before|';
echo $n && (1 / 0);
echo 'after';
