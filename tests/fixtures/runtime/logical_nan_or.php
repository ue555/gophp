<?php
$n = (1e308 * 10) - (1e308 * 10);
echo ($n || $n) + 0;
echo (true || $n) + 0;
echo (false || $n) + 0;
echo ($n || false) + 0, "\n";
